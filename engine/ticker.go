package engine

import (
	"time"

	"github.com/xh3b4sd/objectid"
	"github.com/xh3b4sd/rescue/task"
	"github.com/xh3b4sd/rescue/ticker"
	"github.com/xh3b4sd/tracer"
)

func (e *Engine) Ticker() error {
	var err error

	e.met.Engine.Ticker.Cal.Inc()

	o := func() error {
		err = e.ticker()
		if err != nil {
			return tracer.Mask(err)
		}

		return nil
	}

	err = e.met.Engine.Ticker.Dur.Sin(o)
	if err != nil {
		e.met.Engine.Ticker.Err.Inc()
		return tracer.Mask(err)
	}

	return nil
}

func (e *Engine) ticker() error {
	var err error

	// Emitting scheduled tasks implies certain write operations on the task queue
	// such as adding a new task to a sorted set in redis. Due to such write
	// operations we need to ensure that only one process at a time can write
	// information back to the queue. Otherwise worker behaviour would be
	// inconsistent and the integrity of the queue could not be guaranteed.
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

	var lis []*task.Task
	{
		lis, err = e.searchAll()
		if err != nil {
			return tracer.Mask(err)
		}
	}

	{
		e.met.Task.Inactive.Set(float64(len(lis)))
	}

	if len(lis) == 0 {
		return nil
	}

	// Update Task.Cron for task templates that completed their most recent
	// reconciliation.
	for i, x := range lis {
		// We are looking for tasks which have a schedule. So if there is no
		// schedule defined, then we ignore the task and move on to find another
		// one.
		if x.Cron == nil || !x.Cron.Exi().Aevery() {
			continue
		}

		var now time.Time
		{
			now = e.tim.Ticker()
		}

		// We are looking for tasks that have been scheduled recently based on their
		// next tick. So if the next tick is already in the past, then we ignore the
		// task and move on to find another one.
		if x.Cron.Get().TickP1().Before(now) {
			continue
		}

		var tic *ticker.Ticker
		{
			tic = ticker.New(x.Cron.Get().Aevery(), now)
		}

		// We are looking for tasks that have to complete their reconciliation loop.
		// Our key indicator for that is an out of sync past tick. So if the tick-1
		// in the task template is equal to the currently calculated tick-1, then we
		// ignore the task and move on to find another one.
		if x.Cron.Get().TickM1().Equal(tic.TickM1()) {
			continue
		}

		var exi bool
		for j, y := range lis {
			// Skip the task we are processing right now. Here x and y are equal in
			// case i and j are the same.
			if i == j {
				continue
			}

			if y.Root == nil {
				continue
			}

			// If the scheduled task references the object ID of our task template,
			// then we found the task that is still being reconciled.
			if y.Root.Len() == 1 && y.Root.Has(*x.Core.All(task.Object)) {
				exi = true
				break
			}
		}

		// We are looking for a scheduled task that has been reconciled. The final
		// stage of reconciliation is task deletion. So if the scheduled task does
		// still exist, then we cannot update the template's tick-1 to the most
		// recent up to date point, and we ignore the task and move on to find
		// another one.
		if exi {
			continue
		}

		// We could not find the scheduled task anymore that was previously
		// reconciled. So now we can bring the task template's past tick back into
		// sync, since its most recent reconciliation at the previously defined
		// tick-1 got successfully processed.
		{
			x.Cron.Set().TickM1(tic.TickM1())
		}

		// Update the task template defining Task.Cron.
		{
			k := e.Keyfmt()
			v := task.ToString(x)
			s := x.Core.Get().Object().Float()

			_, err := e.red.Sorted().Update().Value(k, v, s)
			if err != nil {
				return tracer.Mask(err)
			}
		}
	}

	// Create tasks that are due for scheduling.
	for i, x := range lis {
		// We are looking for task templates which have a schedule. So if there is
		// no schedule defined, then we ignore the task and move on to find another
		// one.
		if x.Cron == nil || x.Cron.Get().Aevery() == "" {
			continue
		}

		var now time.Time
		{
			now = e.tim.Ticker()
		}

		// We are looking for task templates that have to be scheduled based on
		// their next tick. So if the next tick is still to come in the future, then
		// we ignore the task and move on to find another one.
		if x.Cron.Get().TickP1().After(now) {
			continue
		}

		var exi bool
		for j, y := range lis {
			// Skip the task we are processing right now. Here x and y are equal in
			// case i and j are the same.
			if i == j {
				continue
			}

			if y.Root == nil {
				continue
			}

			// If the scheduled task references the object ID of our task template,
			// then we found the task that is still being reconciled.
			if y.Root.Len() == 1 && y.Root.Has(*x.Core.All(task.Object)) {
				if y.Node.Get(task.Method) == task.MthdAny {
					exi = true
					break
				} else {
					// Delete y since it was identified to be a scheduled task defining
					// the delivery method "all" or "any", which implies that those tasks
					// become redundant at the new interval that we just crossed. And so
					// we delete the lingering task to replace it with the new one we are
					// about to create below.
					{
						k := e.Keyfmt()
						s := y.Core.Get().Object().Float()

						err = e.red.Sorted().Delete().Score(k, s)
						if err != nil {
							return tracer.Mask(err)
						}
					}
				}
			}
		}

		// We are looking for a scheduled task that has been reconciled. The final
		// stage of reconciliation is task deletion. So if the scheduled task does
		// still exist, then we cannot schedule another one. Regardless, we have to
		// update the task template's tick+1 below.
		if !exi {
			// Create a new scheduled task with the template's Task.Gate, Task.Meta,
			// Task.Root and Task.Sync reference. Note that scheduled tasks have a
			// reserved label reference of the parents object ID in their root
			// directory, pointing to the task template that defines their job
			// description and scheduling information.
			var t *task.Task
			{
				t = &task.Task{
					Core: &task.Core{},
					Gate: x.Gate,
					Meta: x.Meta,
					Node: x.Node,
					Root: &task.Root{
						task.Object: x.Core.Map().Object(),
					},
					Sync: x.Sync,
				}
			}

			var oid objectid.ID
			{
				oid = objectid.Random(objectid.Time(e.tim.Ticker()))
			}

			{
				t.Core.Set().Object(oid)
			}

			if t.Node == nil {
				t.Node = &task.Node{}
			}

			if t.Node.Get(task.Method) == "" {
				t.Node.Set(task.Method, task.MthdAny)
			}

			{
				k := e.Keyfmt()
				v := task.ToString(t)
				s := oid.Float()

				err = e.red.Sorted().Create().Score(k, v, s)
				if err != nil {
					return tracer.Mask(err)
				}
			}
		}

		// Update the task template defining Task.Cron using an up to date ticker
		// instance.
		var tic *ticker.Ticker
		{
			tic = ticker.New(x.Cron.Get().Aevery(), now)
		}

		// The task delivery method "all" does not have a mechanism to acknowledge
		// successful task completion. If all possible workers within the network
		// are involved we have no guarantee of execution everywhere due to
		// potential split brain scenarios and other synchronization issues that
		// distributed systems bring with them. Task processing with task delivery
		// method "all" is time based, meaning that every worker decides for
		// themselves whether to execute such a task based on the point in time the
		// task at hand got created, and the point in time the worker started
		// participating in the network. The network does not know the complete set
		// of workers within the network. What every worker knows for themselves
		// though, is what they should be responsible for. And so for broadcasted
		// tasks, scheduled tasks move tick-1 and tick+1 forward together. That
		// means there is no completion or acknowledgement for scheduled tasks if
		// they are delivered to all workers. We just fire at-least-once, on
		// schedule, and leave the rest to the workers.
		if x.Node.Get(task.Method) == task.MthdAll {
			x.Cron.Set().TickM1(tic.TickM1())
		}

		// We found a scheduled task that got scheduled just now based on its next
		// tick definition. Since the task got just scheduled, we move tick+1
		// forward based on the currently up to date calculation.
		{
			x.Cron.Set().TickP1(tic.TickP1())
		}

		{
			k := e.Keyfmt()
			v := task.ToString(x)
			s := x.Core.Get().Object().Float()

			_, err := e.red.Sorted().Update().Value(k, v, s)
			if err != nil {
				return tracer.Mask(err)
			}
		}
	}

	return nil
}
