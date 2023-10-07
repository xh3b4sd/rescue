package task

import "time"

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
	byt, err := x.MarshalJSON()
	if err != nil {
		panic(err)
	}

	s.Labl[TickM1] = string(byt)
}

func (s *setcrn) TickP1(x time.Time) {
	byt, err := x.MarshalJSON()
	if err != nil {
		panic(err)
	}

	s.Labl[TickP1] = string(byt)
}
