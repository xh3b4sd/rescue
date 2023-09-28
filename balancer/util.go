package balancer

func contains(lis []string, itm string) bool {
	for _, l := range lis {
		if l == itm {
			return true
		}
	}

	return false
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

func max(lis map[string]int) int {
	var max int

	for _, l := range lis {
		if l > max {
			max = l
		}
	}

	return max
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
