package matcher

import (
	"strings"
)

// Has expresses whether the given label set contains all of the given subset.
// The first map represents the full label set to match against. The second map
// represents the subset to match for.
func Has(all map[string]string, sub map[string]string) bool {
	if len(sub) == 1 && sub["*"] == "*" {
		return true
	}

	if len(all) == 0 || len(sub) == 0 {
		return false
	}

	for x, y := range sub {
		m, e := has(all, x)
		if !e {
			return false
		}

		if y == "*" {
			continue
		}

		var f bool

		for _, v := range m {
			if y == v {
				f = true
				break
			}
		}

		if !f {
			return false
		}
	}

	return true
}

func has(all map[string]string, k string) (map[string]string, bool) {
	if k == "*" {
		return all, true
	}

	{
		k = strings.TrimPrefix(k, "*")
		k = strings.TrimSuffix(k, "*")
	}

	if len(k) < 3 {
		return nil, false
	}

	m := map[string]string{}
	for x, y := range all {
		if strings.Contains(x, k) {
			m[x] = y
		}
	}

	if len(m) != 0 {
		return m, true
	}

	return nil, false
}
