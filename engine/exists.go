package engine

import (
	"github.com/xh3b4sd/rescue/metadata"
	"github.com/xh3b4sd/rescue/task"
	"github.com/xh3b4sd/rescue/validate"
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
		err = validate.Empty(tas)
		if err != nil {
			return false, tracer.Mask(err)
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
				e.log.Log(e.ctx, "level", "error", "message", "release failed", "stack", tracer.Stack(err))
			}
		}()
	}

	var lis []*task.Task
	{
		lis, err = e.searchAll()
		if err != nil {
			return false, tracer.Mask(err)
		}

		e.met.Task.Queued.Set(float64(len(lis)))
	}

	{
		for _, t := range lis {
			// When checking for metadata the task fetched from the queue must
			// be given first since it contains all the metadata of the task
			// itself. The task given to Engine.Exists contains only the
			// relevant subset of metadata we want to match against.
			if metadata.Contains(t.Obj.Metadata, tas.Obj.Metadata) {
				return true, nil
			}
		}
	}

	return false, nil
}