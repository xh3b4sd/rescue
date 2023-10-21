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

func (g *getcrn) TickM1() time.Time {
	tic, err := time.Parse(ticker.Layout, g.labl[TickM1])
	if err != nil {
		panic(err)
	}

	return tic
}

func (g *getcrn) TickP1() time.Time {
	tic, err := time.Parse(ticker.Layout, g.labl[TickP1])
	if err != nil {
		panic(err)
	}

	return tic
}
