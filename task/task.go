package task

type Task struct {
	// Meta contains relevant information for task distribution. Below is shown
	// example metadata managed for internal purposes.
	//
	//     task.rescue.io/cycles    4
	//     task.rescue.io/expiry    2023-09-28T14:23:24.16198Z
	//     task.rescue.io/object    1611318984211839461
	//     task.rescue.io/worker    al9qy
	//
	// Any worker should be able to identify if they are able to execute on a task
	// successfully given the task metadata. Upon task creation certain metadata
	// can be set by producers in order to inform consumers about the task's
	// intention. The asterisk may be used as a wildcard for matching any value.
	//
	//     api.naonao.io/action    delete
	//     api.naonao.io/object    *
	//
	Meta map[string]string
}
