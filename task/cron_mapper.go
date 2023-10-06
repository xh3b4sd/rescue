package task

type mapcrn struct {
	Labl map[string]string
}

func (c *Cron) Map() *mapcrn {
	return &mapcrn{
		Labl: *c,
	}
}

func (m *mapcrn) Aevery() string {
	return m.Labl[Aevery]
}

func (m *mapcrn) TickM1() string {
	return m.Labl[TickM1]
}

func (m *mapcrn) TickP1() string {
	return m.Labl[TickP1]
}
