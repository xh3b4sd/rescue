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
	//     ticket.pheobe.io/action    delete
	//     ticket.pheobe.io/refill    true
	//
	Metadata map[string]string
}

func (t *Task) Pref(pre ...string) *Task {
	var tas *Task
	{
		tas = &Task{}
	}

	for k, v := range t.Obj.Metadata {
		if prefix(pre, k) {
			if tas.Obj.Metadata == nil {
				tas.Obj.Metadata = map[string]string{}
			}

			tas.Obj.Metadata[k] = v
		}
	}

	if len(tas.Obj.Metadata) == 0 {
		return nil
	}

	return tas
}

// With returns a task containing all metadata identified by the list of
// provided keys. If metadata does not exist for a key, nil is returned. That
// means that the returned task will be nil, unless the complete range of
// desired metadata keys can be returned all together.
func (t *Task) With(key ...string) *Task {
	var tas *Task
	{
		tas = &Task{}
	}

	for _, k := range key {
		m, e := t.Obj.Metadata[k]
		if e {
			if tas.Obj.Metadata == nil {
				tas.Obj.Metadata = map[string]string{}
			}

			tas.Obj.Metadata[k] = m
		}
	}

	// We do only return a task containing the complete set of desired metadata.
	// If metadata for a given key way not found, we return nil.
	if len(key) != len(tas.Obj.Metadata) {
		return nil
	}

	return tas
}
