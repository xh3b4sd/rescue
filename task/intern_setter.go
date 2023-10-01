package task

import (
	"strconv"
	"time"
)

type setter struct {
	Labl map[string]string
}

func (i *Intern) Set() *setter {
	return &setter{
		Labl: *i,
	}
}

func (s *setter) Bypass(x bool) {
	s.Labl[Bypass] = strconv.FormatBool(x)
}

func (s *setter) Cycles(x int64) {
	s.Labl[Cycles] = strconv.FormatInt(x, 10)
}

func (s *setter) Expiry(x time.Time) {
	byt, err := x.MarshalJSON()
	if err != nil {
		panic(err)
	}

	s.Labl[Expiry] = string(byt)
}

func (s *setter) Object(x int64) {
	s.Labl[Object] = strconv.FormatInt(x, 10)
}

func (s *setter) Worker(x string) {
	s.Labl[Worker] = x
}
