package engine

import (
	"time"

	"github.com/xh3b4sd/rescue/task"
	"github.com/xh3b4sd/tracer"
)

func (e *Engine) Create(tas *task.Task) error {
	var err error

	e.met.Engine.Create.Cal.Inc()

	o := func() error {
		err = e.create(tas)
		if err != nil {
			return tracer.Mask(err)
		}

		return nil
	}

	err = e.met.Engine.Create.Dur.Sin(o)
	if err != nil {
		e.met.Engine.Create.Err.Inc()
		return tracer.Mask(err)
	}

	return nil
}

func (e *Engine) create(tas *task.Task) error {
	var err error

	{
		if tas == nil {
			return tracer.Maskf(taskEmptyError, "Task must not be empty")
		}
		if tas.Core != nil {
			return tracer.Maskf(taskCoreError, "Task.Core must be empty")
		}
		if tas.Meta.Emp() {
			return tracer.Maskf(taskMetaEmptyError, "Task.Meta must not be empty")
		}

		if tas.Meta != nil && tas.Meta.Has(Res()) {
			return tracer.Maskf(labelReservedError, "Task.Meta must not contain reserved scheme rescue.io")
		}
		if tas.Root != nil && tas.Root.Has(Res()) {
			return tracer.Maskf(labelReservedError, "Task.Root must not contain reserved scheme rescue.io")
		}
	}

	// Creating tasks implies certain write operations on the task queue such as
	// adding a new task to a sorted set in redis. Due to such write operations
	// we need to ensure that only one process at a time can write information
	// back to the queue. Otherwise worker behaviour would be inconsistent and
	// the integrity of the queue could not be guaranteed.
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

	var tid int64
	{
		tid = time.Now().UTC().UnixNano()
	}

	{
		tas.Core = &task.Core{}
	}

	{
		tas.Core.Set().Object(tid)
	}

	{
		k := e.Keyfmt()
		v := task.ToString(tas)
		s := float64(tid)

		err = e.red.Sorted().Create().Score(k, v, s)
		if err != nil {
			return tracer.Mask(err)
		}
	}

	return nil
}
