package engine

import (
	"github.com/xh3b4sd/rescue/task"
	"github.com/xh3b4sd/tracer"
)

func (e *Engine) Exists(tas *task.Task) (bool, error) {
	var err error
	var exi bool

	e.met.Engine.Exists.Cal.Inc()

	o := func() error {
		exi, err = e.exists(tas)
		if err != nil {
			return tracer.Mask(err)
		}

		return nil
	}

	err = e.met.Engine.Exists.Dur.Sin(o)
	if err != nil {
		e.met.Engine.Exists.Err.Inc()
		return false, tracer.Mask(err)
	}

	return exi, nil
}

func (e *Engine) exists(tas *task.Task) (bool, error) {
	var err error

	{
		if tas == nil {
			return false, tracer.Maskf(taskEmptyError, "Task must not be empty")
		}
	}

	{
		if tas.Emp() {
			return false, tracer.Maskf(taskMetaEmptyError, "at least one of [Task.Core Task.Cron Task.Gate Task.Host Task.Meta Task.Root] must be configured")
		}
	}

	{
		if tas.Meta != nil && tas.Meta.Has(Res()) {
			return false, tracer.Maskf(labelReservedError, "Task.Meta must not contain reserved scheme rescue.io")
		}
		if tas.Root != nil && tas.Root.Has(Res()) && !tas.Root.Has(Obj()) {
			return false, tracer.Maskf(labelReservedError, "Task.Root must not contain reserved scheme rescue.io")
		}
	}

	// Checking for existing tasks implies certain read operations on the task
	// queue. For consistency reasons we need to ensure that only one process at
	// a time can read information from the queue. Otherwise worker behaviour
	// would be inconsistent and the integrity of the queue could not be
	// guaranteed.
	{
		err := e.loc.Acquire()
		if err != nil {
			return false, tracer.Mask(err)
		}

		defer func() {
			err := e.loc.Release()
			if err != nil {
				e.lerror(tracer.Mask(err))
			}
		}()
	}

	var lis []*task.Task
	{
		lis, err = e.searchAll()
		if err != nil {
			return false, tracer.Mask(err)
		}
	}

	{
		e.met.Task.Inactive.Set(float64(len(lis)))
	}

	for _, x := range lis {
		if x.Has(tas) {
			return true, nil
		}
	}

	return false, nil
}
