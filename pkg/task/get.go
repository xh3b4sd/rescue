package task

import (
	"strconv"

	"github.com/xh3b4sd/rescue/pkg/metadata"
)

func (t *Task) GetBackoff() int {
	bac := t.Obj.Metadata[metadata.Backoff]

	b, err := strconv.Atoi(bac)
	if err != nil {
		panic(err)
	}

	return b
}

func (t *Task) GetExpire() int64 {
	exp := t.Obj.Metadata[metadata.Expire]

	e, err := strconv.ParseInt(exp, 10, 64)
	if err != nil {
		panic(err)
	}

	return e
}

func (t *Task) GetID() float64 {
	tid := t.Obj.Metadata[metadata.ID]

	i, err := strconv.ParseFloat(tid, 64)
	if err != nil {
		panic(err)
	}

	return i
}

func (t *Task) GetOwner() string {
	own := t.Obj.Metadata[metadata.Owner]

	o := own

	return o
}

func (t *Task) GetPrivileged() bool {
	pri := t.Obj.Metadata[metadata.Privileged]

	if pri == "" {
		return false
	}

	p, err := strconv.ParseBool(pri)
	if err != nil {
		panic(err)
	}

	return p
}

func (t *Task) GetVersion() int {
	ver := t.Obj.Metadata[metadata.Version]

	v, err := strconv.Atoi(ver)
	if err != nil {
		panic(err)
	}

	return v
}
