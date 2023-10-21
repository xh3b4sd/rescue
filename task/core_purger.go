package task

type prgcor struct {
	labl map[string]string
}

func (c *Core) Prg() *prgcor {
	return &prgcor{
		labl: *c,
	}
}

func (p *prgcor) Bypass() {
	delete(p.labl, Bypass)
}

func (p *prgcor) Cycles() {
	delete(p.labl, Cycles)
}

func (p *prgcor) Expiry() {
	delete(p.labl, Expiry)
}

func (p *prgcor) Method() {
	delete(p.labl, Method)
}

func (p *prgcor) Object() {
	delete(p.labl, Object)
}

func (p *prgcor) Worker() {
	delete(p.labl, Worker)
}
