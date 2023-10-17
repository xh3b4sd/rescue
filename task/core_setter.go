package task

import (
	"strconv"
	"time"

	"github.com/xh3b4sd/rescue/ticker"
)

type setcor struct {
	Labl map[string]string
}

func (c *Core) Set() *setcor {
	return &setcor{
		Labl: *c,
	}
}

func (s *setcor) Bypass(x bool) {
	s.Labl[Bypass] = strconv.FormatBool(x)
}

func (s *setcor) Cycles(x int64) {
	s.Labl[Cycles] = strconv.FormatInt(x, 10)
}

func (s *setcor) Expiry(x time.Time) {
	s.Labl[Expiry] = x.Format(ticker.Layout)
}

func (s *setcor) Object(x int64) {
	s.Labl[Object] = strconv.FormatInt(x, 10)
}

func (s *setcor) Worker(x string) {
	s.Labl[Worker] = x
}
