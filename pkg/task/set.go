package task

import "strconv"

func (t *Task) SetBackoff(b int) {
	var bac string
	{
		bac = strconv.Itoa(b)
	}

	t.Obj.Metadata["task.rescue.io/backoff"] = bac
}

func (t *Task) SetExpire(e int64) {
	var exp string
	{
		exp = strconv.FormatInt(e, 10)
	}

	t.Obj.Metadata["task.rescue.io/expire"] = exp
}

func (t *Task) SetID(i float64) {
	var tid string
	{
		tid = strconv.FormatFloat(i, 'f', -1, 64)
	}

	t.Obj.Metadata["task.rescue.io/id"] = tid
}

func (t *Task) SetOwner(o string) {
	t.Obj.Metadata["task.rescue.io/owner"] = o
}

func (t *Task) SetVersion(v int) {
	var ver string
	{
		ver = strconv.Itoa(v)
	}

	t.Obj.Metadata["task.rescue.io/version"] = ver
}
