package engine

import (
	"github.com/xh3b4sd/rescue/task"
)

type Interface interface {
	// Create submits a new task to the system. Anyone can create any task any
	// time. The task producer must just have an understanding of what consumers
	// within the system are capable of. Task.Meta and Task.Root of a queued task
	// must always match with a consumer in order to be processed. Scheduled tasks
	// may be created using Task.Cron.
	Create(tas *task.Task) error

	// Delete removes an existing task from the system. Tasks can only be deleted
	// by the workers that own the task they have been assigned to. Task ownership
	// cannot be cherry-picked. Deleting an expired task causes an error on the
	// consumer side, because the worker falsely believing to still be the task
	// owner, is operating based on an outdated copy of the task that changed
	// meanwhile within the system. within the system. Note that task templates
	// defining Task.Cron may be deleted by anyone using the bypass label, since
	// those templates are never owned by any worker.
	Delete(tas *task.Task) error

	// Exists expresses whether a task with the given label set exists within the
	// underlying queue. Given a task was created with metadata a, b and c, Exists
	// will return true if called with metadata b and c. If workers would want to
	// verify whether they still own a task, they could do the following call.
	// Basically, calling `tas.Core.All` returns a label set that matches all the
	// given label keys. That selective label set is then used by Exists to find
	// any task that matches the given query.
	//
	//     Exists(&task.Task{Core: tas.Core.All(task.Object, task.Worker)})
	//
	Exists(tas *task.Task) (bool, error)

	// Expire is a background process that every worker should continously execute
	// in order to revoke ownership from tasks that have not been completed within
	// the specified time of expiry. Expire goes through the full list of
	// available tasks and revokes ownership from every task that was found to be
	// expired. That means that in a cluster of multiple workers, it takes only a
	// single functioning worker to call expire in order to keep existing tasks
	// available to be worked on.
	Expire() error

	// Extend can be called by workers owning a task in order to extend that
	// task's expiry. There should be a good reason to extend an ownership claim.
	// For instance, extending a task's expiry because the amount of work cannot
	// be done by a worker within the specified time may rather indicate a broken
	// system design. Processing tasks should not take forever. Instead,
	// processing tasks should require a limited amount of actions. It is totally
	// legitimate to process a task according to a budget and then return it to
	// the queue for another worker to pick it up later. Such a system design
	// should lead to more resilient and reliable software architectures, simply
	// because resource management is then designed to be eventually reconciled.
	Extend(tas *task.Task) error

	// Keyfmt returns the formatted key for this engine's queue, the underlying
	// sorted set in Redis.
	Keyfmt() string

	// Listen returns the TCP address in the form of host:port which the
	// underlying redis client is connected to.
	Listen() string

	// Lister fetches all existing tasks that match the given metadata. While
	// there are valid use cases for leveraging Lister, these use cases might be
	// rare, and the use of Lister may indicate a more fundamental flaw in the
	// underlying system design.
	Lister(tas *task.Task) ([]*task.Task, error)

	// Search provides the calling worker with an available task.
	Search() (*task.Task, error)

	// Engine.Ticker is an optional background process that every worker can
	// continously execute in order to emit scheduled tasks based on any task
	// template defining Task.Cron. Ticker goes through the full list of available
	// tasks and creates new tasks for any task template that is found to be due
	// for scheduling based on its next tick. That means that in a cluster of
	// multiple workers, it takes only a single functioning worker to call ticker
	// in order to keep scheduling recurring tasks for anyone to work on.
	Ticker() error
}
