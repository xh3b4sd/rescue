package engine

import (
	"time"

	"github.com/xh3b4sd/tracer"

	"github.com/xh3b4sd/rescue/pkg/key"
	"github.com/xh3b4sd/rescue/pkg/task"
	"github.com/xh3b4sd/rescue/pkg/validate"
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
		err = validate.Empty(tas)
		if err != nil {
			return tracer.Mask(err)
		}
		err = validate.Label(tas)
		if err != nil {
			return tracer.Mask(err)
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

	var tid float64
	{
		tid = float64(time.Now().UTC().UnixNano())
	}

	{
		tas.SetBackoff(0)
		tas.SetExpire(0)
		tas.SetID(tid)
		tas.SetOwner("")
		tas.SetVersion(1)
	}

	{
		k := key.Queue(e.que)
		v := task.ToString(tas)
		s := tid

		err = e.red.Sorted().Create().Score(k, v, s)
		if err != nil {
			return tracer.Mask(err)
		}
	}

	return nil
}
