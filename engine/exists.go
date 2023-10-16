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
		cor := tas.Core.Emp()
		crn := tas.Cron.Emp()
		gat := tas.Gate.Emp()
		met := tas.Meta.Emp()
		roo := tas.Root.Emp()

		if cor && crn && gat && met && roo {
			return false, tracer.Maskf(taskMetaEmptyError, "at least one of [Task.Core Task.Cron Task.Gate Task.Meta Task.Root] must be given")
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
		err := e.red.Locker().Acquire()
		if err != nil {
			return false, tracer.Mask(err)
		}

		defer func() {
			err := e.red.Locker().Release()
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

	for _, t := range lis {
		cor := tas.Core.Emp() || (t.Core != nil && t.Core.Has(*tas.Core))
		crn := tas.Cron.Emp() || (t.Cron != nil && t.Cron.Has(*tas.Cron))
		gat := tas.Gate.Emp() || (t.Gate != nil && t.Gate.Has(*tas.Gate))
		met := tas.Meta.Emp() || (t.Meta != nil && t.Meta.Has(*tas.Meta))
		roo := tas.Root.Emp() || (t.Root != nil && t.Root.Has(*tas.Root))

		if cor && crn && gat && met && roo {
			return true, nil
		}
	}

	return false, nil
}
