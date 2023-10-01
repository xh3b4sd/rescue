package task

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
