package engine

import (
	"time"

	"github.com/xh3b4sd/tracer"

	"github.com/xh3b4sd/rescue/pkg/key"
	"github.com/xh3b4sd/rescue/pkg/task"
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
		if len(lis) == 0 {
			return nil
		}
	}

	cur := map[string]int{}
	{
		for _, l := range lis {
			cur[l.GetOwner()]++
		}
	}

	{
		for _, t := range lis {
			// We are looking for tasks which have an owner. So if there is no
			// owner assigned we ignore the task and move on to find another
			// one.
			{
				if t.GetOwner() == "" {
					continue
				}
			}

			var exp int64
			var now int64
			var own string
			{
				exp = t.GetExpire()
				now = time.Now().UTC().UnixNano()
				own = t.GetOwner()
			}

			// We are looking for tasks which are expired already. So if the
			// task we look at is not expired yet, we ignore it and move on to
			// find another one.
			{
				if now < exp {
					continue
				}
			}

			{
				t.IncBackoff(1)
				t.SetExpire(0)
				t.SetOwner("")
				t.IncVersion(1)
			}

			{
				k := key.Queue(e.que)
				v := task.ToString(t)
				s := t.GetID()

				_, err := e.red.Sorted().Update().Score(k, v, s)
				if err != nil {
					return tracer.Mask(err)
				}
			}

			{
				e.met.Task.Expired.Inc()
			}

			{
				cur[own]--
			}
		}
	}

	var des map[string]int
	{
		des = e.bal.Opt(ensure(keys(cur), e.own), sum(cur))
	}

	var dev map[string]int
	{
		dev = e.bal.Dev(cur, des)
	}

	{
		for _, t := range lis {
			// We are looking for tasks which have an owner that is supposed to
			// revoke their ownership. So if there is no revocation indicated
			// for the current owner we ignore the task and move on to find
			// another one.
			{
				cou := dev[t.GetOwner()]
				if cou == 0 {
					continue
				}
			}

			var exp int64
			var now int64
			var own string
			{
				exp = t.GetExpire()
				now = time.Now().UTC().UnixNano()
				own = t.GetOwner()
			}

			// We are looking for tasks which are expired already. So if the
			// task we look at is not expired yet, we ignore it and move on to
			// find another one.
			{
				if now < exp {
					continue
				}
			}

			{
				t.IncBackoff(1)
				t.SetExpire(0)
				t.SetOwner("")
				t.IncVersion(1)
			}

			{
				k := key.Queue(e.que)
				v := task.ToString(t)
				s := t.GetID()

				_, err := e.red.Sorted().Update().Score(k, v, s)
				if err != nil {
					return tracer.Mask(err)
				}
			}

			{
				e.met.Task.Expired.Inc()
			}

			{
				dev[own]--
			}
		}
	}

	{
		if sum(dev) != 0 {
			return tracer.Mask(taskNotRevokedError)
		}
	}

	return nil
}
