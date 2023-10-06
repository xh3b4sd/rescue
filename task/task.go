package task

type Task struct {
	// Core contains systemically relevant information fundamental for task
	// distribution. Below is shown example metadata managed internally.
	//
	//     task.rescue.io/cycles    4
	//     task.rescue.io/expiry    2023-09-28T14:23:24.161982Z
	//     task.rescue.io/object    1611318984211839461
	//     task.rescue.io/worker    90dc68ba-4820-42ac-a924-2450388c15a6
	//
	Core *Core `json:"core,omitempty"`

	// Cron contains optional scheduling information. A task may define to be
	// scheduled at an interval on the clock. Below is an example of a task that
	// is emitted at an interval of every 6 hours. That is every day at 00:00,
	// 06:00, 12:00 and 18:00, measured in UTC.
	//
	//     time.rescue.io/@every    6 hours
	//
	// Upon task creation, tasks defining an optional schedule will reflect the
	// previous and the next tick according to their configured interval.
	//
	//     time.rescue.io/tick-1    2023-09-28T12:00:00.000000Z
	//     time.rescue.io/tick+1    2023-09-28T18:00:00.000000Z
	//
	// The interval definition supports an opinionated set of duration units
	// expressed in more or less natural language. Note that the third column of
	// second and third order definitions for detailed schedules is not
	// implemented at the moment. Only quantity=1 and quantity=x are supported
	// right now. Important to show here right now is that the DSL allows for
	// certain extensions, if desirable.
	//
	//     quantity=1    quantity=x    at / on (at)
	//
	//     minute        15 minutes
	//     hour          4 hours       at **:30
	//     day           5 days        at 08:00
	//     week          2 weeks       on Wednesday (at 06:00)
	//     month         3 months      on the 15th (at 17:00)
	//
	// Note that scheduled tasks are emitted according to their specified
	// interval, never earlier, but arguably later to a neglectable extend.
	// Scheduling will always depend on the current conditions of the underlying
	// system. If hardware is overloaded or no worker process is running, then
	// scheduling might be affected considerably. If workers search for tasks
	// every 5 seconds, then scheduled tasks are likely to be executed with a
	// delay of a couple of seconds.
	Cron *Cron `json:"cron,omitempty"`

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
