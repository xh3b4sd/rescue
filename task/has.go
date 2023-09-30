package task

import "strings"

func (t *Task) Has(met map[string]string) bool {
	if len(met) == 1 && met["*"] == "*" {
		return true
	}

	if len(t.Meta) == 0 || len(met) == 0 {
		return false
	}

	for a, b := range met {
		// TODO unit test, we want to have a key wildcards
		//
		//     ke*: val
		//     *ey: val
		//
		m, e := t.has(a)
		if !e {
			return false
		}

		// TODO unit test, we want to have a value wildcard
		//
		//     key: *
		//
		if b == "*" {
			continue
		}

		var f bool

		for _, y := range m {
			if b == y {
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

func (t *Task) has(a string) (map[string]string, bool) {
	{
		a = strings.TrimPrefix(a, "*")
		a = strings.TrimSuffix(a, "*")
	}

	if len(a) < 3 {
		return nil, false
	}

	m := map[string]string{}
	for x, y := range t.Meta {
		if strings.Contains(x, a) {
			m[x] = y
		}
	}

	if len(m) != 0 {
		return m, true
	}

	return nil, false
}
