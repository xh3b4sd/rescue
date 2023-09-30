package balancer

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func Test_Balancer_Dev(t *testing.T) {
	testCases := []struct {
		cur map[string]int
		des map[string]int
		dev map[string]int
	}{
		// Case 000
		{
			cur: map[string]int{},
			des: map[string]int{},
			dev: nil,
		},
		// Case 001
		{
			cur: map[string]int{
				"a": 1,
				"b": 1,
			},
			des: map[string]int{
				"a": 1,
				"b": 1,
			},
			dev: nil,
		},
		// Case 002
		{
			cur: map[string]int{
				"a": 9,
				"b": 8,
				"c": 8,
			},
			des: map[string]int{
				"a": 9,
				"b": 8,
				"c": 8,
			},
			dev: nil,
		},
		// Case 003
		{
			cur: map[string]int{
				"a": 12,
				"b": 5,
				"c": 8,
			},
			des: map[string]int{
				"a": 9,
				"b": 8,
				"c": 8,
			},
			dev: map[string]int{
				"a": 1,
			},
		},
		// Case 004
		{
			cur: map[string]int{
				"a": 12,
				"b": 3,
				"c": 10,
			},
			des: map[string]int{
				"a": 9,
				"b": 8,
				"c": 8,
			},
			dev: map[string]int{
				"a": 1,
				"c": 1,
			},
		},
		// Case 005
		{
			cur: map[string]int{
				"a": 86,
				"b": 75,
				"c": 44,
			},
			des: map[string]int{
				"a": 69,
				"b": 68,
				"c": 68,
			},
			dev: map[string]int{
				"a": 8,
				"b": 3,
			},
		},
		// Case 006
		{
			cur: map[string]int{
				"a": 10,
				"b": 10,
			},
			des: map[string]int{
				"a": 7,
				"b": 7,
				"c": 6,
			},
			dev: map[string]int{
				"a": 1,
				"b": 1,
			},
		},
		// Case 007
		{
			cur: map[string]int{
				"a": 9,
				"b": 9,
				"c": 2,
			},
			des: map[string]int{
				"a": 7,
				"b": 7,
				"c": 6,
			},
			dev: map[string]int{
				"a": 1,
				"b": 1,
			},
		},
		// Case 008
		{
			cur: map[string]int{
				"a": 8,
				"b": 8,
				"c": 4,
			},
			des: map[string]int{
				"a": 7,
				"b": 7,
				"c": 6,
			},
			dev: map[string]int{
				"a": 1,
				"b": 1,
			},
		},
		// Case 009
		{
			cur: map[string]int{
				"a": 7,
				"b": 7,
				"c": 6,
			},
			des: map[string]int{
				"a": 7,
				"b": 7,
				"c": 6,
			},
			dev: nil,
		},
	}

	for i, tc := range testCases {
		t.Run(fmt.Sprintf("%03d", i), func(t *testing.T) {
			dev := Default().Dev(tc.cur, tc.des)

			if !reflect.DeepEqual(tc.dev, dev) {
				t.Fatalf("\n\n%s\n", cmp.Diff(tc.dev, dev))
			}
		})
	}
}

func Test_Balancer_Opt(t *testing.T) {
	testCases := []struct {
		own []string
		tas int
		opt map[string]int
	}{
		// Case 000
		{
			own: []string{},
			tas: 0,
			opt: nil,
		},
		// Case 001
		{
			own: []string{},
			tas: 3,
			opt: nil,
		},
		// Case 002
		{
			own: []string{
				"a",
			},
			tas: 5,
			opt: map[string]int{
				"a": 5,
			},
		},
		// Case 003
		{
			own: []string{
				"a",
				"b",
			},
			tas: 6,
			opt: map[string]int{
				"a": 3,
				"b": 3,
			},
		},
		// Case 004
		{
			own: []string{
				"a",
				"b",
			},
			tas: 7,
			opt: map[string]int{
				"a": 4,
				"b": 3,
			},
		},
		// Case 005
		{
			own: []string{
				"a",
				"b",
				"c",
			},
			tas: 7,
			opt: map[string]int{
				"a": 3,
				"b": 2,
				"c": 2,
			},
		},
		// Case 006
		{
			own: []string{
				"a",
				"b",
				"b",
				"b",
				"c",
			},
			tas: 7,
			opt: map[string]int{
				"a": 3,
				"b": 2,
				"c": 2,
			},
		},
		// Case 007
		{
			own: []string{
				"a",
				"b",
				"c",
				"d",
				"e",
			},
			tas: 13,
			opt: map[string]int{
				"a": 3,
				"b": 3,
				"c": 3,
				"d": 2,
				"e": 2,
			},
		},
		// Case 008
		{
			own: []string{
				"a",
				"b",
				"c",
			},
			tas: 20,
			opt: map[string]int{
				"a": 7,
				"b": 7,
				"c": 6,
			},
		},
	}

	for i, tc := range testCases {
		t.Run(fmt.Sprintf("%03d", i), func(t *testing.T) {
			opt := Default().Opt(tc.own, tc.tas)

			if !reflect.DeepEqual(tc.opt, opt) {
				t.Fatalf("\n\n%s\n", cmp.Diff(tc.opt, opt))
			}
		})
	}
}
