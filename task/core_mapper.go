package task

type mapcor struct {
	labl map[string]string
}

func (c *Core) Map() *mapcor {
	return &mapcor{
		labl: *c,
	}
}

func (m *mapcor) Bypass() string {
	return m.labl[Bypass]
}

func (m *mapcor) Cycles() string {
	return m.labl[Cycles]
}

func (m *mapcor) Expiry() string {
	return m.labl[Expiry]
}

func (m *mapcor) Method() string {
	return m.labl[Method]
}

func (m *mapcor) Object() string {
	return m.labl[Object]
}

func (m *mapcor) Worker() string {
	return m.labl[Worker]
}
