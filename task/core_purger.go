package task

type purger struct {
	Labl map[string]string
}

func (c *Core) Prg() *purger {
	return &purger{
		Labl: *c,
	}
}

func (p *purger) Bypass() {
	delete(p.Labl, Bypass)
}

func (p *purger) Cycles() {
	delete(p.Labl, Cycles)
}

func (p *purger) Expiry() {
	delete(p.Labl, Expiry)
}

func (p *purger) Object() {
	delete(p.Labl, Object)
}

func (p *purger) Worker() {
	delete(p.Labl, Worker)
}
