package engine

import (
	"github.com/xh3b4sd/rescue/task"
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
		if tas == nil {
			return tracer.Maskf(taskEmptyError, "Task must not be empty")
		}
		if tas.Core.Emp() {
			return tracer.Maskf(taskCoreError, "Task.Core must not be empty")
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
		k := e.Keyfmt()
		s := float64(tas.Core.Get().Object())

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

	// We need to check the user given task against the actually stored tasks in
	// the queue. It might happen that tasks expire, causing ownership to be taken
	// away from workers. If workers try to delete their tasks after their tasks
	// expired within the queue, the attemtped delete operation must be considered
	// invalid. This ensures that the user given task can be picked up again by
	// another worker.
	//
	// Note that the comparison of current and desired state must exclude the
	// bypass, cycles and expiry metadata. In case a task expired there might be a
	// worker who picked up the expired task already, modifying the tasks metadata
	// further. Also, if we would change the metadata in such a case ourselves,
	// the worker having already claimed ownership of the task we are processing,
	// could then not delete the task properly anymore upon successful execution
	// on their end, because the task state this worker knows changed within the
	// system, and we would have broken the integrity of it.
	var equ bool
	{
		tid := cur.Core.Get().Object() == tas.Core.Get().Object()
		own := cur.Core.Get().Worker() == tas.Core.Get().Worker() || tas.Core.Get().Bypass()

		if tid && own {
			equ = true
		}
	}

	{
		if !equ {
			cur.Core.Set().Cycles(cur.Core.Get().Cycles() + 1)
		}
	}

	{
		if !equ {
			k := e.Keyfmt()
			v := task.ToString(cur)
			s := float64(cur.Core.Get().Object())

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
		k := e.Keyfmt()
		s := float64(tas.Core.Get().Object())

		err = e.red.Sorted().Delete().Score(k, s)
		if err != nil {
			return tracer.Mask(err)
		}
	}

	return nil
}
