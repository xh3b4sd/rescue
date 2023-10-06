# rescue

Reconciliation driven resource queue. The main primitives of `rescue` are the
`Engine` and `Task` types. A Task describes some kind of job, a piece of work
that has to be done. The job description is defined by a set of labels, simple
key-value pairs.

```go
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
	// interval, never earlier, but arguably later to a negligible extend.
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
```



### Create Tasks

`Engine.Create` submits a new task to the system. Anyone can create any task any
time. The task producer must just have an understanding of what consumers within
the system are capable of. `Task.Meta` and `Task.Root` of a queued task must
always match with a consumer in order to be processed. Scheduled tasks may be
created using `Task.Cron`.

```go
tas := &task.Task{
	Meta: &task.Meta{
		"test.api.io/action": "delete",
		"test.api.io/object": "1234",
	},
}

err := eng.Create(tas)
if err != nil {
	panic(err)
}
```



### Delete Tasks

`Engine.Delete` removes an existing task from the system. Tasks can only be
deleted by the workers that own the task they have been assigned to. Task
ownership cannot be cherry-picked. Deleting an expired task causes an error on
the consumer side, because the worker falsely believing to still be the task
owner, is operating based on an outdated copy of the task that changed meanwhile
within the system. Note that task templates defining `Task.Cron` may be deleted
by anyone using the bypass label, since those templates are never owned by any
worker.

```go
err := eng.Delete(tas)
if err != nil {
	panic(err)
}
```



### Verify Tasks

`Engine.Exists` expresses whether a task with the given label set exists within
the underlying queue. Given a task was created with metadata a, b and c, Exists
will return true if called with metadata b and c. If workers would want to
verify whether they still own a task, they could do the following call.
Basically, calling `tas.Core.All` returns a label set that matches all the given
label keys. That selective label set is then used by Exists to find any task
that matches the given query.

```go
tas := &task.Task{
	Core: tas.Core.All(task.Object, task.Worker),
}

exi, err := eng.Exists(tas)
if err != nil {
	panic(err)
}
```



### Expire Tasks

`Engine.Expire` is a background process that every worker should continously
execute in order to revoke ownership from tasks that have not been completed
within the specified time of expiry. Expire goes through the full list of
available tasks and revokes ownership from every task that was found to be
expired. That means that in a cluster of multiple workers, it takes only a
single functioning worker to call expire in order to keep existing tasks
available to be worked on.

```go
err := eng.Expire()
if err != nil {
	panic(err)
}
```



### Repeat Tasks

`Engine.Ticker` is an optional background process that every worker can
continously execute in order to emit scheduled tasks based on any task template
defining `Task.Cron`. Ticker goes through the full list of available tasks and
creates new tasks for any task template that is found to be due for scheduling
based on its next tick. That means that in a cluster of multiple workers, it
takes only a single functioning worker to call ticker in order to keep
scheduling recurring tasks for anyone to work on.

```go
err := eng.Ticker()
if err != nil {
	panic(err)
}
```



### Search Tasks

`Engine.Search` provides the calling worker with an available task.

```go
tas, err := eng.Search()
if err != nil {
	panic(err)
}
```



### Worker Interface

A common pattern to select the right task to work on would be some kind of
worker interface like shown below. Since any worker may be assigned to any task
at any time, the correct business logic must be invoked for the current task at
hand. That means a task must be identified using `Filter`, and then be processed
using `Ensure`. The asterisk may be used as a wildcard for matching any key or
value.

```go
func (w *Worker) Ensure(tas *task.Task) error {

	// business logic to cleanup properly

	return nil
}

func (w *Worker) Filter(tas *task.Task) bool {
	met := map[string]string{
		"test.api.io/action": "delete",
		"test.api.io/object": "*",
	}

	return tas.Meta.Has(met)
}
```



### Conformance Tests

```
docker run --rm --name redis-stack -p 6379:6379 -p 8001:8001 redis/redis-stack:latest
```

```
go test ./... -race -tags redis
```



### Redis Port

```
export REDIS_PORT=6380
```

```
docker run --rm --name redis-stack-rescue -p 6380:6379 -p 8002:8001 redis/redis-stack:latest
```
