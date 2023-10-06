package task

type exicrn struct {
	Labl map[string]string
}

func (c *Cron) Exi() *exicrn {
	return &exicrn{
		Labl: *c,
	}
}

func (e *exicrn) Aevery() bool {
	return e.Labl[Aevery] != ""
}

func (e *exicrn) TickM1() bool {
	return e.Labl[TickM1] != ""
}

func (e *exicrn) TickP1() bool {
	return e.Labl[TickP1] != ""
}
