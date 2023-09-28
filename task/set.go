package task

import (
	"strconv"

	"github.com/xh3b4sd/rescue/metadata"
)

func (t *Task) SetBackoff(b int) {
	var bac string
	{
		bac = strconv.Itoa(b)
	}

	t.Obj.Metadata[metadata.Backoff] = bac
}

func (t *Task) SetExpire(e int64) {
	var exp string
	{
		exp = strconv.FormatInt(e, 10)
	}

	t.Obj.Metadata[metadata.Expire] = exp
}

func (t *Task) SetID(i float64) {
	var tid string
	{
		tid = strconv.FormatFloat(i, 'f', -1, 64)
	}

	t.Obj.Metadata[metadata.ID] = tid
}

func (t *Task) SetOwner(o string) {
	t.Obj.Metadata[metadata.Owner] = o
}

func (t *Task) SetPrivileged(p bool) {
	t.Obj.Metadata[metadata.Privileged] = strconv.FormatBool(p)
}

func (t *Task) SetVersion(v int) {
	var ver string
	{
		ver = strconv.Itoa(v)
	}

	t.Obj.Metadata[metadata.Version] = ver
}
