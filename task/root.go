package task

type Root map[string]string

func (r *Root) All(key ...string) *Root {
	roo := Root(All(*r, key...))
	return &roo
}

func (r *Root) Any(key ...string) *Root {
	roo := Root(Any(*r, key...))
	return &roo
}

func (r *Root) Emp() bool {
	if r == nil {
		return true
	}

	roo := *r
	return len(roo) == 0
}

func (r *Root) Get(key string) string {
	roo := *r
	return roo[key]
}

func (r *Root) Has(lab map[string]string) bool {
	return Has(*r, lab)
}

func (r *Root) Set(key string, val string) {
	roo := *r
	roo[key] = val
}
