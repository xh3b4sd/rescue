package task

type prgcrn struct {
	Labl map[string]string
}

func (c *Cron) Prg() *prgcrn {
	return &prgcrn{
		Labl: *c,
	}
}

func (p *prgcrn) Aevery() {
	delete(p.Labl, Aevery)
}

func (p *prgcrn) TickM1() {
	delete(p.Labl, TickM1)
}

func (p *prgcrn) TickP1() {
	delete(p.Labl, TickP1)
}
