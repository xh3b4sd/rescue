package task

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
