package task

import (
	"time"

	"github.com/xh3b4sd/rescue/ticker"
)

type getcrn struct {
	Labl map[string]string
}

func (c *Cron) Get() *getcrn {
	return &getcrn{
		Labl: *c,
	}
}

func (g *getcrn) Aevery() string {
	return g.Labl[Aevery]
}

func (g *getcrn) TickM1() time.Time {
	tim, err := time.Parse(ticker.Layout, g.Labl[TickM1])
	if err != nil {
		panic(err)
	}

	return tim
}

func (g *getcrn) TickP1() time.Time {
	tim, err := time.Parse(ticker.Layout, g.Labl[TickP1])
	if err != nil {
		panic(err)
	}

	return tim
}
