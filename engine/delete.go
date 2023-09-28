package engine

import (
	"github.com/xh3b4sd/rescue/key"
	"github.com/xh3b4sd/rescue/task"
	"github.com/xh3b4sd/rescue/validate"
	"github.com/xh3b4sd/tracer"
)

func (e *Engine) Delete(tas *task.Task) error {
	var err error

	e.met.Engine.Delete.Cal.Inc()

	o := func() error {
		err = e.delete(tas)
		if err != nil {
			return tracer.Mask(err)
		}

		return nil
	}

	err = e.met.Engine.Delete.Dur.Sin(o)
	if err != nil {
		e.met.Engine.Delete.Err.Inc()
		return tracer.Mask(err)
	}

	return nil
}

func (e *Engine) delete(tas *task.Task) error {
	var err error

	{
		err = validate.Empty(tas)
		if err != nil {
			return tracer.Mask(err)
		}
	}

	// Deleting tasks implies certain write operations on the task
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

	var cur *task.Task
	{
		k := key.Queue(e.que)
		s := tas.GetID()

		str, err := e.red.Sorted().Search().Score(k, s, s)
		if err != nil {
			return tracer.Mask(err)
		}

		if len(str) != 1 {
			e.met.Task.NotFound.Inc()
			return tracer.Mask(taskNotFoundError)
		}

		cur = task.FromString(str[0])
	}

	// We need to check tas against our actually stored tasks in the queue. It
	// might happen that tasks expire, causing ownership to be taken away from
	// workers. If workers try to delete their tasks after their tasks expired
	// within the queue, the attemtped delete operation must be considered
	// invalid. This ensures tas to be picked up again by another worker.
	//
	// Note that the comparison of current and desired state must exclude the
	// backoff, expire and version metadata. In case a task expired there might
	// be a worker who picked up the expired task already. If we would change
	// the backoff and version information in such a case, the worker having
	// picked up the expired task meanwhile could not delete the task properly
	// anymore, because the task state it knows changed within the system.
	var equ bool
	{
		tid := cur.GetID() == tas.GetID()
		own := cur.GetOwner() == tas.GetOwner() || tas.GetPrivileged()

		if tid && own {
			equ = true
		}
	}

	{
		if !equ {
			cur.IncBackoff(1)
			cur.IncVersion(1)
		}
	}

	{
		if !equ {
			k := key.Queue(e.que)
			v := task.ToString(cur)
			s := cur.GetID()

			_, err := e.red.Sorted().Update().Score(k, v, s)
			if err != nil {
				return tracer.Mask(err)
			}
		}
	}

	{
		if !equ {
			e.met.Task.Outdated.Inc()
			return tracer.Mask(taskOutdatedError)
		}
	}

	{
		k := key.Queue(e.que)
		s := tas.GetID()

		err = e.red.Sorted().Delete().Score(k, s)
		if err != nil {
			return tracer.Mask(err)
		}
	}

	return nil
}
