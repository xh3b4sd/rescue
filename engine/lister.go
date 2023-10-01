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

	// We verify the given task metadata in order to ensure that no domain
	// metadata specific to the rescue internals are provided. That is to not let
	// arbitrary processes purposfully list tasks by ID because that ability could
	// be abused to take ownership from worker processes that may not be aware of
	// the corruption.
	{
		if tas == nil {
			return nil, tracer.Maskf(taskEmptyError, "Task must not be empty")
		}
		if tas.Core != nil {
			return nil, tracer.Maskf(taskCoreError, "Task.Core must be empty")
		}
		if tas.Meta.Emp() && tas.Root.Emp() {
			return nil, tracer.Maskf(taskMetaEmptyError, "either Task.Meta or Task.Root must not be empty")
		}

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

	var fil []*task.Task
	{
		for _, t := range lis {
			if t.Core != nil && tas.Core != nil && !t.Core.Has(*tas.Core) {
				continue
			}

			if t.Meta != nil && tas.Meta != nil && !t.Meta.Has(*tas.Meta) {
				continue
			}

			// TODO test exists for root
			if t.Root != nil && tas.Root != nil && !t.Root.Has(*tas.Root) {
				continue
			}

			{
				fil = append(fil, t)
			}
		}
	}

	return fil, nil
}
