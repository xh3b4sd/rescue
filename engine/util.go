package engine

import (
	"time"

	"github.com/xh3b4sd/objectid"
)

func contains(l []string, i string) bool {
	for _, x := range l {
		if x == i {
			return true
		}
	}

	return false
}

func ensure(l []string, s ...string) []string {
	return append(remove(l, s...), s...)
}

// expiry returns the earliest recorded expiry of the provided local cache.
func expiry(m map[objectid.ID]*local) time.Time {
	var fir time.Time

	for _, v := range m {
		if !v.exp.IsZero() {
			if fir.IsZero() || v.exp.Before(fir) {
				fir = v.exp
			}
		}
	}

	return fir
}

func keys(m map[string]int) []string {
	var key []string

	for k := range m {
		if k != "" {
			key = append(key, k)
		}
	}

	return key
}

func remove(l []string, s ...string) []string {
	var res []string

	for _, x := range l {
		if !contains(s, x) {
			res = append(res, x)
		}
	}

	return res
}

func sum(m map[string]int) int {
	if len(m) == 0 {
		return 0
	}

	var sum int

	for _, l := range m {
		sum += l
	}

	return sum
}
