package task

type purger struct {
	Meta map[string]string
}

func (t *Task) Prg() Purger {
	return &purger{
		Meta: t.Meta,
	}
}

func (p *purger) Bypass() {
	delete(p.Meta, Bypass)
}

func (p *purger) Cycles() {
	delete(p.Meta, Cycles)
}

func (p *purger) Expiry() {
	delete(p.Meta, Expiry)
}

func (p *purger) Object() {
	delete(p.Meta, Object)
}

func (p *purger) Worker() {
	delete(p.Meta, Worker)
}
