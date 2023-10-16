package engine

import (
	"time"

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
				e.lerror(tracer.Mask(err))
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
		exi := cur.Core.Exi().Worker() && tas.Core.Exi().Worker() || tas.Core.Get().Bypass()
		own := cur.Core.Get().Worker() == tas.Core.Get().Worker() || tas.Core.Get().Bypass()
		tid := cur.Core.Get().Object() == tas.Core.Get().Object()

		if exi && own && tid {
			equ = true
		}
	}

	if !equ {
		cur.Core.Set().Cycles(cur.Core.Get().Cycles() + 1)
	}

	if !equ {
		k := e.Keyfmt()
		v := task.ToString(cur)
		s := float64(cur.Core.Get().Object())

		_, err := e.red.Sorted().Update().Score(k, v, s)
		if err != nil {
			return tracer.Mask(err)
		}
	}

	if !equ {
		e.met.Task.Outdated.Inc()
		return tracer.Mask(taskOutdatedError)
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

	// We want to update all the task templates that define matching keys for the
	// given trigger task inside Task.Gate, but only if the given trigger task
	// defines Task.Gate themselves. Any matching label key will have the
	// corresponding reserved value of either "deleted" or "waiting".
	if tas.Gate != nil && tas.Gate.Has(Tri()) {
		for _, t := range lis {
			// Any task that does not define Task.Gate is not a task template, and so
			// we ignore it and move on to the next task.
			if t.Gate == nil {
				continue
			}

			// Any task that uses the reserved value "trigger" is not a task template,
			// and so we ignore it and move on to the next task.
			if t.Gate.Has(Tri()) {
				continue
			}

			var gat []string
			{
				gat = t.Gate.Any(tas.Gate.Key()...).Key()
			}

			// Any task template that does not contain any of the given trigger task's
			// label keys is not the associated task template that we are looking for,
			// and so we ignore it and move on to the next task.
			if len(gat) == 0 {
				continue
			}

			// Since we found a matching task template that defines the given trigger
			// task's label keys including their corresponding reserved values
			// "waiting", we set the values of those keys to "deleted" and update the
			// system state of the underlying sorted set below.
			for _, x := range gat {
				t.Gate.Set(x, task.Deleted)
			}

			if t.Sync != nil && tas.Sync != nil {
				var syn []string
				{
					syn = t.Sync.Any(tas.Sync.Key()...).Key()
				}

				for _, x := range syn {
					t.Sync.Set(x, tas.Sync.Get(x))
				}
			}

			// Any task template that does not contain any reserved value "waiting"
			// anymore does only contain reserved values "deleted". That means this
			// task template can cause the creation of its trigger task, causing the
			// task template to be reset for the next cycle.
			if !t.Gate.Has(Wai()) {
				var tri *task.Task
				{
					tri = &task.Task{
						Meta: t.Meta,
						Root: &task.Root{
							task.Object: t.Core.Map().Object(),
						},
						Sync: t.Sync,
					}
				}

				var tid int64
				{
					tid = time.Now().UTC().UnixNano()
				}

				{
					tri.Core = &task.Core{}
				}

				{
					tri.Core.Set().Object(tid)
				}

				{
					k := e.Keyfmt()
					v := task.ToString(tri)
					s := float64(tid)

					err = e.red.Sorted().Create().Score(k, v, s)
					if err != nil {
						return tracer.Mask(err)
					}
				}

				// Once all reserved values flipped from "waiting" to "deleted" within a
				// task template and the associated trigger task got created, reset all
				// reserved values back to "waiting" for the next cycle to begin.
				for _, x := range t.Gate.Key() {
					t.Gate.Set(x, task.Waiting)
				}
			}

			// Update the system state of the task template in the underlying sorted
			// set.
			{
				k := e.Keyfmt()
				v := task.ToString(t)
				s := float64(t.Core.Get().Object())

				_, err := e.red.Sorted().Update().Score(k, v, s)
				if err != nil {
					return tracer.Mask(err)
				}
			}
		}
	}

	// Update any task template defining Task.Cron with the scheduled task data
	// specified in Task.Sync, if such data exists.
	if tas.Root != nil && tas.Root.Exi(task.Object) && tas.Sync != nil && !tas.Sync.Emp() {
		for _, t := range lis {
			if t.Core.Map().Object() != tas.Root.Get(task.Object) {
				continue
			}

			{
				t.Sync = tas.Sync
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
