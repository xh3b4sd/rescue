package balancer

import (
	"math"
	"sort"
)

const (
	// DeviationThreshold is the fraction of maximum cumulative deviation allowed
	// between current and desired balance.
	DeviationThreshold = 0.20
	// ReductionParameter is the fraction of reduction applied to any identified
	// deviation. Given a deviation of 5, 1 and 3 for workers a, b, and c, the
	// resulting deviations after applying ReductionParameter would be 2, 1 and 1
	// respectively.
	ReductionParameter = 0.50
)

type Balancer struct{}

func New() *Balancer {
	return &Balancer{}
}

func (b *Balancer) Dev(cur map[string]int, des map[string]int) map[string]int {
	var key []string
	{
		key = append(key, keys(cur)...)
		key = append(key, keys(des)...)
	}

	dev := map[string]int{}
	for _, k := range key {
		c := cur[k]
		d := des[k]

		r := (c - d) * 2

		if r > 0 {
			dev[k] = r
		}
	}

	if sum(cur) > 5 && float64(sum(dev)) >= float64(sum(cur))*DeviationThreshold {
		return b.reduce(dev)
	}

	if sum(cur) > 50 && float64(sum(dev)) >= float64(sum(cur))*DeviationThreshold/2 {
		return b.reduce(dev)
	}

	if sum(cur) > 50 && float64(max(dev)) >= float64(sum(cur))*DeviationThreshold/2 {
		return b.reduce(dev)
	}

	return nil
}

func (b *Balancer) Opt(wrk []string, tas int) map[string]int {
	if len(wrk) == 0 {
		return nil
	}

	var cop []string
	for _, x := range wrk {
		if !contains(cop, x) {
			cop = append(cop, x)
		}
	}

	{
		sort.Strings(cop)
	}

	bal := map[string]int{}

	for {
		if tas == 0 {
			break
		}

		for _, x := range cop {
			bal[x]++

			tas--

			if tas == 0 {
				break
			}
		}
	}

	return bal
}

func (b *Balancer) reduce(dev map[string]int) map[string]int {
	red := map[string]int{}

	for k := range dev {
		r := float64(dev[k]) * ReductionParameter
		n, _ := math.Modf(r)

		if r > 0 && n == 0 {
			red[k] = 1
		}
		if r > 0 && n == 1 {
			red[k] = 1
		}
		if r > 0 && n > 1 {
			red[k] = int(n / 2)
		}
	}

	return red
}
