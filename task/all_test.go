package task

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
				Object,
			},
			all: nil,
		},
		// Case 002
		{
			lab: map[string]string{
				Object: "1",
			},
			key: []string{
				"foo",
			},
			all: nil,
		},
		// Case 003
		{
			lab: map[string]string{
				Object: "1",
			},
			key: []string{
				"foo", // foo not in lab
				Object,
			},
			all: nil,
		},
		// Case 004
		{
			lab: map[string]string{
				Object: "1",
			},
			key: []string{
				Object,
			},
			all: map[string]string{
				Object: "1",
			},
		},
		// Case 005
		{
			lab: map[string]string{
				Cycles: "5",
				Object: "1",
				Worker: "a",
			},
			key: []string{
				Object,
			},
			all: map[string]string{
				Object: "1",
			},
		},
		// Case 006
		{
			lab: map[string]string{
				Cycles: "5",
				Object: "1",
				Worker: "a",
			},
			key: []string{
				Object,
				Worker,
			},
			all: map[string]string{
				Object: "1",
				Worker: "a",
			},
		},
	}

	for i, tc := range testCases {
		t.Run(fmt.Sprintf("%03d", i), func(t *testing.T) {
			all := All(tc.lab, tc.key...)

			if !reflect.DeepEqual(tc.all, all) {
				t.Fatalf("\n\n%s\n", cmp.Diff(tc.all, all))
			}
		})
	}
}
