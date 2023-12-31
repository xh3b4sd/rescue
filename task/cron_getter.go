package task

import (
	"time"

	"github.com/xh3b4sd/rescue/ticker"
)

type getcrn struct {
	labl map[string]string
}

func (c *Cron) Get() *getcrn {
	return &getcrn{
		labl: *c,
	}
}

func (g *getcrn) Aevery() string {
	return g.labl[Aevery]
}

func (g *getcrn) Aexact() time.Time {
	tim, err := time.Parse(ticker.Layout, g.labl[Aexact])
	if err != nil {
		panic(err)
	}

	return tim
}

func (g *getcrn) TickM1() time.Time {
	tim, err := time.Parse(ticker.Layout, g.labl[TickM1])
	if err != nil {
		panic(err)
	}

	return tim
}

func (g *getcrn) TickP1() time.Time {
	tim, err := time.Parse(ticker.Layout, g.labl[TickP1])
	if err != nil {
		panic(err)
	}

	return tim
}
