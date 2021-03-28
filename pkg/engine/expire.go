package engine

import (
	"fmt"
	"time"

	"github.com/xh3b4sd/tracer"

	"github.com/xh3b4sd/rescue/pkg/key"
	"github.com/xh3b4sd/rescue/pkg/task"
)

func (e *Engine) expire() error {
	var err error

	// Expiring task ownership implies certain write operations on the task
	// queue such as updating the owner information. Due to such write
	// operations we need to ensure that only one process at a time can write
	// information back to the queue. Otherwise worker behaviour would be
	// inconsistent and the integrity of the queue could not be guaranteed.
	{
		err := e.redigo.Locker().Acquire()
		if err != nil {
			return tracer.Mask(err)
		}

		defer func() {
			err := e.redigo.Locker().Release()
			if err != nil {
				fmt.Println(err)
			}
		}()
	}

	var tks []*task.Task
	{
		tks, err = e.searchAll()
		if err != nil {
			return tracer.Mask(err)
		}
	}

	{
		for _, t := range tks {
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
			{
				exp = t.GetExpire()
				now = time.Now().UTC().UnixNano()
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
				k := key.Task
				v := task.ToString(t)
				s := t.GetID()

				_, err := e.redigo.Sorted().Update().Value(k, v, s)
				if err != nil {
					return tracer.Mask(err)
				}
			}

			{
				e.metric.Task.Expired.Inc()
			}
		}
	}

	return nil
}
