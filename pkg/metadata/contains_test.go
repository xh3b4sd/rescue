package metadata

import (
	"reflect"
	"strconv"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func Test_Metadata_Contains(t *testing.T) {
	testCases := []struct {
		all map[string]string
		sub map[string]string
		con bool
	}{
		// Case 0 ensures empty input results in false.
		{
			all: nil,
			sub: nil,
			con: false,
		},
		// Case 1 ensures empty input results in false.
		{
			all: map[string]string{},
			sub: map[string]string{},
			con: false,
		},
		// Case 2 ensures empty input results in false.
		{
			all: nil,
			sub: map[string]string{},
			con: false,
		},
		// Case 3 ensures empty input results in false.
		{
			all: map[string]string{},
			sub: nil,
			con: false,
		},
		// Case 4 ensures no match results in false.
		{
			all: map[string]string{
				"key": "val",
			},
			sub: map[string]string{},
			con: false,
		},
		// Case 5 ensures no match results in false.
		{
			all: map[string]string{},
			sub: map[string]string{
				"key": "val",
			},
			con: false,
		},
		// Case 6 ensures no match results in false.
		{
			all: map[string]string{
				"key": "val",
			},
			sub: map[string]string{
				"key": "lav",
			},
			con: false,
		},
		// Case 7 ensures missing matches result in false.
		{
			all: map[string]string{
				"key": "val",
				"baz": "zap",
			},
			sub: map[string]string{
				"key": "val",
				"one": "two",
			},
			con: false,
		},
		// Case 8 ensures single matches result in true.
		{
			all: map[string]string{
				"key": "val",
			},
			sub: map[string]string{
				"key": "val",
			},
			con: true,
		},
		// Case 9 ensures multiple matches result in true.
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
			con: true,
		},
	}

	for i, tc := range testCases {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			con := Contains(tc.all, tc.sub)

			if !reflect.DeepEqual(tc.con, con) {
				t.Fatalf("\n\n%s\n", cmp.Diff(tc.con, con))
			}
		})
	}
}
