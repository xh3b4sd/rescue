package engine

import (
	"strconv"
	"time"

	"github.com/xh3b4sd/objectid"
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

	var oid objectid.ID
	{
		oid = objectid.Random(objectid.Time(e.tim.Create()))
	}

	if tas.Core == nil {
		tas.Core = &task.Core{}
	}

	{
		tas.Core.Set().Object(oid)
	}

	if tas.Cron != nil && tas.Cron.Exi().Aevery() {
		tas.Cron.Set().TickM1(tic.TickM1())
		tas.Cron.Set().TickP1(tic.TickP1())
	}

	if tas.Node == nil {
		tas.Node = &task.Node{}
	}

	if tas.Node.Get(task.Method) == "" {
		tas.Node.Set(task.Method, task.MthdAny)
	}

	{
		k := e.Keyfmt()
		v := task.ToString(tas)
		s := oid.Float()

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
		if tas.Core != nil && !(tas.Core.Len() == 1 && tas.Core.Has(Can())) {
			return nil, tracer.Maskf(taskCoreError, "Task.Core must be empty")
		}
		if tas.Core != nil && !(tas.Core.Has(Can()) && natNum(tas.Core.Map().Cancel())) {
			return nil, tracer.Maskf(taskCoreError, "Task.Core does not define a positive number for %s", task.Cancel)
		}
		if tas.Meta == nil || tas.Meta.Emp() {
			return nil, tracer.Maskf(taskMetaEmptyError, "Task.Meta must not be empty")
		}
	}

	if tas.Node != nil {
		if !tas.Node.Has(Met()) {
			return nil, tracer.Maskf(taskHostError, "Task.Node must contain reserved key [%s]", task.Method)
		}

		var met string
		{
			met = tas.Node.Get(task.Method)
		}

		if met == task.MthdAll && tas.Node.Len() != 1 {
			return nil, tracer.Maskf(taskHostError, `Task.Node must not contain any more labels if delivery method "all" is configured`)
		}

		if met == task.MthdAny && tas.Node.Len() != 1 {
			return nil, tracer.Maskf(taskHostError, `Task.Node must not contain any more labels if delivery method "any" is configured`)
		}

		if met == task.MthdUni && !tas.Node.Exi(task.Worker) {
			return nil, tracer.Maskf(taskHostError, `Task.Node must only contain reserved keys [%s %s] if delivery method "uni" is configured`, task.Method, task.Worker)
		}
		if met == task.MthdUni && tas.Node.Len() != 2 {
			return nil, tracer.Maskf(taskHostError, `Task.Node must only contain reserved keys [%s %s] if delivery method "uni" is configured`, task.Method, task.Worker)
		}

		for k, v := range *tas.Node {
			if k == task.Method && v != task.MthdAll && v != task.MthdAny && v != task.MthdUni {
				return nil, tracer.Maskf(labelValueError, "Task.Node must only contain one of the reserved values [%s %s %s]", task.MthdAll, task.MthdAny, task.MthdUni)
			}
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

	var now time.Time
	{
		now = e.tim.Create()
	}

	var tic *ticker.Ticker
	if tas.Cron != nil {
		if tas.Cron.Exi().Aevery() && tas.Cron.Exi().Aexact() {
			return nil, tracer.Maskf(taskCronError, "Task.Cron must not define @every and @exact together")
		}

		if tas.Cron.Len() != 1 {
			return nil, tracer.Maskf(taskCronError, "Task.Cron must only be configured with one valid format")
		}

		if tas.Cron.Exi().Adefer() {
			return nil, tracer.Maskf(taskCronError, "Task.Cron must define @defer at task creation")
		}

		if tas.Cron.Exi().Aevery() && !tas.Cron.Exi().Aexact() {
			{
				tic = ticker.New(tas.Cron.Get().Aevery(), now)
			}

			if tic.TickP1().IsZero() {
				return nil, tracer.Maskf(taskCronError, "Task.Cron format must be valid, got @every = %q", tas.Cron.Get().Aevery())
			}
		}

		if !tas.Cron.Exi().Aevery() && tas.Cron.Exi().Aexact() {
			tim, err := time.Parse(ticker.Layout, tas.Cron.Map().Aexact())
			if err != nil {
				return nil, tracer.Mask(err)
			}

			if !tim.After(now) {
				return nil, tracer.Maskf(taskCronError, "Task.Cron @exact must be in the future")
			}
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
		if tas.Sync != nil && tas.Sync.All(task.Paging).Len() != tas.Sync.All("*rescue.io*").Len() {
			return nil, tracer.Maskf(labelReservedError, "Task.Sync must not contain reserved scheme rescue.io, other than generic paging pointer")
		}
	}

	return tic, nil
}

func natNum(s string) bool {
	num, err := strconv.Atoi(s)
	return err == nil && num > 0
}
