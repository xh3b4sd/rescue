package random

import (
	"crypto/rand"
	"math/big"
)

const (
	Length = 5
	Number = "0123456789"
	Letter = "abcdefghijklmnopqrstuvwxyz"
)

const (
	all = Number + Letter
)

// New returns a new random string that can be used to generate a worker name.
// Worker names are reflected in task metadata using the owner label.
func New() string {
	var ran string

	for i := 0; i < Length; i++ {
		ind, err := rand.Int(rand.Reader, big.NewInt(int64(len(all))))
		if err != nil {
			panic(err)
		}

		ran += string(all[ind.Int64()])
	}

	return ran
}
