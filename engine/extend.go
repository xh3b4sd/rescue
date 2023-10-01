package engine

import (
	"time"

	"github.com/xh3b4sd/rescue/task"
	"github.com/xh3b4sd/tracer"
)

func (e *Engine) Extend(tas *task.Task) error {
	var err error

	e.met.Engine.Extend.Cal.Inc()

	o := func() error {
		err = e.extend(tas)
		if err != nil {
			return tracer.Mask(err)
		}

		return nil
	}

	err = e.met.Engine.Extend.Dur.Sin(o)
	if err != nil {
		e.met.Engine.Extend.Err.Inc()
		return tracer.Mask(err)
	}

	return nil
}

func (e *Engine) extend(tas *task.Task) error {
	{
		if tas == nil {
			return tracer.Maskf(taskEmptyError, "Task must not be empty")
		}
		if tas.Core.Emp() {
			return tracer.Maskf(taskCoreError, "Task.Core must not be empty")
		}
	}

	// Extending task expiry implies certain write operations on the task queue
	// such as updating the expiry information. Due to such write operations we
	// need to ensure that only one process at a time can write information back
	// to the queue. Otherwise worker behaviour would be inconsistent and the
	// integrity of the queue could not be guaranteed.
	{
		err := e.red.Locker().Acquire()
		if err != nil {
			return tracer.Mask(err)
		}

		defer func() {
			err := e.red.Locker().Release()
			if err != nil {
				e.log.Log(e.ctx, "level", "error", "message", "release failed", "stack", tracer.Stack(err))
			}
		}()
	}

	var cur *task.Task
	{
		k := e.Keyfmt()
		s := float64(tas.Core.Get().Object())

		str, err := e.red.Sorted().Search().Score(k, s, s)
		if err != nil {
			return tracer.Mask(err)
		}

		if len(str) != 1 {
			e.met.Task.NotFound.Inc()
			return tracer.Mask(taskNotFoundError)
		}

		cur = task.FromString(str[0])
	}

	// Tasks can only be extended by owners.
	{
		tid := cur.Core.Get().Object() == tas.Core.Get().Object()
		own := cur.Core.Get().Worker() == tas.Core.Get().Worker()

		if !tid || !own {
			e.met.Task.Outdated.Inc()
			return tracer.Mask(taskOutdatedError)
		}
	}

	{
		cur.Core.Set().Expiry(time.Now().UTC().Add(e.exp))
	}

	{
		k := e.Keyfmt()
		v := task.ToString(cur)
		s := float64(cur.Core.Get().Object())

		_, err := e.red.Sorted().Update().Score(k, v, s)
		if err != nil {
			return tracer.Mask(err)
		}
	}

	return nil
}
