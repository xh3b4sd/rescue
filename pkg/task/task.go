package task

type Task struct {
	Obj TaskObj
}

type TaskObj struct {
	// Metadata contains relevant information for task distribution. Below is
	// shown example metadata managed for internal purposes.
	//
	//     task.rescue.io/backoff          2
	//     task.rescue.io/expire           1612612238292387243
	//     task.rescue.io/id               1611318984211839461
	//     task.rescue.io/owner            al9qy
	//     task.rescue.io/version          4
	//
	// Any worker should be able to identify if they are able to execute on a
	// task successfully given the task metadata. Upon task creation certain
	// metadata can be set by producers in order to inform consumers about the
	// task's intention.
	//
	//     task.venturemark.co/action      delete
	//     task.venturemark.co/resource    timeline
	//
	Metadata map[string]string
}
