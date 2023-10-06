package task

type exicor struct {
	Labl map[string]string
}

func (c *Core) Exi() *exicor {
	return &exicor{
		Labl: *c,
	}
}

func (e *exicor) Bypass() bool {
	return e.Labl[Bypass] != ""
}

func (e *exicor) Cycles() bool {
	return e.Labl[Cycles] != ""
}

func (e *exicor) Expiry() bool {
	return e.Labl[Expiry] != ""
}

func (e *exicor) Object() bool {
	return e.Labl[Object] != ""
}

func (e *exicor) Worker() bool {
	return e.Labl[Worker] != ""
}
