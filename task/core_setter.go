package task

import (
	"strconv"
	"time"

	"github.com/xh3b4sd/rescue/ticker"
)

type setcor struct {
	labl map[string]string
}

func (c *Core) Set() *setcor {
	return &setcor{
		labl: *c,
	}
}

func (s *setcor) Bypass(x bool) {
	s.labl[Bypass] = strconv.FormatBool(x)
}

func (s *setcor) Cycles(x int64) {
	s.labl[Cycles] = strconv.FormatInt(x, 10)
}

func (s *setcor) Expiry(x time.Time) {
	s.labl[Expiry] = x.Format(ticker.Layout)
}

func (s *setcor) Method(x string) {
	s.labl[Method] = x
}

func (s *setcor) Object(x int64) {
	s.labl[Object] = strconv.FormatInt(x, 10)
}

func (s *setcor) Worker(x string) {
	s.labl[Worker] = x
}
