package matcher

// All returns a label set that matches all the given label keys. If any of the
// given label keys does not match, nil is returned. That means that the
// returned label set will be nil, unless the complete list of the given label
// keys matches against the given label set.
func All(all map[string]string, key ...string) map[string]string {
	lab := map[string]string{}

	for _, x := range key {
		m, e := has(all, x)
		if !e {
			return nil
		}

		for k, v := range m {
			lab[k] = v
		}
	}

	return lab
}
