package random

import (
	"testing"
)

func Test_Random_MustNew(t *testing.T) {
	for i := 0; i < 1000; i++ {
		go func() {
			a := MustNew()
			b := MustNew()

			if a == b {
				panic("random strings must be unique")
			}
		}()
	}
}
