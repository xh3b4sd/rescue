package task

type exicor struct {
	labl map[string]string
}

func (c *Core) Exi() *exicor {
	return &exicor{
		labl: *c,
	}
}

func (e *exicor) Bypass() bool {
	return e.labl[Bypass] != ""
}

func (e *exicor) Cycles() bool {
	return e.labl[Cycles] != ""
}

func (e *exicor) Expiry() bool {
	return e.labl[Expiry] != ""
}

func (e *exicor) Method() bool {
	return e.labl[Method] != ""
}

func (e *exicor) Object() bool {
	return e.labl[Object] != ""
}

func (e *exicor) Worker() bool {
	return e.labl[Worker] != ""
}
