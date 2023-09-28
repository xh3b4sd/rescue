package task

func All() *Task {
	return &Task{
		Obj: TaskObj{
			Metadata: map[string]string{
				"*": "*",
			},
		},
	}
}
