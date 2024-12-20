package task

import (
	"strconv"
	"time"

	"github.com/xh3b4sd/objectid"
	"github.com/xh3b4sd/rescue/ticker"
)

type getcor struct {
	labl map[string]string
}

func (c *Core) Get() *getcor {
	return &getcor{
		labl: *c,
	}
}

func (g *getcor) Bypass() bool {
	if g.labl[Bypass] == "" {
		return false
	}

	byp, err := strconv.ParseBool(g.labl[Bypass])
	if err != nil {
		panic(err)
	}

	return byp
}

func (g *getcor) Cancel() int64 {
	if g.labl[Cancel] == "" {
		return 0
	}

	can, err := strconv.ParseInt(g.labl[Cancel], 10, 64)
	if err != nil {
		panic(err)
	}

	return can
}

func (g *getcor) Cycles() int64 {
	if g.labl[Cycles] == "" {
		return 0
	}

	cyc, err := strconv.ParseInt(g.labl[Cycles], 10, 64)
	if err != nil {
		panic(err)
	}

	return cyc
}

func (g *getcor) Expiry() time.Time {
	exp, err := time.Parse(ticker.Layout, g.labl[Expiry])
	if err != nil {
		panic(err)
	}

	return exp
}

func (g *getcor) Object() objectid.ID {
	return objectid.ID(g.labl[Object])
}

func (g *getcor) Worker() string {
	return g.labl[Worker]
}
