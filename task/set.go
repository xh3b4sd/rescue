package task

import (
	"strconv"
	"time"
)

type setter struct {
	Meta map[string]string
}

func (t *Task) Set() Setter {
	return &setter{
		Meta: t.Meta,
	}
}

func (s *setter) Bypass(x bool) {
	s.Meta[Bypass] = strconv.FormatBool(x)
}

func (s *setter) Cycles(x int64) {
	s.Meta[Cycles] = strconv.FormatInt(x, 10)
}

func (s *setter) Expiry(x time.Time) {
	byt, err := x.MarshalJSON()
	if err != nil {
		panic(err)
	}

	s.Meta[Expiry] = string(byt)
}

func (s *setter) Object(x int64) {
	s.Meta[Object] = strconv.FormatInt(x, 10)
}

func (s *setter) Worker(x string) {
	s.Meta[Worker] = x
}
