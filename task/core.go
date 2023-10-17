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

func (c *Core) Eql(x *Core) bool {
	return c != nil && x != nil && c.Len() == x.Len() && c.Has(*x)
}

func (c *Core) Has(lab map[string]string) bool {
	return Has(*c, lab)
}

func (c *Core) Key() []string {
	if c == nil {
		return nil
	}

	var key []string

	for k := range *c {
		key = append(key, k)
	}

	return key
}

func (c *Core) Len() int {
	if c == nil {
		return 0
	}

	cor := *c
	return len(cor)
}
