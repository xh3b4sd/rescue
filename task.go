package rescue

type Task struct {
	Obj TaskObj
}

type TaskObj struct {
	// Metadata contains relevant information for task distribution and provides
	// a sort of instruction set. Any worker should be able to identify if they
	// are able to execute on a task successfully given the metadata.
	//
	//     task.venturemark.co/action      delete
	//     task.venturemark.co/expire      1611612297
	//     task.venturemark.co/id          fg4uu
	//     task.venturemark.co/owner       al9qy
	//     task.venturemark.co/resource    timeline
	//     task.venturemark.co/version     4
	//
	Metadata map[string]string
}
