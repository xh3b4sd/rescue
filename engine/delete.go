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

	var loc *local
	{
		loc = e.cac[tas.Core.Map().Object()]
	}

	// Allow the local deletion of any broadcasted task that is not a task
	// template.
	if loc != nil {
		all := tas.Node.Get(task.Method) == task.MthdAll
		byp := tas.Core.Exi().Bypass()
		crn := tas.Cron == nil
		gat := tas.Gate == nil

		if all && !byp && crn && gat {
			// We set this worker's internal time pointer to the expiry of the oldest
			// local task that we track internally. We do this to respect the expiry
			// of broadcasted tasks indexed locally. Tasks may fail and have to be
			// picked up again. Any more broadcasted tasks defining the delivery
			// method "all" may be processed as well if they got created after the
			// task that we just completed, because we are processing everything in
			// first-in-first-out fashion.
			{
				e.pnt = first(unix(expiry(e.cac)))
			}

			// Since this worker did its part in processing the broadcasted task, we
			// can mark this task's local copy as done.
			{
				loc.don = true
			}

			{
				e.cac[tas.Core.Map().Object()] = loc
			}

			return nil
		}
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
			return tracer.Maskf(taskNotFoundError, tas.Core.Map().Object())
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
		return tracer.Maskf(taskOutdatedError, cur.Core.Map().Object())
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
		for _, x := range lis {
			// Any task that does not define Task.Gate is not a task template, and so
			// we ignore it and move on to the next task.
			if x.Gate == nil {
				continue
			}

			// Any task that uses the reserved value "trigger" is not a task template,
			// and so we ignore it and move on to the next task.
			if x.Gate.Has(Tri()) {
				continue
			}

			var gat []string
			{
				gat = x.Gate.Any(tas.Gate.Key()...).Key()
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
			for _, y := range gat {
				x.Gate.Set(y, task.Deleted)
			}

			if x.Sync != nil && tas.Sync != nil {
				var syn []string
				{
					syn = x.Sync.Any(tas.Sync.Key()...).Key()
				}

				for _, y := range syn {
					x.Sync.Set(y, tas.Sync.Get(y))
				}
			}

			// Any task template that does not contain any reserved value "waiting"
			// anymore does only contain reserved values "deleted". That means this
			// task template can cause the creation of its trigger task, causing the
			// task template to be reset for the next cycle.
			if !x.Gate.Has(Wai()) {
				var t *task.Task
				{
					t = &task.Task{
						Core: &task.Core{},
						Meta: x.Meta,
						Node: x.Node,
						Root: &task.Root{
							task.Object: x.Core.Map().Object(),
						},
						Sync: x.Sync,
					}
				}

				var tid int64
				{
					tid = e.tim.Delete().UnixNano()
				}

				{
					t.Core.Set().Object(tid)
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
					s := float64(tid)

					err = e.red.Sorted().Create().Score(k, v, s)
					if err != nil {
						return tracer.Mask(err)
					}
				}

				// Once all reserved values flipped from "waiting" to "deleted" within a
				// task template and the associated trigger task got created, reset all
				// reserved values back to "waiting" for the next cycle to begin.
				for _, y := range x.Gate.Key() {
					x.Gate.Set(y, task.Waiting)
				}
			}

			// Update the system state of the task template in the underlying sorted
			// set.
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
	}

	// Update any task template defining Task.Cron with the scheduled task data
	// specified in Task.Sync, if such data exists.
	if tas.Root != nil && tas.Root.Exi(task.Object) && tas.Sync != nil && !tas.Sync.Emp() {
		for _, x := range lis {
			if x.Core.Map().Object() != tas.Root.Get(task.Object) {
				continue
			}

			{
				x.Sync = tas.Sync
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
