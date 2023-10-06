package task

type prgcor struct {
	Labl map[string]string
}

func (c *Core) Prg() *prgcor {
	return &prgcor{
		Labl: *c,
	}
}

func (p *prgcor) Bypass() {
	delete(p.Labl, Bypass)
}

func (p *prgcor) Cycles() {
	delete(p.Labl, Cycles)
}

func (p *prgcor) Expiry() {
	delete(p.Labl, Expiry)
}

func (p *prgcor) Object() {
	delete(p.Labl, Object)
}

func (p *prgcor) Worker() {
	delete(p.Labl, Worker)
}
