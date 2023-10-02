# rescue

Reconciliation driven resource queue. The main primitive of `rescue` is a
`Task`. A Task describes some kind of job, a piece of work that has to be done.
The job description is defined by a set of labels, simple key-value pairs.

```go
type Task struct {
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
```



### Create Tasks

Anyone can create any task any time. The task producer must just have an
understanding of what consumers within the system are capable of. `Meta` and
`Root` of a queued task must always match with a consumer in order to be
processed.

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

Tasks can only be deleted by the workers that own the task they have been
assigned to. Task ownership cannot be cherry-picked. Deleting an expired task
causes an error on the consumer side.

```go
err := eng.Delete(tas)
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



### Search Tasks

Search provides the calling worker with an available task.

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
using `Ensure`.

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
docker run --rm --name redis-stack-redigo -p 6380:6379 -p 8002:8001 redis/redis-stack:latest
```
