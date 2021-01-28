package random

import (
	"crypto/rand"
	"time"

	"github.com/xh3b4sd/budget"
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

	var b budget.Interface
	{
		c := budget.ConstantConfig{
			Budget:   3,
			Duration: 1 * time.Second,
		}

		b, err = budget.NewConstant(c)
		if err != nil {
			panic(err)
		}
	}

	var r random.Interface
	{
		c := random.Config{
			Budget:     b,
			RandFunc:   rand.Int,
			RandReader: rand.Reader,

			Timeout: 1 * time.Second,
		}

		r, err = random.New(c)
		if err != nil {
			panic(err)
		}
	}

	l, err := r.NMax(Len, len(chars))
	if err != nil {
		panic(err)
	}

	var pass string

	for _, i := range l {
		pass += string(chars[i])
	}

	return pass
}
