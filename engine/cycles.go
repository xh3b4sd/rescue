package engine

import (
	"fmt"

	"github.com/xh3b4sd/rescue/task"
	"github.com/xh3b4sd/tracer"
)

func (e *Engine) Cycles(tas *task.Task) error {
	var err error

	e.met.Engine.Cycles.Cal.Inc()

	o := func() error {
		err = e.cycles(tas)
		if err != nil {
			return tracer.Mask(err)
		}

		return nil
	}

	err = e.met.Engine.Cycles.Dur.Sin(o)
	if err != nil {
		e.met.Engine.Create.Err.Inc()
		return tracer.Mask(err)
	}

	return nil
}

func (e *Engine) cycles(tas *task.Task) error {
	var err error

	{
		if tas == nil {
			return tracer.Maskf(taskEmptyError, "Task must not be empty")
		}
		if tas.Core.Emp() {
			return tracer.Maskf(taskCoreError, "Task.Core must not be empty")
		}
		if !tas.Core.Has(Can()) {
			return tracer.Maskf(taskCoreError, "Task.Core does not define %s", task.Cancel)
		}
	}

	// Creating tasks implies certain write operations on the task queue such as
	// adding a new task to a sorted set in redis. Due to such write operations
	// we need to ensure that only one process at a time can write information
	// back to the queue. Otherwise worker behaviour would be inconsistent and
	// the integrity of the queue could not be guaranteed.
	{
		err := e.loc.Acquire()
		if err != nil {
			return tracer.Mask(err)
		}

		defer func() {
			err := e.loc.Release()
			if err != nil {
				e.lerror(tracer.Mask(err))
			}
		}()
	}

	var jsn []string
	{
		k := e.Keyfmt()
		s := tas.Core.Get().Object().Float()

		jsn, err = e.red.Sorted().Search().Score(k, s, s)
		if err != nil {
			return tracer.Mask(err)
		}
	}

	if len(jsn) != 1 || jsn[0] == "" {
		return tracer.Mask(fmt.Errorf("no task found for object ID %q", tas.Core.Map().Object()))
	}

	var upd *task.Task
	{
		upd = task.FromString(jsn[0])
	}

	{
		upd.Core.Prg().Expiry()
		upd.Core.Prg().Worker()
		upd.Core.Prg().Cycles()
	}

	{
		k := e.Keyfmt()
		v := task.ToString(upd)
		s := upd.Core.Get().Object().Float()

		_, err := e.red.Sorted().Update().Score(k, v, s)
		if err != nil {
			return tracer.Mask(err)
		}
	}

	return nil
}
