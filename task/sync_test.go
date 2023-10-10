package task

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func Test_Task_Sync_Emp(t *testing.T) {
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
				Sync: &Sync{},
			},
			emp: true,
		},
		// Case 003
		{
			tas: &Task{
				Sync: &Sync{"foo": "bar"},
			},
			emp: false,
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
			emp: true,
		},
	}

	for i, tc := range testCases {
		t.Run(fmt.Sprintf("%03d", i), func(t *testing.T) {
			emp := tc.tas.Sync.Emp()

			if emp != tc.emp {
				t.Fatalf("\n\n%s\n", cmp.Diff(tc.emp, emp))
			}
		})
	}
}

func Test_Task_Sync_Exi(t *testing.T) {
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
				Sync: &Sync{},
			},
			exi: false,
		},
		// Case 003
		{
			tas: &Task{
				Sync: &Sync{"foo": "bar"},
			},
			exi: true,
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
			exi: false,
		},
	}

	for i, tc := range testCases {
		t.Run(fmt.Sprintf("%03d", i), func(t *testing.T) {
			exi := tc.tas.Sync.Exi("foo")

			if exi != tc.exi {
				t.Fatalf("\n\n%s\n", cmp.Diff(tc.exi, exi))
			}
		})
	}
}

func Test_Task_Sync_Get(t *testing.T) {
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
				Sync: &Sync{},
			},
			get: "",
		},
		// Case 003
		{
			tas: &Task{
				Sync: &Sync{"foo": "bar"},
			},
			get: "bar",
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
			get: "",
		},
	}

	for i, tc := range testCases {
		t.Run(fmt.Sprintf("%03d", i), func(t *testing.T) {
			get := tc.tas.Sync.Get("foo")

			if get != tc.get {
				t.Fatalf("\n\n%s\n", cmp.Diff(tc.get, get))
			}
		})
	}
}

func Test_Task_Sync_Set(t *testing.T) {
	testCases := []struct {
		tas *Task
		key string
		val string
		set *Task
	}{
		// Case 000
		{
			tas: &Task{
				Sync: &Sync{"foo": "bar"},
			},
			key: "foo",
			val: "zap",
			set: &Task{
				Sync: &Sync{"foo": "zap"},
			},
		},
		// Case 001
		{
			tas: &Task{
				Sync: &Sync{"foo": "bar", "one": "two"},
			},
			key: "one",
			val: "thr",
			set: &Task{
				Sync: &Sync{"foo": "bar", "one": "thr"},
			},
		},
	}

	for i, tc := range testCases {
		t.Run(fmt.Sprintf("%03d", i), func(t *testing.T) {
			tc.tas.Sync.Set(tc.key, tc.val)

			set := tc.set

			if !reflect.DeepEqual(set, tc.tas) {
				t.Fatalf("\n\n%s\n", cmp.Diff(tc.tas, set))
			}
		})
	}
}
