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

func Test_Task_Root_Exi(t *testing.T) {
	testCases := []struct {
		tas *Task
		exi bool
	}{
		// Case 000
		{
			tas: &Task{},
			exi: false,
		},
		// Case 001
		{
			tas: &Task{
				Core: &Core{"foo": "bar"},
			},
			exi: false,
		},
		// Case 002
		{
			tas: &Task{
				Meta: &Meta{},
			},
			exi: false,
		},
		// Case 003
		{
			tas: &Task{
				Meta: &Meta{"foo": "bar"},
			},
			exi: false,
		},
		// Case 004
		{
			tas: &Task{
				Root: &Root{},
			},
			exi: false,
		},
		// Case 005
		{
			tas: &Task{
				Root: &Root{"foo": "bar"},
			},
			exi: true,
		},
	}

	for i, tc := range testCases {
		t.Run(fmt.Sprintf("%03d", i), func(t *testing.T) {
			exi := tc.tas.Root.Exi("foo")

			if !reflect.DeepEqual(tc.exi, exi) {
				t.Fatalf("\n\n%s\n", cmp.Diff(tc.exi, exi))
			}
		})
	}
}

func Test_Task_Root_Get(t *testing.T) {
	testCases := []struct {
		tas *Task
		get string
	}{
		// Case 000
		{
			tas: &Task{},
			get: "",
		},
		// Case 001
		{
			tas: &Task{
				Core: &Core{"foo": "bar"},
			},
			get: "",
		},
		// Case 002
		{
			tas: &Task{
				Meta: &Meta{},
			},
			get: "",
		},
		// Case 003
		{
			tas: &Task{
				Meta: &Meta{"foo": "bar"},
			},
			get: "",
		},
		// Case 004
		{
			tas: &Task{
				Root: &Root{},
			},
			get: "",
		},
		// Case 005
		{
			tas: &Task{
				Root: &Root{"foo": "bar"},
			},
			get: "bar",
		},
	}

	for i, tc := range testCases {
		t.Run(fmt.Sprintf("%03d", i), func(t *testing.T) {
			get := tc.tas.Root.Get("foo")

			if !reflect.DeepEqual(tc.get, get) {
				t.Fatalf("\n\n%s\n", cmp.Diff(tc.get, get))
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

			if !reflect.DeepEqual(set, tc.tas) {
				t.Fatalf("\n\n%s\n", cmp.Diff(tc.tas, set))
			}
		})
	}
}
