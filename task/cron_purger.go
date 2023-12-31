package task

type prgcrn struct {
	labl map[string]string
}

func (c *Cron) Prg() *prgcrn {
	return &prgcrn{
		labl: *c,
	}
}

func (p *prgcrn) Aevery() {
	delete(p.labl, Aevery)
}

func (p *prgcrn) Aexact() {
	delete(p.labl, Aexact)
}

func (p *prgcrn) TickM1() {
	delete(p.labl, TickM1)
}

func (p *prgcrn) TickP1() {
	delete(p.labl, TickP1)
}
