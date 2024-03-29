package task

import (
	"fmt"
	"reflect"
	"slices"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func Test_Task_Cron_Emp(t *testing.T) {
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
				Meta: &Meta{"foo": "bar"},
			},
			emp: true,
		},
		// Case 002
		{
			tas: &Task{
				Cron: &Cron{},
			},
			emp: true,
		},
		// Case 003
		{
			tas: &Task{
				Cron: &Cron{"foo": "bar"},
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
			emp := tc.tas.Cron.Emp()

			if emp != tc.emp {
				t.Fatalf("\n\n%s\n", cmp.Diff(tc.emp, emp))
			}
		})
	}
}

func Test_Task_Cron_Eql(t *testing.T) {
	testCases := []struct {
		tas *Task
		crn *Cron
		eql bool
	}{
		// Case 000
		{
			tas: &Task{},
			crn: nil,
			eql: false,
		},
		// Case 001
		{
			tas: &Task{},
			crn: &Cron{},
			eql: false,
		},
		// Case 002
		{
			tas: &Task{
				Cron: &Cron{"foo": "bar"},
			},
			crn: &Cron{},
			eql: false,
		},
		// Case 003
		{
			tas: &Task{
				Meta: &Meta{"foo": "bar"},
			},
			crn: &Cron{"foo": "bar"},
			eql: false,
		},
		// Case 004
		{
			tas: &Task{
				Cron: &Cron{"foo": "bar"},
			},
			crn: &Cron{"foo": "bar"},
			eql: true,
		},
		// Case 005
		{
			tas: &Task{
				Cron: &Cron{"foo": "bar", "baz": "zap"},
			},
			crn: &Cron{"foo": "bar"},
			eql: false,
		},
		// Case 006
		{
			tas: &Task{
				Cron: &Cron{"foo": "bar", "baz": "zap"},
			},
			crn: &Cron{"baz": "zap"},
			eql: false,
		},
		// Case 007
		{
			tas: &Task{
				Cron: &Cron{"foo": "bar", "baz": "zap"},
			},
			crn: &Cron{"foo": "bar", "baz": "zap"},
			eql: true,
		},
		// Case 008
		{
			tas: &Task{
				Cron: &Cron{"foo": "", "baz": "zap"},
			},
			crn: &Cron{"foo": "", "baz": "zap"},
			eql: true,
		},
		// Case 009
		{
			tas: &Task{
				Cron: &Cron{"foo": "", "baz": "zap"},
			},
			crn: &Cron{"foo": "bar", "baz": ""},
			eql: false,
		},
	}

	for i, tc := range testCases {
		t.Run(fmt.Sprintf("%03d", i), func(t *testing.T) {
			eql := tc.tas.Cron.Eql(tc.crn)

			if eql != tc.eql {
				t.Fatalf("\n\n%s\n", cmp.Diff(tc.eql, eql))
			}
		})
	}
}

func Test_Task_Cron_Key(t *testing.T) {
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
				Cron: &Cron{
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
				Cron: &Cron{
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
			key := tc.tas.Cron.Key()

			slices.Sort(key)
			slices.Sort(tc.key)

			if !reflect.DeepEqual(key, tc.key) {
				t.Fatalf("\n\n%s\n", cmp.Diff(tc.key, key))
			}
		})
	}
}
