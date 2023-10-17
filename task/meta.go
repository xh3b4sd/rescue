package task

type Meta map[string]string

func (m *Meta) All(key ...string) *Meta {
	met := Meta(All(*m, key...))
	return &met
}

func (m *Meta) Any(key ...string) *Meta {
	met := Meta(Any(*m, key...))
	return &met
}

func (m *Meta) Emp() bool {
	return m.Len() == 0
}

func (m *Meta) Eql(x *Meta) bool {
	return m != nil && x != nil && m.Len() == x.Len() && m.Has(*x)
}

func (m *Meta) Exi(key string) bool {
	if m == nil {
		return false
	}

	met := *m
	return key != "" && met[key] != ""
}

func (m *Meta) Get(key string) string {
	if m == nil {
		return ""
	}

	met := *m
	return met[key]
}

func (m *Meta) Has(lab map[string]string) bool {
	return Has(*m, lab)
}

func (m *Meta) Key() []string {
	if m == nil {
		return nil
	}

	var key []string

	for k := range *m {
		key = append(key, k)
	}

	return key
}

func (m *Meta) Len() int {
	if m == nil {
		return 0
	}

	met := *m
	return len(met)
}

func (m *Meta) Set(key string, val string) {
	met := *m
	met[key] = val
}
