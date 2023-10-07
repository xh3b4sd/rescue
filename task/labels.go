package task

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
