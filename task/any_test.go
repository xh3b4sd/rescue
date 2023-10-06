package task

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func Test_Task_Any(t *testing.T) {
	testCases := []struct {
		lab map[string]string
		key []string
		any map[string]string
	}{
		// Case 000
		{
			lab: map[string]string{},
			key: []string{},
			any: nil,
		},
		// Case 001
		{
			lab: map[string]string{
				"test.domain.io/id":      "1",
				"test.domain.io/owner":   "a",
				"test.domain.io/version": "5",
				"some.bucket.io/key":     "foo",
			},
			key: []string{},
			any: nil,
		},
		// Case 002
		{
			lab: map[string]string{
				"test.domain.io/id":      "1",
				"test.domain.io/owner":   "a",
				"test.domain.io/version": "5",
				"some.bucket.io/key":     "foo",
			},
			key: []string{
				"test*",
			},
			any: map[string]string{
				"test.domain.io/id":      "1",
				"test.domain.io/owner":   "a",
				"test.domain.io/version": "5",
			},
		},
		// Case 003
		{
			lab: map[string]string{
				"test.domain.io/id":      "1",
				"test.domain.io/owner":   "a",
				"test.domain.io/version": "5",
				"some.bucket.io/key":     "foo",
			},
			key: []string{
				"this.buc*",
			},
			any: nil,
		},
		// Case 004
		{
			lab: map[string]string{
				"test.domain.io/id":      "1",
				"test.domain.io/owner":   "a",
				"test.domain.io/version": "5",
				"some.bucket.io/key":     "foo",
			},
			key: []string{
				"some.bucket.io",
				"this.buck*",
			},
			any: map[string]string{
				"some.bucket.io/key": "foo",
			},
		},
		// Case 005
		{
			lab: map[string]string{
				"test.domain.io/id":      "1",
				"test.domain.io/owner":   "a",
				"test.domain.io/version": "5",
				"some.bucket.io/key":     "foo",
			},
			key: []string{
				"some.others.io",
				"this.buck*",
			},
			any: nil,
		},
	}

	for i, tc := range testCases {
		t.Run(fmt.Sprintf("%03d", i), func(t *testing.T) {
			any := Any(tc.lab, tc.key...)

			if !reflect.DeepEqual(any, tc.any) {
				t.Fatalf("\n\n%s\n", cmp.Diff(tc.any, any))
			}
		})
	}
}
