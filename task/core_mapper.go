package task

type mapcor struct {
	Labl map[string]string
}

func (c *Core) Map() *mapcor {
	return &mapcor{
		Labl: *c,
	}
}

func (m *mapcor) Bypass() string {
	return m.Labl[Bypass]
}

func (m *mapcor) Cycles() string {
	return m.Labl[Cycles]
}

func (m *mapcor) Expiry() string {
	return m.Labl[Expiry]
}

func (m *mapcor) Object() string {
	return m.Labl[Object]
}

func (m *mapcor) Worker() string {
	return m.Labl[Worker]
}
