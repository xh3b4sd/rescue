package task

type Core map[string]string

func (c *Core) All(key ...string) *Core {
	cor := Core(All(*c, key...))
	return &cor
}

func (c *Core) Any(key ...string) *Core {
	cor := Core(Any(*c, key...))
	return &cor
}

func (c *Core) Emp() bool {
	if c == nil {
		return true
	}

	cor := *c
	return len(cor) == 0
}

func (c *Core) Has(lab map[string]string) bool {
	return Has(*c, lab)
}
