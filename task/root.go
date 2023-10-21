package task

import "github.com/xh3b4sd/rescue/matcher"

type Root map[string]string

func (r *Root) All(key ...string) *Root {
	roo := Root(matcher.All(*r, key...))
	return &roo
}

func (r *Root) Any(key ...string) *Root {
	roo := Root(matcher.Any(*r, key...))
	return &roo
}

func (r *Root) Emp() bool {
	return r.Len() == 0
}

func (r *Root) Eql(x *Root) bool {
	return r != nil && x != nil && r.Len() == x.Len() && r.Has(*x)
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
	return matcher.Has(*r, lab)
}

func (r *Root) Key() []string {
	if r == nil {
		return nil
	}

	var key []string

	for k := range *r {
		key = append(key, k)
	}

	return key
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
