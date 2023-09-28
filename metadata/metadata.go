package metadata

const (
	// Backoff is the number of attempts workers tried to execute a given task.
	// This number is being incremented e.g. after ownership expiration, resulting
	// in rescheduling so that other workers can take over task execution. The
	// number of backoffs should be monitored so that producers and consumers can
	// be adjusted accordingly.
	Backoff = "task.rescue.io/backoff"

	// Expire is the unix timestamp in nanoseconds normalized to the UTC timezone.
	// This expiration time gets set once a worker takes the task for execution.
	// The worker taking the task becomes the owner and has a certain time window
	// to successfully execute the task. Should the task be executed successfully
	// within time the owner of the task can delete the task from the queue.
	// Should the task expire before being removed from the queue the worker
	// executing the task loses ownership of the task, causing the task to be
	// rescheduled.
	Expire = "task.rescue.io/expire"

	// ID is the task ID within the queue.
	ID = "task.rescue.io/id"

	// Owner is the worker executing the task.
	Owner = "task.rescue.io/owner"

	// Privileged is to bypass certain design specific safeguards. Privileged may
	// never be used, unless very good reasons demand it for special use cases.
	Privileged = "task.rescue.io/privileged"

	// Version is the number of changes made to the task during the task's
	// lifetime. Any write operation to the queue increments this number. The
	// initial task version is 1.
	Version = "task.rescue.io/version"
)
