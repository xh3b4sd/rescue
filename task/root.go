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
	return r.Len() == 0
}

func (r *Root) Exi(key string) bool {
	if r == nil {
		return false
	}

	roo := *r
	return key != "" && roo[key] != ""
}

func (r *Root) Get(key string) string {
	if r == nil {
		return ""
	}

	roo := *r
	return roo[key]
}

func (r *Root) Has(lab map[string]string) bool {
	return Has(*r, lab)
}

func (r *Root) Len() int {
	if r == nil {
		return 0
	}

	roo := *r
	return len(roo)
}

func (r *Root) Set(key string, val string) {
	roo := *r
	roo[key] = val
}
