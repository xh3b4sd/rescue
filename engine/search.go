package engine

import (
	"github.com/xh3b4sd/rescue/task"
	"github.com/xh3b4sd/tracer"
)

func (e *Engine) Search() (*task.Task, error) {
	var err error
	var tas *task.Task

	e.met.Engine.Search.Cal.Inc()

	o := func() error {
		tas, err = e.search()
		if err != nil {
			return tracer.Mask(err)
		}

		return nil
	}

	err = e.met.Engine.Search.Dur.Sin(o)
	if err != nil {
		e.met.Engine.Search.Err.Inc()
		return nil, tracer.Mask(err)
	}

	return tas, nil
}

func (e *Engine) search() (*task.Task, error) {
	var err error

	// Searching for new tasks implies certain write operations on the task
	// queue such as updating the owner information. Due to such write
	// operations we need to ensure that only one process at a time can write
	// information back to the queue. Otherwise worker behaviour would be
	// inconsistent and the integrity of the queue could not be guaranteed.
	{
		err := e.red.Locker().Acquire()
		if err != nil {
			return nil, tracer.Mask(err)
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
			return nil, tracer.Mask(err)
		}
	}

	{
		e.met.Task.Inactive.Set(float64(len(lis)))
	}

	if len(lis) == 0 {
		e.met.Task.NotFound.Inc()
		return nil, tracer.Mask(taskNotFoundError)
	}

	// Find all tasks that have a Task.Root defined. If the root task exists,
	// delete the task that defines it, because the existing root task is meant to
	// cover all the business logic that its nested tasks are responsible for.
	// Note that we collect the list indices of the redundant tasks that we delete
	// from the underlying sorted set.
	var rem []int
	for i, x := range lis {
		// Remove all scheduled task templates for further processing.
		if x.Cron != nil {
			rem = append(rem, i)
			continue
		}

		if x.Root == nil {
			continue
		}

		// It is not permitted to set reserved labels to Task.Root from the outside.
		// The system though does that for scheduled tasks that are emitted on the
		// basis of task templates specifying Task.Cron. Scheduled tasks will
		// contain the tree root's object ID, referencing the task template.
		// Scheduled tasks are not obsolete based on their tree structure and
		// template reference. So if we find a scheduled task we ignore them,
		// because we do not want to delete those.
		if x.Root.Len() == 1 && x.Root.Exi(task.Object) {
			continue
		}

		for j, y := range lis {
			// Skip the task we are processing right now. Here x and y are equal in
			// case i and j are the same.
			if i == j {
				continue
			}

			// Skip all the tasks that do not match the root description.
			if !y.Meta.Has(*x.Root) {
				continue
			}

			// Delete x since it was identified to be a nested task under the root
			// that is represented by task y.
			{
				k := e.Keyfmt()
				s := float64(x.Core.Get().Object())

				err = e.red.Sorted().Delete().Score(k, s)
				if err != nil {
					return nil, tracer.Mask(err)
				}
			}

			{
				e.met.Task.Obsolete.Inc()
			}

			{
				rem = append(rem, i)
			}
		}
	}

	// Each of the redundant task must be removed from our local copy once we
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
		e.met.Task.NotFound.Inc()
		return nil, tracer.Mask(taskNotFoundError)
	}

	// Calculate the balanced ownership that workers can claim.
	cur := map[string]int{}
	{
		for _, l := range lis {
			cur[l.Core.Get().Worker()]++
		}

		var des map[string]int
		{
			des = e.bal.Opt(ensure(keys(cur), e.wrk), sum(cur))
		}

		var dev int
		{
			dev = des[e.wrk] - cur[e.wrk]
		}

		if dev <= 0 {
			e.met.Task.NotFound.Inc()
			return nil, tracer.Mask(taskNotFoundError)
		}
	}

	var tas *task.Task
	for _, t := range lis {
		// We are looking for tasks which do not yet have an owner. So if
		// there is an owner assigned we ignore the task and move on to find
		// another one.
		{
			if t.Core.Get().Worker() != "" {
				continue
			}
		}

		{
			tas = t
			break
		}
	}

	if tas == nil {
		e.met.Task.NotFound.Inc()
		return nil, tracer.Mask(taskNotFoundError)
	}

	{
		tas.Core.Set().Expiry(e.tim.Search().Add(e.exp))
		tas.Core.Set().Worker(e.wrk)
	}

	{
		k := e.Keyfmt()
		v := task.ToString(tas)
		s := float64(tas.Core.Get().Object())

		_, err := e.red.Sorted().Update().Score(k, v, s)
		if err != nil {
			return nil, tracer.Mask(err)
		}
	}

	{
		e.met.Task.Parallel.Set(float64(cur[tas.Core.Get().Worker()] + 1))
	}

	return tas, nil
}

func (e *Engine) searchAll() ([]*task.Task, error) {
	var err error

	var str []string
	{
		k := e.Keyfmt()

		str, err = e.red.Sorted().Search().Order(k, 0, -1)
		if err != nil {
			return nil, tracer.Mask(err)
		}
	}

	var lis []*task.Task
	for _, s := range str {
		lis = append(lis, task.FromString(s))
	}

	return lis, nil
}
