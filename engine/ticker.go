package engine

import (
	"time"

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
		if x.Cron == nil || x.Cron.Get().Aevery() == "" {
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
			s := float64(x.Core.Get().Object())

			_, err := e.red.Sorted().Update().Score(k, v, s)
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
				exi = true
				break
			}
		}

		// We are looking for a scheduled task that has been reconciled. The final
		// stage of reconciliation is task deletion. So if the scheduled task does
		// still exist, then we cannot schedule another one. Regardless, we have to
		// update the task template's tick+1 below.
		if !exi {
			// Create a new scheduled task with the template's Task.Meta, Task.Root
			// and Task.Sync reference. Note that scheduled tasks have a reserved
			// label reference of the parents object ID in their root directory,
			// pointing to the task template that defines their job description and
			// scheduling information.
			var t *task.Task
			{
				t = &task.Task{
					Meta: x.Meta,
					Root: &task.Root{
						task.Object: x.Core.Map().Object(),
					},
					Sync: x.Sync,
				}
			}

			var tid int64
			{
				tid = time.Now().UTC().UnixNano()
			}

			{
				t.Core = &task.Core{}
			}

			{
				t.Core.Set().Object(tid)
			}

			{
				k := e.Keyfmt()
				v := task.ToString(t)
				s := float64(tid)

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

		// We found a scheduled task that got scheduled just now based on its next
		// tick definition. Since the task got just scheduled, we move tick+1
		// forward based on the currently up to date calculation.
		{
			x.Cron.Set().TickP1(tic.TickP1())
		}

		{
			k := e.Keyfmt()
			v := task.ToString(x)
			s := float64(x.Core.Get().Object())

			_, err := e.red.Sorted().Update().Score(k, v, s)
			if err != nil {
				return tracer.Mask(err)
			}
		}
	}

	return nil
}
