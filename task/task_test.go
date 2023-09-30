package task

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func Test_Task_All(t *testing.T) {
	testCases := []struct {
		tas *Task
		key []string
		all *Task
	}{
		// Case 000
		{
			tas: &Task{},
			key: []string{},
			all: &Task{},
		},
		// Case 001
		{
			tas: &Task{},
			key: []string{
				"foo",
				Object,
			},
			all: nil,
		},
		// Case 002
		{
			tas: &Task{
				Meta: map[string]string{
					Object: "1",
				},
			},
			key: []string{
				"foo",
			},
			all: nil,
		},
		// Case 003
		{
			tas: &Task{
				Meta: map[string]string{
					Object: "1",
				},
			},
			key: []string{
				"foo", // foo not in tas
				Object,
			},
			all: nil,
		},
		// Case 004
		{
			tas: &Task{
				Meta: map[string]string{
					Object: "1",
				},
			},
			key: []string{
				Object,
			},
			all: &Task{
				Meta: map[string]string{
					Object: "1",
				},
			},
		},
		// Case 005
		{
			tas: &Task{
				Meta: map[string]string{
					Cycles: "5",
					Object: "1",
					Worker: "a",
				},
			},
			key: []string{
				Object,
			},
			all: &Task{
				Meta: map[string]string{
					Object: "1",
				},
			},
		},
		// Case 006
		{
			tas: &Task{
				Meta: map[string]string{
					Cycles: "5",
					Object: "1",
					Worker: "a",
				},
			},
			key: []string{
				Object,
				Worker,
			},
			all: &Task{
				Meta: map[string]string{
					Object: "1",
					Worker: "a",
				},
			},
		},
	}

	for i, tc := range testCases {
		t.Run(fmt.Sprintf("%03d", i), func(t *testing.T) {
			all := tc.tas.All(tc.key...)

			if !reflect.DeepEqual(tc.all, all) {
				t.Fatalf("\n\n%s\n", cmp.Diff(tc.all, all))
			}
		})
	}
}

func Test_Task_Any(t *testing.T) {
	testCases := []struct {
		tas *Task
		key []string
		any *Task
	}{
		// Case 000
		{
			tas: &Task{},
			key: []string{},
			any: nil,
		},
		// Case 001
		{
			tas: &Task{
				Meta: map[string]string{
					"test.domain.io/id":      "1",
					"test.domain.io/owner":   "a",
					"test.domain.io/version": "5",
					"some.bucket.io/key":     "foo",
				},
			},
			key: []string{},
			any: nil,
		},
		// Case 002
		{
			tas: &Task{
				Meta: map[string]string{
					"test.domain.io/id":      "1",
					"test.domain.io/owner":   "a",
					"test.domain.io/version": "5",
					"some.bucket.io/key":     "foo",
				},
			},
			key: []string{
				"test*",
			},
			any: &Task{
				Meta: map[string]string{
					"test.domain.io/id":      "1",
					"test.domain.io/owner":   "a",
					"test.domain.io/version": "5",
				},
			},
		},
		// Case 003
		{
			tas: &Task{
				Meta: map[string]string{
					"test.domain.io/id":      "1",
					"test.domain.io/owner":   "a",
					"test.domain.io/version": "5",
					"some.bucket.io/key":     "foo",
				},
			},
			key: []string{
				"this.buc*",
			},
			any: nil,
		},
		// Case 004
		{
			tas: &Task{
				Meta: map[string]string{
					"test.domain.io/id":      "1",
					"test.domain.io/owner":   "a",
					"test.domain.io/version": "5",
					"some.bucket.io/key":     "foo",
				},
			},
			key: []string{
				"some.bucket.io",
				"this.buck*",
			},
			any: &Task{
				Meta: map[string]string{
					"some.bucket.io/key": "foo",
				},
			},
		},
		// Case 005
		{
			tas: &Task{
				Meta: map[string]string{
					"test.domain.io/id":      "1",
					"test.domain.io/owner":   "a",
					"test.domain.io/version": "5",
					"some.bucket.io/key":     "foo",
				},
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
			any := tc.tas.Any(tc.key...)

			if !reflect.DeepEqual(tc.any, any) {
				t.Fatalf("\n\n%s\n", cmp.Diff(tc.any, any))
			}
		})
	}
}

func Test_Task_Has(t *testing.T) {
	testCases := []struct {
		all map[string]string
		sub map[string]string
		has bool
	}{
		// Case 000 ensures empty input results in false.
		{
			all: nil,
			sub: nil,
			has: false,
		},
		// Case 001 ensures empty input results in false.
		{
			all: map[string]string{},
			sub: map[string]string{},
			has: false,
		},
		// Case 002 ensures empty input results in false.
		{
			all: nil,
			sub: map[string]string{},
			has: false,
		},
		// Case 003 ensures empty input results in false.
		{
			all: map[string]string{},
			sub: nil,
			has: false,
		},
		// Case 004 ensures no match results in false.
		{
			all: map[string]string{
				"key": "val",
			},
			sub: map[string]string{},
			has: false,
		},
		// Case 005 ensures no match results in false.
		{
			all: map[string]string{},
			sub: map[string]string{
				"key": "val",
			},
			has: false,
		},
		// Case 006 ensures no match results in false.
		{
			all: map[string]string{
				"key": "val",
			},
			sub: map[string]string{
				"key": "lav",
			},
			has: false,
		},
		// Case 007 ensures missing matches result in false.
		{
			all: map[string]string{
				"key": "val",
				"baz": "zap",
			},
			sub: map[string]string{
				"key": "val",
				"one": "two",
			},
			has: false,
		},
		// Case 008 ensures single matches result in true.
		{
			all: map[string]string{
				"key": "val",
			},
			sub: map[string]string{
				"key": "val",
			},
			has: true,
		},
		// Case 009 ensures multiple matches result in true.
		{
			all: map[string]string{
				"key": "val",
				"one": "two",
				"baz": "zap",
			},
			sub: map[string]string{
				"key": "val",
				"one": "two",
			},
			has: true,
		},
		// Case 010 ensures that the catch all returns true for any metadat.
		{
			all: map[string]string{
				"key": "val",
				"one": "val",
				"baz": "val",
			},
			sub: map[string]string{
				"*": "*",
			},
			has: true,
		},
		// Case 011 ensures that the catch all returns true for empty metadat.
		{
			all: map[string]string{},
			sub: map[string]string{
				"*": "*",
			},
			has: true,
		},
	}

	for i, tc := range testCases {
		t.Run(fmt.Sprintf("%03d", i), func(t *testing.T) {
			var tas *Task
			{
				tas = &Task{
					Meta: tc.all,
				}
			}

			has := tas.Has(tc.sub)

			if !reflect.DeepEqual(tc.has, has) {
				t.Fatalf("\n\n%s\n", cmp.Diff(tc.has, has))
			}
		})
	}
}
