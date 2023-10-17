package task

type Cron map[string]string

func (c *Cron) All(key ...string) *Cron {
	crn := Cron(All(*c, key...))
	return &crn
}

func (c *Cron) Any(key ...string) *Cron {
	crn := Cron(Any(*c, key...))
	return &crn
}

func (c *Cron) Emp() bool {
	return c.Len() == 0
}

func (c *Cron) Eql(x *Cron) bool {
	return c != nil && x != nil && c.Len() == x.Len() && c.Has(*x)
}

func (c *Cron) Has(lab map[string]string) bool {
	return Has(*c, lab)
}

func (c *Cron) Key() []string {
	if c == nil {
		return nil
	}

	var key []string

	for k := range *c {
		key = append(key, k)
	}

	return key
}

func (c *Cron) Len() int {
	if c == nil {
		return 0
	}

	crn := *c
	return len(crn)
}
