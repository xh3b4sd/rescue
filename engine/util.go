package engine

func contains(lis []string, itm string) bool {
	for _, l := range lis {
		if l == itm {
			return true
		}
	}

	return false
}

func ensure(lis []string, str ...string) []string {
	return append(remove(lis, str...), str...)
}

func keys(lis map[string]int) []string {
	var key []string

	for k := range lis {
		if k == "" {
			continue
		}

		key = append(key, k)
	}

	return key
}

func remove(lis []string, str ...string) []string {
	var res []string

	for _, l := range lis {
		if contains(str, l) {
			continue
		}

		res = append(res, l)
	}

	return res
}

func sum(lis map[string]int) int {
	if len(lis) == 0 {
		return 0
	}

	var sum int

	for _, l := range lis {
		sum += l
	}

	return sum
}
