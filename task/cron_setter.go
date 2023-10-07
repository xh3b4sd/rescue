package task

import (
	"time"

	"github.com/xh3b4sd/rescue/ticker"
)

type setcrn struct {
	Labl map[string]string
}

func (c *Cron) Set() *setcrn {
	return &setcrn{
		Labl: *c,
	}
}

func (s *setcrn) Aevery(x string) {
	s.Labl[Aevery] = x
}

func (s *setcrn) TickM1(x time.Time) {
	s.Labl[TickM1] = x.Format(ticker.Layout)
}

func (s *setcrn) TickP1(x time.Time) {
	s.Labl[TickP1] = x.Format(ticker.Layout)
}
