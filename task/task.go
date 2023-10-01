package task

type Task struct {
	// TODO figure out whether to use Brkr and how

	// Brkr is an optional channel transparently carrying a mechanism to signal
	// the end of a deadline. Especially long running tasks may be at risk of
	// running way beyond the given task expiry. Claiming new tasks from the queue
	// should not be affected by the processing of claimed tasks. So in order to
	// signal task expiry to the outside and the inside, Brkr might be used to
	// coordinate relevant processes.
	Brkr <-chan struct{} `json:"-"`

	// TODO use Core for internal labels only

	// Core contains systemically relevant information fundamental for task
	// distribution. Below is shown example metadata managed internally.
	//
	//     task.rescue.io/cycles    4
	//     task.rescue.io/expiry    2023-09-28T14:23:24.16198Z
	//     task.rescue.io/object    1611318984211839461
	//     task.rescue.io/worker    al9qy
	//
	Core *Core `json:"core,omitempty"`

	// Meta contains task specific information defined by the user. Any worker
	// should be able to identify whether they are able to execute on a task
	// successfully, given the task metadata. Upon task creation, certain metadata
	// can be set by producers in order to inform consumers about the task's
	// intention. The asterisk may be used as a wildcard for matching any key or
	// value.
	//
	//        *naonao.io/action    delete
	//     api.naonao.io/object    *
	//
	Meta *Meta `json:"meta,omitempty"`

	// TODO implemnent Root as a new feature

	// Root allows to manage a tree of dependencies. Consider task x and y, where
	// x is the root of y.
	//
	//     x.api.io/object    1234
	//
	//     └   y.api.io/object    2345
	//     └   y.api.io/object    3456
	//     └   y.api.io/object    4567
	//
	// If a task for x is present it makes y obsolete. Scheduling and processing y
	// if x is present may cause conflicts that are hard to resolve. So y may
	// define x as root, causing Engine.Create and Engine.Search to neither
	// schedule nor process y if x happens to exist. In the described example,
	// task y may define x as root like shown below, for y to be discarded by the
	// system, if it happens to exist alongside x.
	//
	//     x.api.io/object    *
	//
	Root *Root `json:"root,omitempty"`
}
