package task

type Public map[string]string

func (p *Public) All(key ...string) *Public {
	pub := Public(All(*p, key...))
	return &pub
}

func (p *Public) Any(key ...string) *Public {
	pub := Public(Any(*p, key...))
	return &pub
}

func (p *Public) Emp() bool {
	if p == nil {
		return true
	}

	pub := *p
	return len(pub) == 0
}

func (p *Public) Get(key string) string {
	pub := *p
	return pub[key]
}

func (p *Public) Has(lab map[string]string) bool {
	return Has(*p, lab)
}

func (p *Public) Set(key string, val string) {
	pub := *p
	pub[key] = val
	p = &pub
}
