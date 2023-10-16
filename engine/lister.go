package engine

import (
	"github.com/xh3b4sd/rescue/task"
	"github.com/xh3b4sd/tracer"
)

func (e *Engine) Lister(tas *task.Task) ([]*task.Task, error) {
	var err error
	var lis []*task.Task

	e.met.Engine.Lister.Cal.Inc()

	o := func() error {
		lis, err = e.lister(tas)
		if err != nil {
			return tracer.Mask(err)
		}

		return nil
	}

	err = e.met.Engine.Lister.Dur.Sin(o)
	if err != nil {
		e.met.Engine.Lister.Err.Inc()
		return nil, tracer.Mask(err)
	}

	return lis, nil
}

func (e *Engine) lister(tas *task.Task) ([]*task.Task, error) {
	var err error

	// We verify the given task labels to ensure that no core metadata specific to
	// the rescue internals are provided. That is to not let arbitrary processes
	// purposfully list tasks by ID because that ability could be abused to take
	// ownership from worker processes that may not be aware of the corruption.
	{
		if tas == nil {
			return nil, tracer.Maskf(taskEmptyError, "Task must not be empty")
		}
		if tas.Core != nil {
			return nil, tracer.Maskf(taskCoreError, "Task.Core must be empty")
		}
	}

	{
		crn := (tas.Cron == nil || tas.Cron.Emp())
		gat := (tas.Gate == nil || tas.Gate.Emp())
		met := (tas.Meta == nil || tas.Meta.Emp())
		roo := (tas.Root == nil || tas.Root.Emp())

		if crn && gat && met && roo {
			return nil, tracer.Maskf(taskMetaEmptyError, "at least one of [Task.Cron Task.Gate Task.Meta Task.Root] must be given")
		}
	}

	{
		if tas.Meta != nil && tas.Meta.Has(Res()) {
			return nil, tracer.Maskf(labelReservedError, "Task.Meta must not contain reserved scheme rescue.io")
		}
		if tas.Root != nil && tas.Root.Has(Res()) {
			return nil, tracer.Maskf(labelReservedError, "Task.Root must not contain reserved scheme rescue.io")
		}
	}

	// Listing all existing tasks implies certain read operations on the task
	// queue. For consistency reasons we need to ensure that only one process at
	// a time can read information from the queue. Otherwise worker behaviour
	// would be inconsistent and the integrity of the queue could not be
	// guaranteed.
	{
		err := e.red.Locker().Acquire()
		if err != nil {
			return nil, tracer.Mask(err)
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
			return nil, tracer.Mask(err)
		}
	}

	{
		e.met.Task.Inactive.Set(float64(len(lis)))
	}

	var fil []*task.Task
	for _, t := range lis {
		crn := tas.Cron.Emp() || (t.Cron != nil && t.Cron.Has(*tas.Cron))
		gat := tas.Gate.Emp() || (t.Gate != nil && t.Gate.Has(*tas.Gate))
		met := tas.Meta.Emp() || (t.Meta != nil && t.Meta.Has(*tas.Meta))
		roo := tas.Root.Emp() || (t.Root != nil && t.Root.Has(*tas.Root))

		if crn && gat && met && roo {
			fil = append(fil, t)
		}
	}

	return fil, nil
}
