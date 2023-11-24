package engine

import (
	"time"

	"github.com/xh3b4sd/rescue/task"
	"github.com/xh3b4sd/tracer"
)

func (e *Engine) Expire() error {
	var err error

	e.met.Engine.Expire.Cal.Inc()

	o := func() error {
		err = e.expire()
		if err != nil {
			return tracer.Mask(err)
		}

		return nil
	}

	err = e.met.Engine.Expire.Dur.Sin(o)
	if err != nil {
		e.met.Engine.Expire.Err.Inc()
		return tracer.Mask(err)
	}

	return nil
}

func (e *Engine) expire() error {
	var err error

	// Expiring task ownership implies certain write operations on the task
	// queue such as updating the owner information. Due to such write
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

	var rem []int
	for i, x := range lis {
		// Skip all scheduled task templates for further processing. Any task
		// template defining Task.Cron is meant to trigger time based task
		// scheduling for child tasks originating from that template. The template
		// itself is not meant to be processed by workers.
		if x.Cron != nil {
			rem = append(rem, i)
			continue
		}

		// Skip all trigger task templates for further processing. Any task
		// template defining Task.Gate is meant to trigger event based task
		// scheduling for child tasks originating from that template. The template
		// itself is not meant to be processed by workers. Note that we are checking
		// whether Task.Gate has not any value "trigger", which means that if
		// Task.Gate is not empty, then its values can only either be "waiting" or
		// "deleted", which defines the trigger templates. Scheduled tasks defining
		// any value "trigger" in Task.Gate are the very tasks that workers should
		// process, because completion of those processed trigger tasks is what
		// causes the trigger template to create the gated task that is being onhold
		// until all triggers completed.
		if x.Gate != nil && !x.Gate.Has(Tri()) {
			rem = append(rem, i)
			continue
		}

		// Derive this task's creation timestamp from its object ID.
		var tim time.Time
		{
			tim = created(x.Core.Get().Object())
		}

		// Any lingering task is removed from the internal state once it is older
		// the configured retention period, which defaults to 1 week. This is just a
		// random guess on what is sensible, and since we want to do some house
		// keeping in order to prevent unnecessary state bloat, we just get rid of
		// it eventually. The assumption here right now is that tasks to be
		// processed by all workers within the network are either already processed,
		// or not relevant anymore beyond 1 week of creation.
		if e.tim.Search().Sub(tim) > e.cln {
			// Remove the irrelevant task from memory, if any.
			{
				delete(e.cac, x.Core.Map().Object())
			}

			// Remove the irrelevant task from the underlying queue.
			{
				k := e.Keyfmt()
				s := float64(x.Core.Get().Object())

				err = e.red.Sorted().Delete().Score(k, s)
				if err != nil {
					return tracer.Mask(err)
				}
			}

			{
				rem = append(rem, i)
			}
		}
	}

	// Each of the deleted tasks must be removed from our local copy once we
	// deleted the respective elements from the underlying sorted set.
	for i, x := range rem {
		j := x - i
		if j < len(lis)-1 {
			copy(lis[j:], lis[j+1:])
		}
		lis[len(lis)-1] = nil
		lis = lis[:len(lis)-1]
	}

	if len(lis) == 0 {
		return nil
	}

	cur := map[string]int{}
	for _, l := range lis {
		cur[l.Core.Get().Worker()]++
	}

	for _, x := range lis {
		// We are looking for tasks which have an owner. So if there is no
		// owner assigned we ignore the task and move on to find another
		// one.
		if x.Core.Get().Worker() == "" {
			continue
		}

		var exp time.Time
		var now time.Time
		var wrk string
		{
			exp = x.Core.Get().Expiry()
			now = e.tim.Expire()
			wrk = x.Core.Get().Worker()
		}

		// We are looking for tasks which are expired already. So if the task we
		// look at is not expired yet, we ignore it and move on to find another
		// one. In other words, if the current task's expiry is still about to
		// happen after the current time, then the task is not yet expired, and we
		// continue with the next task.
		if exp.After(now) {
			continue
		}

		{
			x.Core.Prg().Expiry()
			x.Core.Prg().Worker()
			x.Core.Set().Cycles(x.Core.Get().Cycles() + 1)
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

		{
			e.met.Task.Expired.Inc()
		}

		{
			cur[wrk]--
		}
	}

	var des map[string]int
	{
		des = e.bal.Opt(ensure(keys(cur), e.wrk), sum(cur))
	}

	var dev map[string]int
	{
		dev = e.bal.Dev(cur, des)
	}

	for _, x := range lis {
		// We are looking for tasks which have an owner that is supposed to
		// revoke their ownership. So if there is no revocation indicated
		// for the current owner we ignore the task and move on to find
		// another one.
		if dev[x.Core.Get().Worker()] == 0 {
			continue
		}

		var exp time.Time
		var now time.Time
		var wrk string
		{
			exp = x.Core.Get().Expiry()
			now = e.tim.Expire()
			wrk = x.Core.Get().Worker()
		}

		// We are looking for tasks which are expired already. So if the task we
		// look at is not expired yet, we ignore it and move on to find another
		// one. In other words, if the current task's expiry is still about to
		// happen after the current time, then the task is not yet expired, and we
		// continue with the next task.
		if exp.After(now) {
			continue
		}

		{
			x.Core.Prg().Expiry()
			x.Core.Prg().Worker()
			x.Core.Set().Cycles(x.Core.Get().Cycles() + 1)
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

		{
			e.met.Task.Expired.Inc()
		}

		{
			dev[wrk]--
		}
	}

	if sum(dev) != 0 {
		return tracer.Mask(taskNotRevokedError)
	}

	return nil
}
