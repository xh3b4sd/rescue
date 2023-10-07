package task

import (
	"strconv"
	"time"
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
	byt, err := x.MarshalJSON()
	if err != nil {
		panic(err)
	}

	s.Labl[Expiry] = string(byt)
}

func (s *setcor) Object(x int64) {
	s.Labl[Object] = strconv.FormatInt(x, 10)
}

func (s *setcor) Worker(x string) {
	s.Labl[Worker] = x
}
