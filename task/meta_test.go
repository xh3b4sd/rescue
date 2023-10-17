package task

import (
	"fmt"
	"reflect"
	"slices"
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
			emp := tc.tas.Meta.Emp()

			if emp != tc.emp {
				t.Fatalf("\n\n%s\n", cmp.Diff(tc.emp, emp))
			}
		})
	}
}

func Test_Task_Meta_Eql(t *testing.T) {
	testCases := []struct {
		tas *Task
		met *Meta
		eql bool
	}{
		// Case 000
		{
			tas: &Task{},
			met: nil,
			eql: false,
		},
		// Case 001
		{
			tas: &Task{},
			met: &Meta{},
			eql: false,
		},
		// Case 002
		{
			tas: &Task{
				Meta: &Meta{"foo": "bar"},
			},
			met: &Meta{},
			eql: false,
		},
		// Case 003
		{
			tas: &Task{
				Core: &Core{"foo": "bar"},
			},
			met: &Meta{"foo": "bar"},
			eql: false,
		},
		// Case 004
		{
			tas: &Task{
				Meta: &Meta{"foo": "bar"},
			},
			met: &Meta{"foo": "bar"},
			eql: true,
		},
		// Case 005
		{
			tas: &Task{
				Meta: &Meta{"foo": "bar", "baz": "zap"},
			},
			met: &Meta{"foo": "bar"},
			eql: false,
		},
		// Case 006
		{
			tas: &Task{
				Meta: &Meta{"foo": "bar", "baz": "zap"},
			},
			met: &Meta{"baz": "zap"},
			eql: false,
		},
		// Case 007
		{
			tas: &Task{
				Meta: &Meta{"foo": "bar", "baz": "zap"},
			},
			met: &Meta{"foo": "bar", "baz": "zap"},
			eql: true,
		},
		// Case 008
		{
			tas: &Task{
				Meta: &Meta{"foo": "", "baz": "zap"},
			},
			met: &Meta{"foo": "", "baz": "zap"},
			eql: true,
		},
		// Case 009
		{
			tas: &Task{
				Meta: &Meta{"foo": "", "baz": "zap"},
			},
			met: &Meta{"foo": "bar", "baz": ""},
			eql: false,
		},
	}

	for i, tc := range testCases {
		t.Run(fmt.Sprintf("%03d", i), func(t *testing.T) {
			eql := tc.tas.Meta.Eql(tc.met)

			if eql != tc.eql {
				t.Fatalf("\n\n%s\n", cmp.Diff(tc.eql, eql))
			}
		})
	}
}

func Test_Task_Meta_Exi(t *testing.T) {
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
			exi := tc.tas.Meta.Exi("foo")

			if exi != tc.exi {
				t.Fatalf("\n\n%s\n", cmp.Diff(tc.exi, exi))
			}
		})
	}
}

func Test_Task_Meta_Get(t *testing.T) {
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
			get := tc.tas.Meta.Get("foo")

			if get != tc.get {
				t.Fatalf("\n\n%s\n", cmp.Diff(tc.get, get))
			}
		})
	}
}

func Test_Task_Meta_Key(t *testing.T) {
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
				Root: &Root{
					"foo": "bar",
				},
			},
			key: nil,
		},
		// Case 002
		{
			tas: &Task{
				Meta: &Meta{
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
				Meta: &Meta{
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
			key := tc.tas.Meta.Key()

			slices.Sort(key)
			slices.Sort(tc.key)

			if !reflect.DeepEqual(key, tc.key) {
				t.Fatalf("\n\n%s\n", cmp.Diff(tc.key, key))
			}
		})
	}
}

func Test_Task_Meta_Set(t *testing.T) {
	testCases := []struct {
		tas *Task
		key string
		val string
		set *Task
	}{
		// Case 000
		{
			tas: &Task{
				Meta: &Meta{"foo": "bar"},
			},
			key: "foo",
			val: "zap",
			set: &Task{
				Meta: &Meta{"foo": "zap"},
			},
		},
		// Case 001
		{
			tas: &Task{
				Meta: &Meta{"foo": "bar", "one": "two"},
			},
			key: "one",
			val: "thr",
			set: &Task{
				Meta: &Meta{"foo": "bar", "one": "thr"},
			},
		},
	}

	for i, tc := range testCases {
		t.Run(fmt.Sprintf("%03d", i), func(t *testing.T) {
			tc.tas.Meta.Set(tc.key, tc.val)

			set := tc.set

			if !reflect.DeepEqual(set, tc.tas) {
				t.Fatalf("\n\n%s\n", cmp.Diff(tc.tas, set))
			}
		})
	}
}
