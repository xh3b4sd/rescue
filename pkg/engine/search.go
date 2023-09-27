package engine

import (
	"sort"
	"time"

	"github.com/xh3b4sd/tracer"

	"github.com/xh3b4sd/rescue/pkg/key"
	"github.com/xh3b4sd/rescue/pkg/task"
)

func (e *Engine) Search() (*task.Task, error) {
	var err error
	var tas *task.Task

	e.met.Engine.Search.Cal.Inc()

	o := func() error {
		tas, err = e.search()
		if err != nil {
			return tracer.Mask(err)
		}

		return nil
	}

	err = e.met.Engine.Search.Dur.Sin(o)
	if err != nil {
		e.met.Engine.Search.Err.Inc()
		return nil, tracer.Mask(err)
	}

	return tas, nil
}

func (e *Engine) search() (*task.Task, error) {
	var err error

	// Searching for new tasks implies certain write operations on the task
	// queue such as updating the owner information. Due to such write
	// operations we need to ensure that only one process at a time can write
	// information back to the queue. Otherwise worker behaviour would be
	// inconsistent and the integrity of the queue could not be guaranteed.
	{
		err := e.red.Locker().Acquire()
		if err != nil {
			return nil, tracer.Mask(err)
		}

		defer func() {
			err := e.red.Locker().Release()
			if err != nil {
				e.log.Log(e.ctx, "level", "error", "message", "release failed", "stack", tracer.Stack(err))
			}
		}()
	}

	var lis []*task.Task
	{
		lis, err = e.searchAll()
		if err != nil {
			return nil, tracer.Mask(err)
		}

		e.met.Task.Queued.Set(float64(len(lis)))
	}

	{
		if len(lis) == 0 {
			e.met.Task.NotFound.Inc()
			return nil, tracer.Mask(taskNotFoundError)
		}
	}

	cur := map[string]int{}
	{
		for _, l := range lis {
			cur[l.GetOwner()]++
		}
	}

	var des map[string]int
	{
		des = e.bal.Opt(ensure(keys(cur), e.own), sum(cur))
	}

	var dev int
	{
		dev = des[e.own] - cur[e.own]
	}

	{
		if dev <= 0 {
			e.met.Task.NotFound.Inc()
			return nil, tracer.Mask(taskNotFoundError)
		}
	}

	var tas *task.Task
	{
		for _, t := range lis {
			// We are looking for tasks which do not yet have an owner. So if
			// there is an owner assigned we ignore the task and move on to find
			// another one.
			{
				if t.GetOwner() != "" {
					continue
				}
			}

			{
				tas = t
				break
			}
		}
	}

	{
		if tas == nil {
			e.met.Task.NotFound.Inc()
			return nil, tracer.Mask(taskNotFoundError)
		}
	}

	{
		tas.SetExpire(time.Now().UTC().UnixNano() + int64(e.ttl))
		tas.SetOwner(e.own)
		tas.IncVersion(1)
	}

	{
		k := key.Queue(e.que)
		v := task.ToString(tas)
		s := tas.GetID()

		_, err := e.red.Sorted().Update().Index(k, v, s)
		if err != nil {
			return nil, tracer.Mask(err)
		}
	}

	return tas, nil
}

func (e *Engine) searchAll() ([]*task.Task, error) {
	var err error

	var str []string
	{
		k := key.Queue(e.que)

		str, err = e.red.Sorted().Search().Order(k, 0, -1)
		if err != nil {
			return nil, tracer.Mask(err)
		}

		sort.Strings(str)
	}

	var lis []*task.Task
	{
		for _, s := range str {
			lis = append(lis, task.FromString(s))
		}
	}

	return lis, nil
}
