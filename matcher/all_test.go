package matcher

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func Test_Task_All(t *testing.T) {
	testCases := []struct {
		lab map[string]string
		key []string
		all map[string]string
	}{
		// Case 000
		{
			lab: map[string]string{},
			key: []string{},
			all: map[string]string{},
		},
		// Case 001
		{
			lab: map[string]string{},
			key: []string{
				"foo",
				"test.api.io/object",
			},
			all: nil,
		},
		// Case 002
		{
			lab: map[string]string{
				"test.api.io/object": "1",
			},
			key: []string{
				"foo",
			},
			all: nil,
		},
		// Case 003
		{
			lab: map[string]string{
				"test.api.io/object": "1",
			},
			key: []string{
				"foo", // foo not in lab
				"test.api.io/object",
			},
			all: nil,
		},
		// Case 004
		{
			lab: map[string]string{
				"test.api.io/object": "1",
			},
			key: []string{
				"test.api.io/object",
			},
			all: map[string]string{
				"test.api.io/object": "1",
			},
		},
		// Case 005
		{
			lab: map[string]string{
				"test.api.io/cycles": "5",
				"test.api.io/object": "1",
				"test.api.io/worker": "a",
			},
			key: []string{
				"test.api.io/object",
			},
			all: map[string]string{
				"test.api.io/object": "1",
			},
		},
		// Case 006
		{
			lab: map[string]string{
				"test.api.io/cycles": "5",
				"test.api.io/object": "1",
				"test.api.io/worker": "a",
			},
			key: []string{
				"test.api.io/object",
				"test.api.io/worker",
			},
			all: map[string]string{
				"test.api.io/object": "1",
				"test.api.io/worker": "a",
			},
		},
		// Case 007
		{
			lab: map[string]string{
				"foo":                  "bar",
				"test.api.io/cycles":   "5",
				"test.other.io/object": "1",
				"test.api.io/worker":   "a",
			},
			key: []string{
				"*.api.io/*",
			},
			all: map[string]string{
				"test.api.io/cycles": "5",
				"test.api.io/worker": "a",
			},
		},
	}

	for i, tc := range testCases {
		t.Run(fmt.Sprintf("%03d", i), func(t *testing.T) {
			all := All(tc.lab, tc.key...)

			if !reflect.DeepEqual(all, tc.all) {
				t.Fatalf("\n\n%s\n", cmp.Diff(tc.all, all))
			}
		})
	}
}
