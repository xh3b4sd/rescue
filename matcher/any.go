package matcher

// Any returns a label set that matches any of the given label keys. If any of
// the given label keys does not match, it is simply ignored. That means that
// the returned label set might be nil if not a single label key matches against
// the given label set. If some of the given label keys match, a label set with
// the matching labels is returned.
func Any(all map[string]string, key ...string) map[string]string {
	lab := map[string]string{}

	for _, x := range key {
		m, e := has(all, x)
		if !e {
			continue
		}

		for k, v := range m {
			lab[k] = v
		}
	}

	if len(lab) == 0 {
		return nil
	}

	return lab
}
