package task

type Gate map[string]string

func (g *Gate) All(key ...string) *Gate {
	gat := Gate(All(*g, key...))
	return &gat
}

func (g *Gate) Any(key ...string) *Gate {
	gat := Gate(Any(*g, key...))
	return &gat
}

func (g *Gate) Emp() bool {
	return g.Len() == 0
}

func (g *Gate) Eql(x *Gate) bool {
	return g != nil && x != nil && g.Len() == x.Len() && g.Has(*x)
}

func (g *Gate) Exi(key string) bool {
	if g == nil {
		return false
	}

	gat := *g
	return key != "" && gat[key] != ""
}

func (g *Gate) Get(key string) string {
	if g == nil {
		return ""
	}

	gat := *g
	return gat[key]
}

func (g *Gate) Has(lab map[string]string) bool {
	return Has(*g, lab)
}

func (g *Gate) Key() []string {
	if g == nil {
		return nil
	}

	var key []string

	for k := range *g {
		key = append(key, k)
	}

	return key
}

func (g *Gate) Len() int {
	if g == nil {
		return 0
	}

	gat := *g
	return len(gat)
}

func (g *Gate) Set(key string, val string) {
	gat := *g
	gat[key] = val
}
