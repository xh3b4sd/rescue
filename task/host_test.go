package task

import (
	"fmt"
	"reflect"
	"slices"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func Test_Task_Host_Emp(t *testing.T) {
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
				Host: &Host{},
			},
			emp: true,
		},
		// Case 003
		{
			tas: &Task{
				Host: &Host{"foo": "bar"},
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
			emp := tc.tas.Host.Emp()

			if emp != tc.emp {
				t.Fatalf("\n\n%s\n", cmp.Diff(tc.emp, emp))
			}
		})
	}
}

func Test_Task_Host_Eql(t *testing.T) {
	testCases := []struct {
		tas *Task
		hos *Host
		eql bool
	}{
		// Case 000
		{
			tas: &Task{},
			hos: nil,
			eql: false,
		},
		// Case 001
		{
			tas: &Task{},
			hos: &Host{},
			eql: false,
		},
		// Case 002
		{
			tas: &Task{
				Host: &Host{"foo": "bar"},
			},
			hos: &Host{},
			eql: false,
		},
		// Case 003
		{
			tas: &Task{
				Core: &Core{"foo": "bar"},
			},
			hos: &Host{"foo": "bar"},
			eql: false,
		},
		// Case 004
		{
			tas: &Task{
				Host: &Host{"foo": "bar"},
			},
			hos: &Host{"foo": "bar"},
			eql: true,
		},
		// Case 005
		{
			tas: &Task{
				Host: &Host{"foo": "bar", "baz": "zap"},
			},
			hos: &Host{"foo": "bar"},
			eql: false,
		},
		// Case 006
		{
			tas: &Task{
				Host: &Host{"foo": "bar", "baz": "zap"},
			},
			hos: &Host{"baz": "zap"},
			eql: false,
		},
		// Case 007
		{
			tas: &Task{
				Host: &Host{"foo": "bar", "baz": "zap"},
			},
			hos: &Host{"foo": "bar", "baz": "zap"},
			eql: true,
		},
		// Case 008
		{
			tas: &Task{
				Host: &Host{"foo": "", "baz": "zap"},
			},
			hos: &Host{"foo": "", "baz": "zap"},
			eql: true,
		},
		// Case 009
		{
			tas: &Task{
				Host: &Host{"foo": "", "baz": "zap"},
			},
			hos: &Host{"foo": "bar", "baz": ""},
			eql: false,
		},
	}

	for i, tc := range testCases {
		t.Run(fmt.Sprintf("%03d", i), func(t *testing.T) {
			eql := tc.tas.Host.Eql(tc.hos)

			if eql != tc.eql {
				t.Fatalf("\n\n%s\n", cmp.Diff(tc.eql, eql))
			}
		})
	}
}

func Test_Task_Host_Exi(t *testing.T) {
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
				Host: &Host{},
			},
			exi: false,
		},
		// Case 003
		{
			tas: &Task{
				Host: &Host{"foo": "bar"},
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
			exi := tc.tas.Host.Exi("foo")

			if exi != tc.exi {
				t.Fatalf("\n\n%s\n", cmp.Diff(tc.exi, exi))
			}
		})
	}
}

func Test_Task_Host_Get(t *testing.T) {
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
				Host: &Host{},
			},
			get: "",
		},
		// Case 003
		{
			tas: &Task{
				Host: &Host{"foo": "bar"},
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
			get := tc.tas.Host.Get("foo")

			if get != tc.get {
				t.Fatalf("\n\n%s\n", cmp.Diff(tc.get, get))
			}
		})
	}
}

func Test_Task_Host_Key(t *testing.T) {
	testCases := []struct {
		tas *Task
		key []string
	}{
		// Case 000
		{
			tas: &Task{},
			key: nil,
		},
		// Case 001
		{
			tas: &Task{
				Meta: &Meta{
					"foo": "bar",
				},
			},
			key: nil,
		},
		// Case 002
		{
			tas: &Task{
				Host: &Host{
					"foo": "bar",
				},
			},
			key: []string{
				"foo",
			},
		},
		// Case 003
		{
			tas: &Task{
				Host: &Host{
					"foo": "bar",
					"baz": "foo",
					"key": "baz",
				},
			},
			key: []string{
				"foo",
				"baz",
				"key",
			},
		},
	}

	for i, tc := range testCases {
		t.Run(fmt.Sprintf("%03d", i), func(t *testing.T) {
			key := tc.tas.Host.Key()

			slices.Sort(key)
			slices.Sort(tc.key)

			if !reflect.DeepEqual(key, tc.key) {
				t.Fatalf("\n\n%s\n", cmp.Diff(tc.key, key))
			}
		})
	}
}

func Test_Task_Host_Set(t *testing.T) {
	testCases := []struct {
		tas *Task
		key string
		val string
		set *Task
	}{
		// Case 000
		{
			tas: &Task{
				Host: &Host{"foo": "bar"},
			},
			key: "foo",
			val: "zap",
			set: &Task{
				Host: &Host{"foo": "zap"},
			},
		},
		// Case 001
		{
			tas: &Task{
				Host: &Host{"foo": "bar", "one": "two"},
			},
			key: "one",
			val: "thr",
			set: &Task{
				Host: &Host{"foo": "bar", "one": "thr"},
			},
		},
	}

	for i, tc := range testCases {
		t.Run(fmt.Sprintf("%03d", i), func(t *testing.T) {
			tc.tas.Host.Set(tc.key, tc.val)

			set := tc.set

			if !reflect.DeepEqual(set, tc.tas) {
				t.Fatalf("\n\n%s\n", cmp.Diff(tc.tas, set))
			}
		})
	}
}
