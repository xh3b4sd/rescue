package task

import (
	"strconv"
	"time"

	"github.com/xh3b4sd/rescue/ticker"
)

type getcor struct {
	Labl map[string]string
}

func (c *Core) Get() *getcor {
	return &getcor{
		Labl: *c,
	}
}

func (g *getcor) Bypass() bool {
	if g.Labl[Bypass] == "" {
		return false
	}

	byp, err := strconv.ParseBool(g.Labl[Bypass])
	if err != nil {
		panic(err)
	}

	return byp
}

func (g *getcor) Cycles() int64 {
	if g.Labl[Cycles] == "" {
		return 0
	}

	cyc, err := strconv.ParseInt(g.Labl[Cycles], 10, 64)
	if err != nil {
		panic(err)
	}

	return cyc
}

func (g *getcor) Expiry() time.Time {
	tim, err := time.Parse(ticker.Layout, g.Labl[Expiry])
	if err != nil {
		panic(err)
	}

	return tim
}

func (g *getcor) Object() int64 {
	cyc, err := strconv.ParseInt(g.Labl[Object], 10, 64)
	if err != nil {
		panic(err)
	}

	return cyc
}

func (g *getcor) Worker() string {
	return g.Labl[Worker]
}
