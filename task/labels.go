package task

const (
	// Method is the addressing strategy to deliver a task within the network of
	// worker nodes. Every single task may be processed with varying guarantees of
	// delivery, but always at-least-once.
	//
	//     all    delivered to all workers within the network (timeline based)
	//     any    delivered to any worker within the network (default method)
	//     mny    delivered to many specific workers within the network (not implemented)
	//     uni    delivered to a single specific worker within the network (worker based)
	//
	Method = "addr.rescue.io/method"
)

const (
	// Paging is the requeuing indicator a task may carry to allow workers to
	// process tasks again with the given paging pointer. This paging pointer is a
	// progress indicator that can be used to inform workers at which point work
	// has to be picked up again.
	Paging = "sync.rescue.io/paging"
)

const (
	// Bypass is to work around certain design specific safeguards. Bypass may
	// never be used, unless very good reasons demand it for special use cases.
	Bypass = "task.rescue.io/bypass"

	// Cycles is the number of attempts that workers tried to execute a given
	// task. This number is being incremented e.g. after ownership expiration,
	// resulting in rescheduling so that other workers can take over task
	// execution. The number of cycles should be monitored so that issues in
	// system design can be addressed, ensuring that producers and consumers are
	// balanced accordingly.
	Cycles = "task.rescue.io/cycles"

	// Expiry is the unix timestamp of a tasks expiration time. This expiration
	// time gets set once a worker takes the task for execution. The worker taking
	// the task becomes the owner and has a certain time window to successfully
	// execute the task. Should the task be executed successfully within the
	// apiration time window, then the owner of the task can delete the task from
	// the queue. Should the task expire before being removed from the queue, then
	// the worker executing the task loses ownership of the task, causing the task
	// to be rescheduled. Workers losing ownership should be self aware of that
	// fact and stop executing on the expired task.
	Expiry = "task.rescue.io/expiry"

	// Object is the identifier of the task within the queue.
	Object = "task.rescue.io/object"

	// Worker is the name of the worker executing the task.
	Worker = "task.rescue.io/worker"
)

const (
	// Aevery is the optional scheduling information that users may provide. See
	// Task.Cron for more information.
	Aevery = "time.rescue.io/@every"

	// TickM1 is the unix timestamp of the most recent interval at which a
	// scheduled task got emitted and then eventually resolved. TickM1 is only
	// updated to the next interval of tick-1, if the task emitted most recently
	// got successfully resolved. It may happen that the task emitted most
	// recently did not get resolved before scheduling for the next tick was
	// supposed to happen. Then TickM1 remains unchanged until the system was able
	// to resolve certain underlying issues.
	TickM1 = "time.rescue.io/tick-1"

	// TickP1 is the unix timestamp of the next upcoming interval at which a
	// scheduled task may be created again. TickP1 is always updated to the next
	// interval of tick+1, regardless of the conditions in the underlying system.
	TickP1 = "time.rescue.io/tick+1"
)

const (
	// MthdAll is the addressing method to deliver a task to every worker within
	// the network at the time of that particular task creation. Consider workers
	// A, B and C participating in the network at any point in time. Consider task
	// T to be delivered to A, B and C at point in time X. Workers A and B have
	// been online before time X. As soon as workers A and B see task T, they will
	// process it, because their local pointers are before time X. Once workers A
	// and B processed task T, they will update their local pointers and know that
	// their part was done. Worker C comes online after time X and does not need
	// to process task T.
	MthdAll = "all"

	// MthdAny is the addressing method to deliver a task to any worker within the
	// network. It does not matter which worker is processing the given task, as
	// long as it is being processed. This is the default method and does not have
	// to be specified.
	MthdAny = "any"

	// MthdUni is the addressing method to deliver a task to a specific worker
	// within the network. Using "uni" requires the accompanied usage of the core
	// label key "task.rescue.io/worker" for specifying a particular identifier.
	// Tasks routed via "uni" may expire like "any" task. In fact the only
	// difference between "any" and "uni" is that the "uni" task has a sticky task
	// ownership requirement, while "any" task may be picked up by some arbitrary
	// worker after expiry.
	MthdUni = "uni"
)

const (
	// Deleted is a reserved Task.Gate value that the system applies internally to
	// task templates. Any trigger task carrying a matching Task.Gate key will
	// update the accounting of the task template's watchlist. So upon trigger
	// task completion, the matching label flips from "waiting" to "deleted" as
	// soon as the "trigger" value was received for a given Task.Gate key.
	Deleted = "deleted"

	// Trigger is a reserved Task.Gate value specified by trigger tasks. Any
	// trigger task defining Task.Gate may define any key with the corresponding
	// reserved value "trigger".
	//
	// The keys defined by task templates may correspond with keys defined by
	// trigger task definitions of Task.Gate. If the corresponding reserved values
	// "trigger" and "waiting" match based on the keys of task templates and
	// trigger tasks, then the task template's value flips from "waiting" to
	// "deleted".
	Trigger = "trigger"

	// Waiting is a reserved Task.Gate value specified by trigger tasks. Any task
	// template defining Task.Gate may define any key with the corresponding
	// reserved value "waiting".
	//
	// The keys defined by task templates may correspond with keys defined by
	// trigger task definitions of Task.Gate. If the corresponding reserved values
	// "trigger" and "waiting" match based on the keys of task templates and
	// trigger tasks, then the task template's value flips from "waiting" to
	// "deleted".
	Waiting = "waiting"
)
