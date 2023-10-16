package task

type Sync map[string]string

func (s *Sync) All(key ...string) *Sync {
	syn := Sync(All(*s, key...))
	return &syn
}

func (s *Sync) Any(key ...string) *Sync {
	syn := Sync(Any(*s, key...))
	return &syn
}

func (s *Sync) Emp() bool {
	return s.Len() == 0
}

func (s *Sync) Exi(key string) bool {
	if s == nil {
		return false
	}

	syn := *s
	return key != "" && syn[key] != ""
}

func (s *Sync) Get(key string) string {
	if s == nil {
		return ""
	}

	syn := *s
	return syn[key]
}

func (s *Sync) Has(lab map[string]string) bool {
	return Has(*s, lab)
}

func (s *Sync) Key() []string {
	if s == nil {
		return nil
	}

	var key []string

	for k := range *s {
		key = append(key, k)
	}

	return key
}

func (s *Sync) Len() int {
	if s == nil {
		return 0
	}

	syn := *s
	return len(syn)
}

func (s *Sync) Set(key string, val string) {
	syn := *s
	syn[key] = val
}
