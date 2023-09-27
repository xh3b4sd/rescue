package random

import (
	"crypto/rand"

	"github.com/xh3b4sd/budget/v3"
	"github.com/xh3b4sd/budget/v3/pkg/breaker"
	"github.com/xh3b4sd/random"
)

const (
	Len = 5
)

const (
	Digits = "0123456789"
	Lower  = "abcdefghijklmnopqrstuvwxyz"
)

const (
	chars = Digits + Lower
)

func MustNew() string {
	var err error

	var bre budget.Interface
	{
		c := breaker.Config{
			Limiter: breaker.Limiter{
				Budget: Len,
			},
		}

		bre, err = breaker.New(c)
		if err != nil {
			panic(err)
		}
	}

	var ran random.Interface
	{
		c := random.Config{
			Budget:     bre,
			RandFunc:   rand.Int,
			RandReader: rand.Reader,
		}

		ran, err = random.New(c)
		if err != nil {
			panic(err)
		}
	}

	var pas string
	{
		lis, err := ran.NMax(Len, len(chars))
		if err != nil {
			panic(err)
		}

		for _, i := range lis {
			pas += string(chars[i])
		}
	}

	return pas
}
