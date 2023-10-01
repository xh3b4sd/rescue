package task

import (
	"strconv"
	"time"
)

type getter struct {
	Labl map[string]string
}

func (i *Intern) Get() *getter {
	return &getter{
		Labl: *i,
	}
}

func (g *getter) Bypass() bool {
	if g.Labl[Bypass] == "" {
		return false
	}

	byp, err := strconv.ParseBool(g.Labl[Bypass])
	if err != nil {
		panic(err)
	}

	return byp
}

func (g *getter) Cycles() int64 {
	if g.Labl[Cycles] == "" {
		return 0
	}

	cyc, err := strconv.ParseInt(g.Labl[Cycles], 10, 64)
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

	err := tim.UnmarshalJSON([]byte(g.Labl[Expiry]))
	if err != nil {
		panic(err)
	}

	return *tim
}

func (g *getter) Object() int64 {
	cyc, err := strconv.ParseInt(g.Labl[Object], 10, 64)
	if err != nil {
		panic(err)
	}

	return cyc
}

func (g *getter) Worker() string {
	return g.Labl[Worker]
}
