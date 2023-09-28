package random

import (
	"sync"
	"testing"
)

func Test_Random_New(t *testing.T) {
	mut := sync.Mutex{}
	see := map[string]struct{}{}

	for i := 0; i < 1000; i++ {
		go func() {
			a := New()
			b := New()

			if a == "" || b == "" {
				panic("random strings must not be empty")
			}

			if len(a) != Length || len(b) != Length {
				panic("random strings must have the right length")
			}

			if a == b {
				panic("random strings must be unique")
			}

			{
				mut.Lock()
			}

			{
				_, exi := see[a]
				if exi {
					panic("random strings must not be duplicated")
				}

				see[a] = struct{}{}
			}

			{
				_, exi := see[b]
				if exi {
					panic("random strings must not be duplicated")
				}

				see[b] = struct{}{}
			}

			{
				mut.Unlock()
			}
		}()
	}
}
