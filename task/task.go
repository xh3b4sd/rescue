package task

type Task struct {
	// Core contains systemically relevant information fundamental for task
	// distribution. Below is shown example metadata managed internally.
	//
	//     task.rescue.io/cycles    4
	//     task.rescue.io/expiry    2023-09-28T14:23:24.16198Z
	//     task.rescue.io/object    1611318984211839461
	//     task.rescue.io/worker    90dc68ba-4820-42ac-a924-2450388c15a6
	//
	Core *Core `json:"core,omitempty"`

	// Meta contains task specific information defined by the user. Any worker
	// should be able to identify whether they are able to execute on a task
	// successfully, given the task metadata. Upon task creation, certain metadata
	// can be set by producers in order to inform consumers about the task's
	// intention.
	//
	//     x.api.io/action    delete
	//     x.api.io/object    1234
	//
	Meta *Meta `json:"meta,omitempty"`

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
	//     x.api.io/object    1234
	//
	Root *Root `json:"root,omitempty"`
}
