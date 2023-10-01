package task

type Intern map[string]string

func (i *Intern) All(key ...string) *Intern {
	val := Intern(All(*i, key...))
	return &val
}

func (i *Intern) Any(key ...string) *Intern {
	val := Intern(Any(*i, key...))
	return &val
}

func (i *Intern) Emp() bool {
	if i == nil {
		return true
	}

	val := *i
	return len(val) == 0
}

func (i *Intern) Has(lab map[string]string) bool {
	return Has(*i, lab)
}
