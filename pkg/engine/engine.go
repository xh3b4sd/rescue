package engine

import (
	"fmt"
	"strings"
	"time"

	"github.com/xh3b4sd/logger"
	"github.com/xh3b4sd/redigo"
	"github.com/xh3b4sd/tracer"

	"github.com/xh3b4sd/rescue/pkg/key"
	"github.com/xh3b4sd/rescue/pkg/random"
	"github.com/xh3b4sd/rescue/pkg/task"
)

type Config struct {
	Logger logger.Interface
	Redigo redigo.Interface

	Expire time.Duration
	Owner  string
}

type Engine struct {
	logger logger.Interface
	redigo redigo.Interface

	expire time.Duration
	owner  string
}

func New(config Config) (*Engine, error) {
	if config.Logger == nil {
		return nil, tracer.Maskf(invalidConfigError, "%T.Logger must not be empty", config)
	}
	if config.Redigo == nil {
		return nil, tracer.Maskf(invalidConfigError, "%T.Redigo must not be empty", config)
	}

	if config.Expire == 0 {
		config.Expire = 30 * time.Second
	}
	if config.Owner == "" {
		config.Owner = random.MustNew()
	}

	e := &Engine{
		logger: config.Logger,
		redigo: config.Redigo,

		expire: config.Expire,
		owner:  config.Owner,
	}

	return e, nil
}

func (e *Engine) Create(tsk *task.Task) error {
	var err error

	{
		if tsk == nil {
			return tracer.Maskf(invalidTaskError, "task must not be nil")
		}
		if len(tsk.Obj.Metadata) == 0 {
			return tracer.Maskf(invalidTaskError, "metadata must not be empty")
		}
		for k := range tsk.Obj.Metadata {
			if strings.HasPrefix(k, "task.rescue.io") {
				return tracer.Maskf(invalidTaskError, "metadata must not contain reserved scheme task.rescue.io")
			}
		}
	}

	// Creating tasks implies certain write operations on the task queue such as
	// adding a new task to a sorted set in redis. Due to such write operations
	// we need to ensure that only one process at a time can write information
	// back to the queue. Otherwise worker behaviour would be inconsistent and
	// the integrity of the queue could not be guaranteed.
	{
		err := e.redigo.Locker().Acquire()
		if err != nil {
			return tracer.Mask(err)
		}

		defer func() {
			err := e.redigo.Locker().Release()
			if err != nil {
				fmt.Println(err)
			}
		}()
	}

	var tid float64
	{
		tid = float64(time.Now().UTC().UnixNano())
	}

	{
		tsk.SetBackoff(0)
		tsk.SetExpire(0)
		tsk.SetID(tid)
		tsk.SetOwner("")
		tsk.SetVersion(1)
	}

	{
		k := key.Task
		v := task.ToString(tsk)
		s := tid

		err = e.redigo.Sorted().Create().Element(k, v, s)
		if err != nil {
			return tracer.Mask(err)
		}
	}

	return nil
}

func (e *Engine) Delete(tsk *task.Task) error {
	var err error

	// Deleting tasks implies certain write operations on the task
	// queue such as updating the owner information. Due to such write
	// operations we need to ensure that only one process at a time can write
	// information back to the queue. Otherwise worker behaviour would be
	// inconsistent and the integrity of the queue could not be guaranteed.
	{
		err := e.redigo.Locker().Acquire()
		if err != nil {
			return tracer.Mask(err)
		}

		defer func() {
			err := e.redigo.Locker().Release()
			if err != nil {
				fmt.Println(err)
			}
		}()
	}

	var cur *task.Task
	{
		k := key.Task
		s := tsk.GetID()

		str, err := e.redigo.Sorted().Search().Score(k, s, s)
		if err != nil {
			return tracer.Mask(err)
		}

		if len(str) != 1 {
			return tracer.Mask(searchFailedError)
		}

		cur = task.FromString(str[0])
	}

	// We need to check tsk against our actually stored tasks in the queue. It
	// might happen that tasks expire, causing ownership to be taken away from
	// workers. If workers try to delete their tasks after their tasks expired
	// within the queue, the attemtped delete operation must be considered
	// invalid. This causes tsk to be picked up again by another worker.
	//
	// Note that the comparison of current and desired tasks must exclude the
	// backoff and version metadata. In case a task expired there might be a
	// worker who picked up the expired task already. If we would change the
	// backoff and version information in such a case, the worker having picked
	// up the expired task meanwhile could not delete the task properly anymore,
	// because the task state it knows changed within the system. This is why
	// backoff and version metadata is excluded from the comparison below.
	var equ bool
	{
		exp := cur.GetExpire() == tsk.GetExpire()
		tid := cur.GetID() == tsk.GetID()
		own := cur.GetOwner() == tsk.GetOwner()

		if exp && tid && own {
			equ = true
		}
	}

	{
		if !equ {
			cur.IncBackoff(1)
			cur.IncVersion(1)
		}
	}

	{
		if !equ {
			k := key.Task
			v := task.ToString(cur)
			s := cur.GetID()

			_, err := e.redigo.Sorted().Update().Value(k, v, s)
			if err != nil {
				return tracer.Mask(err)
			}
		}
	}

	{
		if !equ {
			return tracer.Mask(taskOutdatedError)
		}
	}

	{
		k := key.Task
		s := tsk.GetID()

		err = e.redigo.Sorted().Delete().Score(k, s)
		if err != nil {
			return tracer.Mask(err)
		}
	}

	return nil
}

