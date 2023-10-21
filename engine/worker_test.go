package engine

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/xh3b4sd/redigo"
)

func Test_Engine_Worker(t *testing.T) {
	testCases := []struct {
		wrk string
	}{
		// Case 000
		{
			wrk: "foo",
		},
		// Case 001
		{
			wrk: "bar",
		},
		// Case 002
		{
			wrk: "123",
		},
	}

	for i, tc := range testCases {
		t.Run(fmt.Sprintf("%03d", i), func(t *testing.T) {
			var e *Engine
			{
				e = New(Config{
					Redigo: redigo.Fake(),
					Worker: tc.wrk,
				})
			}

			wrk := e.Worker()

			if !reflect.DeepEqual(wrk, tc.wrk) {
				t.Fatalf("\n\n%s\n", cmp.Diff(tc.wrk, wrk))
			}
		})
	}
}
