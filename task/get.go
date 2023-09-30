package task

import (
	"strconv"
	"time"
)

type getter struct {
	Meta map[string]string
}

func (t *Task) Get() Getter {
	return &getter{
		Meta: t.Meta,
	}
}

func (g *getter) Bypass() bool {
	if g.Meta[Bypass] == "" {
		return false
	}

	byp, err := strconv.ParseBool(g.Meta[Bypass])
	if err != nil {
		panic(err)
	}

	return byp
}

func (g *getter) Cycles() int64 {
	if g.Meta[Cycles] == "" {
		return 0
	}

	cyc, err := strconv.ParseInt(g.Meta[Cycles], 10, 64)
	if err != nil {
		panic(err)
	}

	return cyc
}

func (g *getter) Expiry() time.Time {
	var tim *time.Time
	{
		tim = &time.Time{}
	}

	err := tim.UnmarshalJSON([]byte(g.Meta[Expiry]))
	if err != nil {
		panic(err)
	}

	return *tim
}

func (g *getter) Object() int64 {
	cyc, err := strconv.ParseInt(g.Meta[Object], 10, 64)
	if err != nil {
		panic(err)
	}

	return cyc
}

func (g *getter) Worker() string {
	return g.Meta[Worker]
}