func (e *Engine) Expire() error {
	var err error

	// Expiring task ownership implies certain write operations on the task
	// queue such as updating the owner information. Due to such write
	// operations we need to ensure that only one process at a time can write
	// information back to the queue. Otherwise worker behaviour would be
	// inconsistent and the integrity of the queue could not be guaranteed.
	{
		err := e.redigo.Locker().Acquire()
		if err != nil {
			return tracer.Mask(err)
		}

		defer func() {
			err := e.redigo.Locker().Release()
			if err != nil {
				fmt.Println(err)
			}
		}()
	}

	var tks []*task.Task
	{
		tks, err = e.searchAll()
		if err != nil {
			return tracer.Mask(err)
		}
	}

	{
		for _, t := range tks {
			// We are looking for tasks which have an owner. So if there is no
			// owner assigned we ignore the task and move on to find another
			// one.
			{
				if t.GetOwner() == "" {
					continue
				}
			}

			var exp int64
			var now int64
			{
				exp = t.GetExpire()
				now = time.Now().UTC().UnixNano()
			}

			// We are looking for tasks which are expired already. So if the
			// task we look at is not expired yet, we ignore it and move on to
			// find another one.
			{
				if now < exp {
					continue
				}
			}

			{
				t.IncBackoff(1)
				t.SetExpire(0)
				t.SetOwner("")
				t.IncVersion(1)
			}

			{
				k := key.Task
				v := task.ToString(t)
				s := t.GetID()

				_, err := e.redigo.Sorted().Update().Value(k, v, s)
				if err != nil {
					return tracer.Mask(err)
				}
			}
		}
	}

	return nil
}

func (e *Engine) Search() (*task.Task, error) {
	var err error

	// Searching for new tasks implies certain write operations on the task
	// queue such as updating the owner information. Due to such write
	// operations we need to ensure that only one process at a time can write
	// information back to the queue. Otherwise worker behaviour would be
	// inconsistent and the integrity of the queue could not be guaranteed.
	{
		err := e.redigo.Locker().Acquire()
		if err != nil {
			return nil, tracer.Mask(err)
		}

		defer func() {
			err := e.redigo.Locker().Release()
			if err != nil {
				fmt.Println(err)
			}
		}()
	}

	var tks []*task.Task
	{
		tks, err = e.searchAll()
		if err != nil {
			return nil, tracer.Mask(err)
		}
	}

	var tsk *task.Task
	{
		for _, t := range tks {
			// We are looking for tasks which do not yet have an owner. So if
			// there is an owner assigned we ignore the task and move on to find
			// another one.
			{
				if t.GetOwner() != "" {
					continue
				}
			}

			{
				tsk = t
				break
			}
		}
	}

	{
		if tsk == nil {
			return nil, tracer.Mask(noTaskError)
		}
	}

	{
		tsk.IncExpire(time.Now().UTC().UnixNano() + int64(e.expire))
		tsk.SetOwner(e.owner)
		tsk.IncVersion(1)
	}

	{
		k := key.Task
		v := task.ToString(tsk)
		s := tsk.GetID()

		_, err := e.redigo.Sorted().Update().Value(k, v, s)
		if err != nil {
			return nil, tracer.Mask(err)
		}
	}

	return tsk, nil
}

func (e *Engine) searchAll() ([]*task.Task, error) {
	var err error

	var str []string
	{
		k := key.Task

		str, err = e.redigo.Sorted().Search().Index(k, 0, -1)
		if err != nil {
			return nil, tracer.Mask(err)
		}
	}

	var tks []*task.Task
	{
		for _, s := range str {
			tks = append(tks, task.FromString(s))
		}
	}

	return tks, nil
}
