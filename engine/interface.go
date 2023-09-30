package engine

import (
	"github.com/xh3b4sd/rescue/task"
)

type Interface interface {
	// Create submits a new task to the system.
	Create(tas *task.Task) error

	// Delete removes an existing task from the system.
	Delete(tas *task.Task) error

	// Exists expresses whether a task with the given metadata exists within the
	// queue. Given a task was created with metadata a, b and c. Exists will
	// return true if called with metadata b and c. If workers would want to
	// verify whether they still own a task, they could do the following call.
	//
	//     Exists(tas.All(metadata.ID, metadata.Owner))
	//
	Exists(tas *task.Task) (bool, error)

	// Expire is a background process that every worker should continously execute
	// in order to revoke ownership from tasks that have not been completed within
	// the specified time. Expire goes through the full list of available tasks
	// and revokes ownership from every task that was found to be expired. That
	// means that in a cluster of multiple workers, it takes only a single
	// functioning worker to call expire in order to keep existing tasks available
	// to be worked on.
	Expire() error

	// Extend can be called by task owners in order to extend the task's
	// expiration.
	Extend(tas *task.Task) error

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
}
