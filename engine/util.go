package engine

import (
	"slices"
	"time"
)

func contains(l []string, i string) bool {
	for _, x := range l {
		if x == i {
			return true
		}
	}

	return false
}

// created receives a task ID of nanoseconds. Nanoseconds in Go are nullified,
// which leaves us with microsecond precision. Internally, created divides by
// 1000, to call time.UnixMicro with the result.
func created(i int64) time.Time {
	//
	//     seconds    1697575809
	//     milli      1697575809 629
	//     micro      1697575809 629 451
	//     nano       1697575809 629 451 000
	//
	return time.UnixMicro(i / 1000).UTC()
}

func ensure(l []string, s ...string) []string {
	return append(remove(l, s...), s...)
}

func first(l []int64) time.Time {
	slices.Sort(l)
	return created(l[0])
}

func unix(l []time.Time) []int64 {
	var uni []int64

	for _, x := range l {
		uni = append(uni, x.UnixNano())
	}

	return uni
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

func values(m map[string]*local) []time.Time {
	var val []time.Time

	for _, v := range m {
		if !v.exp.IsZero() {
			val = append(val, v.exp)
		}
	}

	return val
}
