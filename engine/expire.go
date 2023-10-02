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
		err := e.red.Locker().Acquire()
		if err != nil {
			return tracer.Mask(err)
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
			return tracer.Mask(err)
		}
	}

	{
		e.met.Task.Inactive.Set(float64(len(lis)))
	}

	if len(lis) == 0 {
		return nil
	}

	cur := map[string]int{}
	for _, l := range lis {
		cur[l.Core.Get().Worker()]++
	}

	for _, t := range lis {
		// We are looking for tasks which have an owner. So if there is no
		// owner assigned we ignore the task and move on to find another
		// one.
		if t.Core.Get().Worker() == "" {
			continue
		}

		var exp time.Time
		var now time.Time
		var wrk string
		{
			exp = t.Core.Get().Expiry()
			now = time.Now().UTC()
			wrk = t.Core.Get().Worker()
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
			t.Core.Prg().Expiry()
			t.Core.Prg().Worker()
			t.Core.Set().Cycles(t.Core.Get().Cycles() + 1)
		}

		{
			k := e.Keyfmt()
			v := task.ToString(t)
			s := float64(t.Core.Get().Object())

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

	for _, t := range lis {
		// We are looking for tasks which have an owner that is supposed to
		// revoke their ownership. So if there is no revocation indicated
		// for the current owner we ignore the task and move on to find
		// another one.
		if dev[t.Core.Get().Worker()] == 0 {
			continue
		}

		var exp time.Time
		var now time.Time
		var wrk string
		{
			exp = t.Core.Get().Expiry()
			now = time.Now().UTC()
			wrk = t.Core.Get().Worker()
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
			t.Core.Prg().Expiry()
			t.Core.Prg().Worker()
			t.Core.Set().Cycles(t.Core.Get().Cycles() + 1)
		}

		{
			k := e.Keyfmt()
			v := task.ToString(t)
			s := float64(t.Core.Get().Object())

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
