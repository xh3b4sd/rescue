package engine

import (
	"time"

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
		err := e.loc.Acquire()
		if err != nil {
			return nil, tracer.Mask(err)
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

	// Search for any task that defines the task delivery method "all". Such tasks
	// are meant to be processed by every worker within the network. We prioritize
	// such tasks and return them first, if we find them.
	for _, x := range lis {
		// Skip all scheduled task templates for further processing. Any task
		// template defining Task.Cron is meant to trigger time based task
		// scheduling for child tasks originating from that template. The template
		// itself is not meant to be processed by workers.
		if x.Cron != nil {
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
			continue
		}

		// Skip any task that does not define the task delivery method "all".
		if x.Node.Get(task.Method) != task.MthdAll {
			continue
		}

		var loc *local
		{
			loc = e.cac[x.Core.Get().Object()]
		}

		// Skip any task from our local copy that we already processed.
		if loc != nil && loc.don {
			continue
		}

		// Derive this task's creation timestamp from its object ID.
		var tim time.Time
		{
			tim = x.Core.Get().Object().Time()
		}

		// Skip any task that got created before this worker started to participate
		// within the network. Engine.pnt is the earliest point in time at which the
		// worker process came online, or the latest point in time of having
		// processed the oldest task broadcasted througout the network. If that
		// pointer is equal to, or after the creation time of the current task that
		// we do not track in our local cache already, then our rule is to not
		// process it. And so we skip the task that got created before the current
		// worker came online, and move on to the next task.
		if loc == nil && !e.pnt.Before(tim) {
			continue
		}

		var now time.Time
		{
			now = e.tim.Search()
		}

		// Skip any task that we are already processing within its specified time of
		// expiry. The tasks we are skipping here are either still being processed,
		// or failed, in which case we will pick them up again after local expiry.
		if loc != nil && loc.exp.After(now) {
			continue
		}

		// Remember the broadcasted task that this worker is processing right now
		// without assigning worker ownership within the underlying system. Also
		// remember the current expiry of this broadcasted task, so that we can
		// expire it locally and retry if necessary.
		{
			e.cac[x.Core.Get().Object()] = &local{exp: now.Add(e.exp)}
		}

		return x, nil
	}

	// Filter all tasks that have Task.Cron, Task.Gate or Task.Root defined.
	// Further, if the root task exists, delete the leaf task that defines it,
	// because the existing root task is meant to cover all the business logic
	// that its nested tasks are responsible for. Note that we collect the list
	// indices of the redundant tasks that we want to delete from the underlying
	// sorted set.
	var rem []int
	for i, x := range lis {
		// Remove all tasks with an active circuit breaker. Any task hitting their
		// defined maximum amount of execution attempts is effectively on hold until
		// their cycles counts are reset.
		if x.Core.Exi().Cancel() && x.Core.Get().Cycles() >= x.Core.Get().Cancel() {
			rem = append(rem, i)
			continue
		}

		// Remove all broadcasted tasks for further processing. Any task defining
		// delivery method "all" must have been addressed already above.
		if x.Node.Get(task.Method) == task.MthdAll {
			rem = append(rem, i)
			continue
		}

		// Remove all scheduled task templates for further processing. Any task
		// template defining Task.Cron and @every is meant to trigger time based
		// task scheduling for child tasks originating from that template. The
		// template itself is not meant to be processed by workers.
		if x.Cron != nil && x.Cron.Exi().Aevery() {
			rem = append(rem, i)
			continue
		}

		var now time.Time
		{
			now = e.tim.Search()
		}

		// Remove all scheduled tasks that define their exact execution to be in the
		// future. Any task template defining Task.Cron and @exact is meant to be
		// executed "exactly" at the specified time. So if the specified time is
		// pointing to the future still, we ignore it here.
		if x.Cron != nil && x.Cron.Exi().Aexact() && x.Cron.Get().Aexact().After(now) {
			rem = append(rem, i)
			continue
		}

		// Remove all deferred tasks that define their next execution to be in the
		// future. Any task defining Task.Cron and tick+1 is meant to be executed
		// after the specified time. So if the specified time is pointing to the
		// future still, we ignore it here.
		if x.Cron != nil && x.Cron.Exi().Adefer() && x.Cron.Exi().TickP1() && x.Cron.Get().TickP1().After(now) {
			rem = append(rem, i)
			continue
		}

		// Remove all trigger task templates for further processing. Any task
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

		if x.Root == nil {
			continue
		}

		// It is not permitted to set reserved labels to Task.Root from the outside.
		// The system though does that for scheduled tasks that are emitted on the
		// basis of task templates specifying Task.Cron and Task.Gate. Scheduled
		// tasks will contain the tree root's object ID, referencing the task
		// template. Scheduled tasks are not obsolete based on their tree structure
		// and template reference. So if we find a scheduled task we ignore it,
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
				s := x.Core.Get().Object().Float()

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

	// Each of the redundant tasks must be removed from our local copy once we
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

	for _, x := range lis {
		// We are looking for tasks which do not yet have an owner. So if there is
		// an owner assigned we ignore the task and move on to find another one.
		if x.Core.Get().Worker() != "" {
			continue
		}

		// The current task is not assigned to any worker. If this task's delivery
		// method is now set to "uni" and its target worker address is this current
		// worker, then we simply take it and assign it to the this current worker.
		// Note that we want to give tasks priority that are specifically addressed
		// to a particular worker. Tasks that can be processed by anyone are of
		// secondary importance in our system.
		if x.Node.Get(task.Method) == task.MthdUni && x.Node.Get(task.Worker) == e.wrk {
			tas = x
			break
		}
	}

	if tas == nil {
		for _, x := range lis {
			// We are looking for tasks which do not yet have an owner. So if there is
			// an owner assigned we ignore the task and move on to find another one.
			if x.Core.Get().Worker() != "" {
				continue
			}

			// The current task is not assigned to any worker. If this task's delivery
			// method is now set to "any", then we simply take it and assign it to this
			// current worker.
			if x.Node.Get(task.Method) == task.MthdAny {
				tas = x
				break
			}
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
		s := tas.Core.Get().Object().Float()

		_, err := e.red.Sorted().Update().Value(k, v, s)
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
