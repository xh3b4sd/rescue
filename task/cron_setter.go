package task

import (
	"time"

	"github.com/xh3b4sd/rescue/ticker"
)

type setcrn struct {
	labl map[string]string
}

func (c *Cron) Set() *setcrn {
	return &setcrn{
		labl: *c,
	}
}

func (s *setcrn) Aevery(x string) {
	s.labl[Aevery] = x
}

func (s *setcrn) Aexact(x string) {
	s.labl[Aexact] = x
}

func (s *setcrn) TickM1(x time.Time) {
	s.labl[TickM1] = x.Format(ticker.Layout)
}

func (s *setcrn) TickP1(x time.Time) {
	s.labl[TickP1] = x.Format(ticker.Layout)
}
