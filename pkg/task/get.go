package task

import (
	"strconv"
)

func (t *Task) GetBackoff() int {
	bac := t.Obj.Metadata["task.rescue.io/backoff"]

	b, err := strconv.Atoi(bac)
	if err != nil {
		panic(err)
	}

	return b
}

func (t *Task) GetExpire() int64 {
	exp := t.Obj.Metadata["task.rescue.io/expire"]

	e, err := strconv.ParseInt(exp, 10, 64)
	if err != nil {
		panic(err)
	}

	return e
}

func (t *Task) GetID() float64 {
	tid := t.Obj.Metadata["task.rescue.io/id"]

	i, err := strconv.ParseFloat(tid, 64)
	if err != nil {
		panic(err)
	}

	return i
}

func (t *Task) GetOwner() string {
	own := t.Obj.Metadata["task.rescue.io/owner"]

	o := own

	return o
}

func (t *Task) GetVersion() int {
	ver := t.Obj.Metadata["task.rescue.io/version"]

	v, err := strconv.Atoi(ver)
	if err != nil {
		panic(err)
	}

	return v
}
