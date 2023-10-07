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
	return c.Len() == 0
}

func (c *Core) Has(lab map[string]string) bool {
	return Has(*c, lab)
}

func (c *Core) Len() int {
	if c == nil {
		return 0
	}

	cor := *c
	return len(cor)
}
