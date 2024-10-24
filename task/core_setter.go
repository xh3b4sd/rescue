package task

import (
	"strconv"
	"time"

	"github.com/xh3b4sd/objectid"
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

func (s *setcor) Cancel(x int64) {
	s.labl[Cancel] = strconv.FormatInt(x, 10)
}

func (s *setcor) Cycles(x int64) {
	s.labl[Cycles] = strconv.FormatInt(x, 10)
}

func (s *setcor) Expiry(x time.Time) {
	s.labl[Expiry] = x.Format(ticker.Layout)
}

func (s *setcor) Object(x objectid.ID) {
	s.labl[Object] = string(x)
}

func (s *setcor) Worker(x string) {
	s.labl[Worker] = x
}
