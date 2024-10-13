package task

type exicrn struct {
	labl map[string]string
}

func (c *Cron) Exi() *exicrn {
	return &exicrn{
		labl: *c,
	}
}

func (e *exicrn) Adefer() bool {
	return e.labl[Adefer] != ""
}

func (e *exicrn) Aevery() bool {
	return e.labl[Aevery] != ""
}

func (e *exicrn) Aexact() bool {
	return e.labl[Aexact] != ""
}

func (e *exicrn) TickM1() bool {
	return e.labl[TickM1] != ""
}

func (e *exicrn) TickP1() bool {
	return e.labl[TickP1] != ""
}
