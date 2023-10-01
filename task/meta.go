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
	if m == nil {
		return true
	}

	met := *m
	return len(met) == 0
}

func (m *Meta) Get(key string) string {
	met := *m
	return met[key]
}

func (m *Meta) Has(lab map[string]string) bool {
	return Has(*m, lab)
}

func (m *Meta) Set(key string, val string) {
	met := *m
	met[key] = val
}
