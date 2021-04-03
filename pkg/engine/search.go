package engine

import (
	"fmt"
	"time"

	"github.com/xh3b4sd/tracer"

	"github.com/xh3b4sd/rescue/pkg/key"
	"github.com/xh3b4sd/rescue/pkg/task"
)

func (e *Engine) search() (*task.Task, error) {
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

		e.metric.Task.Queued.Set(float64(len(tks)))
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
			e.metric.Task.NotFound.Inc()
			return nil, tracer.Mask(noTaskError)
		}
	}

	{
		tsk.IncExpire(time.Now().UTC().UnixNano() + int64(e.ttl))
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

		str, err = e.redigo.Sorted().Search().Order(k, 0, -1)
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
