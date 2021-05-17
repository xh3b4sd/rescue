package metadata

func Contains(all map[string]string, sub map[string]string) bool {
	if len(all) == 0 || len(sub) == 0 {
		return false
	}

	for k, v := range sub {
		s, ok := all[k]
		if !ok {
			return false
		}

		if s != v {
			return false
		}
	}

	return true
}
