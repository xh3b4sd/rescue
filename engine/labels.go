package engine

import "github.com/xh3b4sd/rescue/task"

func All() *task.Task {
	return &task.Task{
		Meta: &task.Meta{
			"*": "*",
		},
	}
}

func Can() map[string]string {
	return map[string]string{
		task.Cancel: "*",
	}
}

func Del() map[string]string {
	return map[string]string{
		"*": task.Deleted,
	}
}

func Met() map[string]string {
	return map[string]string{
		task.Method: "*",
	}
}

func Obj() map[string]string {
	return map[string]string{
		task.Object: "*",
	}
}

func Res() map[string]string {
	return map[string]string{
		"*rescue.io*": "*",
	}
}

func Tri() map[string]string {
	return map[string]string{
		"*": task.Trigger,
	}
}

func Wai() map[string]string {
	return map[string]string{
		"*": task.Waiting,
	}
}
