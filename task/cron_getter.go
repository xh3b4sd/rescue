package task

import "time"

type getcrn struct {
	Labl map[string]string
}

func (c *Cron) Get() *getcrn {
	return &getcrn{
		Labl: *c,
	}
}

func (g *getcrn) Aevery() string {
	return g.Labl[Aevery]
}

func (g *getcrn) TickM1() time.Time {
	var tim *time.Time
	{
		tim = &time.Time{}
	}

	err := tim.UnmarshalJSON([]byte(g.Labl[TickM1]))
	if err != nil {
		panic(err)
	}

	return *tim
}

func (g *getcrn) TickP1() time.Time {
	var tim *time.Time
	{
		tim = &time.Time{}
	}

	err := tim.UnmarshalJSON([]byte(g.Labl[TickP1]))
	if err != nil {
		panic(err)
	}

	return *tim
}
