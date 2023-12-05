package task

import (
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func Test_Task_Pag(t *testing.T) {
	testCases := []struct {
		tas *Task
		pag bool
	}{
		// Case 000
		{
			tas: &Task{},
			pag: false,
		},
		// Case 001
		{
			tas: &Task{
				Core: &Core{"foo": "bar"},
			},
			pag: false,
		},
		// Case 002
		{
			tas: &Task{
				Sync: &Sync{},
			},
			pag: false,
		},
		// Case 003
		{
			tas: &Task{
				Sync: &Sync{"foo": "bar"},
			},
			pag: false,
		},
		// Case 004
		{
			tas: &Task{
				Root: &Root{},
			},
			pag: false,
		},
		// Case 005
		{
			tas: &Task{
				Root: &Root{"foo": "bar"},
			},
			pag: false,
		},
		// Case 006
		{
			tas: &Task{
				Sync: &Sync{},
			},
			pag: false,
		},
		// Case 007
		{
			tas: &Task{
				Sync: &Sync{Paging: ""},
			},
			pag: false,
		},
		// Case 008
		{
			tas: &Task{
				Sync: &Sync{Paging: " "},
			},
			pag: false,
		},
		// Case 009
		{
			tas: &Task{
				Sync: &Sync{Paging: "0"},
			},
			pag: false,
		},
		// Case 010
		{
			tas: &Task{
				Sync: &Sync{Paging: "1"},
			},
			pag: true,
		},
		// Case 011
		{
			tas: &Task{
				Sync: &Sync{Paging: "2376452"},
			},
			pag: true,
		},
		// Case 012
		{
			tas: &Task{
				Sync: &Sync{Paging: "foo"},
			},
			pag: true,
		},
	}

	for i, tc := range testCases {
		t.Run(fmt.Sprintf("%03d", i), func(t *testing.T) {
			pag := tc.tas.Pag()

			if pag != tc.pag {
				t.Fatalf("\n\n%s\n", cmp.Diff(tc.pag, pag))
			}
		})
	}
}
