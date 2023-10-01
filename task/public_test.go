package task

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func Test_Task_Meta_Emp(t *testing.T) {
	testCases := []struct {
		tas *Task
		emp bool
	}{
		// Case 000
		{
			tas: &Task{},
			emp: true,
		},
		// Case 001
		{
			tas: &Task{
				Core: &Intern{"foo": "bar"},
			},
			emp: true,
		},
		// Case 002
		{
			tas: &Task{
				Meta: &Public{},
			},
			emp: true,
		},
		// Case 003
		{
			tas: &Task{
				Meta: &Public{"foo": "bar"},
			},
			emp: false,
		},
		// Case 004
		{
			tas: &Task{
				Root: &Public{},
			},
			emp: true,
		},
		// Case 005
		{
			tas: &Task{
				Root: &Public{"foo": "bar"},
			},
			emp: true,
		},
	}

	for i, tc := range testCases {
		t.Run(fmt.Sprintf("%03d", i), func(t *testing.T) {
			emp := tc.tas.Meta.Emp()

			if !reflect.DeepEqual(tc.emp, emp) {
				t.Fatalf("\n\n%s\n", cmp.Diff(tc.emp, emp))
			}
		})
	}
}

func Test_Task_Root_Emp(t *testing.T) {
	testCases := []struct {
		tas *Task
		emp bool
	}{
		// Case 000
		{
			tas: &Task{},
			emp: true,
		},
		// Case 001
		{
			tas: &Task{
				Core: &Intern{"foo": "bar"},
			},
			emp: true,
		},
		// Case 002
		{
			tas: &Task{
				Meta: &Public{},
			},
			emp: true,
		},
		// Case 003
		{
			tas: &Task{
				Meta: &Public{"foo": "bar"},
			},
			emp: true,
		},
		// Case 004
		{
			tas: &Task{
				Root: &Public{},
			},
			emp: true,
		},
		// Case 005
		{
			tas: &Task{
				Root: &Public{"foo": "bar"},
			},
			emp: false,
		},
	}

	for i, tc := range testCases {
		t.Run(fmt.Sprintf("%03d", i), func(t *testing.T) {
			emp := tc.tas.Root.Emp()

			if !reflect.DeepEqual(tc.emp, emp) {
				t.Fatalf("\n\n%s\n", cmp.Diff(tc.emp, emp))
			}
		})
	}
}
