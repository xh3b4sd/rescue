package task

type mapcrn struct {
	labl map[string]string
}

func (c *Cron) Map() *mapcrn {
	return &mapcrn{
		labl: *c,
	}
}

func (m *mapcrn) Aevery() string {
	return m.labl[Aevery]
}

func (m *mapcrn) TickM1() string {
	return m.labl[TickM1]
}

func (m *mapcrn) TickP1() string {
	return m.labl[TickP1]
}
