package engine

import (
	"fmt"

	"github.com/xh3b4sd/tracer"

	"github.com/xh3b4sd/rescue/pkg/key"
	"github.com/xh3b4sd/rescue/pkg/task"
)

func (e *Engine) delete(tsk *task.Task) error {
	var err error

	// Deleting tasks implies certain write operations on the task
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

	var cur *task.Task
	{
		k := key.Task
		s := tsk.GetID()

		str, err := e.redigo.Sorted().Search().Score(k, s, s)
		if err != nil {
			return tracer.Mask(err)
		}

		if len(str) != 1 {
			return tracer.Mask(searchFailedError)
		}

		cur = task.FromString(str[0])
	}

	// We need to check tsk against our actually stored tasks in the queue. It
	// might happen that tasks expire, causing ownership to be taken away from
	// workers. If workers try to delete their tasks after their tasks expired
	// within the queue, the attemtped delete operation must be considered
	// invalid. This ensures tsk to be picked up again by another worker.
	//
	// Note that the comparison of current and desired state must exclude the
	// backoff and version metadata. In case a task expired there might be a
	// worker who picked up the expired task already. If we would change the
	// backoff and version information in such a case, the worker having picked
	// up the expired task meanwhile could not delete the task properly anymore,
	// because the task state it knows changed within the system. This is why
	// backoff and version metadata is excluded from the comparison below.
	var equ bool
	{
		exp := cur.GetExpire() == tsk.GetExpire()
		tid := cur.GetID() == tsk.GetID()
		own := cur.GetOwner() == tsk.GetOwner()

		if exp && tid && own {
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
			k := key.Task
			v := task.ToString(cur)
			s := cur.GetID()

			_, err := e.redigo.Sorted().Update().Value(k, v, s)
			if err != nil {
				return tracer.Mask(err)
			}
		}
	}

	{
		if !equ {
			e.metric.Task.Outdated.Inc()
			return tracer.Mask(taskOutdatedError)
		}
	}

	{
		k := key.Task
		s := tsk.GetID()

		err = e.redigo.Sorted().Delete().Score(k, s)
		if err != nil {
			return tracer.Mask(err)
		}
	}

	return nil
}
