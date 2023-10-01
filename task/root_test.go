package task

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/google/go-cmp/cmp"
)

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
				Core: &Core{"foo": "bar"},
			},
			emp: true,
		},
		// Case 002
		{
			tas: &Task{
				Meta: &Meta{},
			},
			emp: true,
		},
		// Case 003
		{
			tas: &Task{
				Meta: &Meta{"foo": "bar"},
			},
			emp: true,
		},
		// Case 004
		{
			tas: &Task{
				Root: &Root{},
			},
			emp: true,
		},
		// Case 005
		{
			tas: &Task{
				Root: &Root{"foo": "bar"},
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

func Test_Task_Root_Set(t *testing.T) {
	testCases := []struct {
		tas *Task
		key string
		val string
		set *Task
	}{
		// Case 000
		{
			tas: &Task{
				Root: &Root{"foo": "bar"},
			},
			key: "foo",
			val: "zap",
			set: &Task{
				Root: &Root{"foo": "zap"},
			},
		},
		// Case 001
		{
			tas: &Task{
				Root: &Root{"foo": "bar", "one": "two"},
			},
			key: "one",
			val: "thr",
			set: &Task{
				Root: &Root{"foo": "bar", "one": "thr"},
			},
		},
	}

	for i, tc := range testCases {
		t.Run(fmt.Sprintf("%03d", i), func(t *testing.T) {
			tc.tas.Root.Set(tc.key, tc.val)

			set := tc.set

			if !reflect.DeepEqual(tc.tas, set) {
				t.Fatalf("\n\n%s\n", cmp.Diff(tc.tas, set))
			}
		})
	}
}
