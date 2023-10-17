package engine

import (
	"time"

	"github.com/xh3b4sd/rescue/task"
	"github.com/xh3b4sd/rescue/ticker"
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

	var tic *ticker.Ticker
	{
		tic, err = e.verCre(tas)
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
				e.lerror(tracer.Mask(err))
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

	if tas.Cron != nil {
		tas.Cron.Set().TickM1(tic.TickM1())
		tas.Cron.Set().TickP1(tic.TickP1())
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

func (e *Engine) verCre(tas *task.Task) (*ticker.Ticker, error) {
	{
		if tas == nil {
			return nil, tracer.Maskf(taskEmptyError, "Task must not be empty")
		}
		if tas.Core != nil {
			return nil, tracer.Maskf(taskCoreError, "Task.Core must be empty")
		}
		if tas.Meta == nil || tas.Meta.Emp() {
			return nil, tracer.Maskf(taskMetaEmptyError, "Task.Meta must not be empty")
		}
	}

	{
		if tas.Cron != nil && tas.Root != nil {
			return nil, tracer.Maskf(taskCronError, "Task.Cron and Task.Root must not be configured together")
		}
		if tas.Gate != nil && tas.Root != nil {
			return nil, tracer.Maskf(taskGateError, "Task.Gate and Task.Root must not be configured together")
		}
	}

	{
		if tas.Cron != nil && tas.Cron.Len() != 1 {
			return nil, tracer.Maskf(taskCronError, "Task.Cron must only be configured with one valid format")
		}
	}

	var tic *ticker.Ticker
	if tas.Cron != nil {
		tic = ticker.New(tas.Cron.Get().Aevery(), e.tim.Create())
	}

	{
		if tas.Cron != nil && tic.TickP1().IsZero() {
			return nil, tracer.Maskf(taskCronError, "Task.Cron format must be valid")
		}
	}

	if tas.Gate != nil {
		if tas.Gate.Has(Del()) {
			return nil, tracer.Maskf(labelReservedError, "Task.Gate must not contain reserved value [deleted]")
		}

		var tri bool
		var wai bool
		for _, v := range *tas.Gate {
			if v == task.Trigger {
				tri = true
			}

			if v == task.Waiting {
				wai = true
			}

			if tri && wai {
				return nil, tracer.Maskf(labelValueError, "Task.Gate must not contain both of the reserved values [trigger waiting] together")
			}
			if wai && tas.Cron != nil {
				return nil, tracer.Maskf(labelValueError, "Task.Gate must not contain the reserved value [waiting] if Task.Cron is configured")
			}
			if v != task.Trigger && v != task.Waiting {
				return nil, tracer.Maskf(labelValueError, "Task.Gate must only contain one of the reserved values [trigger waiting]")
			}
		}
	}

	{
		if tas.Gate != nil && tas.Gate.Has(Res()) {
			return nil, tracer.Maskf(labelReservedError, "Task.Gate must not contain reserved scheme rescue.io")
		}
		if tas.Meta != nil && tas.Meta.Has(Res()) {
			return nil, tracer.Maskf(labelReservedError, "Task.Meta must not contain reserved scheme rescue.io")
		}
		if tas.Root != nil && tas.Root.Has(Res()) {
			return nil, tracer.Maskf(labelReservedError, "Task.Root must not contain reserved scheme rescue.io")
		}
		if tas.Sync != nil && tas.Sync.Has(Res()) {
			return nil, tracer.Maskf(labelReservedError, "Task.Sync must not contain reserved scheme rescue.io")
		}
	}

	return tic, nil
}
