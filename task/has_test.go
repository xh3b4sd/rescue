package task

import (
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
)

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
		// Case 012
		{
			all: map[string]string{
				"key": "val",
				"one": "two",
				"baz": "zap",
			},
			sub: map[string]string{
				"key": "*",
				"one": "two",
			},
			has: true,
		},
		// Case 013
		{
			all: map[string]string{
				"key": "val",
				"one": "two",
				"baz": "zap",
			},
			sub: map[string]string{
				"key": "*",
				"one": "*",
			},
			has: true,
		},
		// Case 014
		{
			all: map[string]string{
				"key": "val",
				"one": "two",
				"baz": "zap",
			},
			sub: map[string]string{
				"foo": "*",
				"one": "two",
			},
			has: false,
		},
		// Case 015
		{
			all: map[string]string{
				"some.key.io": "val",
				"some.one.io": "two",
				"some.baz.io": "zap",
			},
			sub: map[string]string{
				"*key.io": "val",
			},
			has: true,
		},
		// Case 016
		{
			all: map[string]string{
				"some.key.io": "val",
				"some.one.io": "two",
				"some.baz.io": "zap",
			},
			sub: map[string]string{
				"some*":       "val",
				"some.one.io": "two",
			},
			has: true,
		},
		// Case 017
		{
			all: map[string]string{
				"some.key.io": "val",
				"some.one.io": "two",
				"some.baz.io": "zap",
			},
			sub: map[string]string{
				"some*":       "val",
				"some.one.io": "wng",
			},
			has: false,
		},
		// Case 018
		{
			all: map[string]string{
				"test.api.io/num": "3",
			},
			sub: map[string]string{
				"*task.rescue.io*": "*",
			},
			has: false,
		},
		// Case 019
		{
			all: map[string]string{
				"test.rescue.io/num": "3",
			},
			sub: map[string]string{
				"*rescue.io*": "*",
			},
			has: true,
		},
	}

	for i, tc := range testCases {
		t.Run(fmt.Sprintf("%03d", i), func(t *testing.T) {
			has := Has(tc.all, tc.sub)

			if has != tc.has {
				t.Fatalf("\n\n%s\n", cmp.Diff(tc.has, has))
			}
		})
	}
}
