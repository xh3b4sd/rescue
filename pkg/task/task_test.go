package task

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/xh3b4sd/rescue/pkg/metadata"
)

func Test_Task_Pref(t *testing.T) {
	testCases := []struct {
		tas *Task
		pre []string
		wit *Task
	}{
		// case 0
		{
			tas: &Task{},
			pre: []string{},
			wit: nil,
		},
		// case 1
		{
			tas: &Task{
				Obj: TaskObj{
					Metadata: map[string]string{
						"test.domain.io/id":      "1",
						"test.domain.io/owner":   "a",
						"test.domain.io/version": "5",
						"some.bucket.io/key":     "foo",
					},
				},
			},
			pre: []string{},
			wit: nil,
		},
		// case 2
		{
			tas: &Task{
				Obj: TaskObj{
					Metadata: map[string]string{
						"test.domain.io/id":      "1",
						"test.domain.io/owner":   "a",
						"test.domain.io/version": "5",
						"some.bucket.io/key":     "foo",
					},
				},
			},
			pre: []string{
				"test",
			},
			wit: &Task{
				Obj: TaskObj{
					Metadata: map[string]string{
						"test.domain.io/id":      "1",
						"test.domain.io/owner":   "a",
						"test.domain.io/version": "5",
					},
				},
			},
		},
		// case 3
		{
			tas: &Task{
				Obj: TaskObj{
					Metadata: map[string]string{
						"test.domain.io/id":      "1",
						"test.domain.io/owner":   "a",
						"test.domain.io/version": "5",
						"some.bucket.io/key":     "foo",
					},
				},
			},
			pre: []string{
				"this.buc",
			},
			wit: nil,
		},
		// case 4
		{
			tas: &Task{
				Obj: TaskObj{
					Metadata: map[string]string{
						"test.domain.io/id":      "1",
						"test.domain.io/owner":   "a",
						"test.domain.io/version": "5",
						"some.bucket.io/key":     "foo",
					},
				},
			},
			pre: []string{
				"some.bucket.io",
				"this.buck",
			},
			wit: &Task{
				Obj: TaskObj{
					Metadata: map[string]string{
						"some.bucket.io/key": "foo",
					},
				},
			},
		},
		// case 5
		{
			tas: &Task{
				Obj: TaskObj{
					Metadata: map[string]string{
						"test.domain.io/id":      "1",
						"test.domain.io/owner":   "a",
						"test.domain.io/version": "5",
						"some.bucket.io/key":     "foo",
					},
				},
			},
			pre: []string{
				"some.others.io",
				"this.buck",
			},
			wit: nil,
		},
	}

	for i, tc := range testCases {
		t.Run(fmt.Sprintf("%03d", i), func(t *testing.T) {
			wit := tc.tas.Pref(tc.pre...)

			if !reflect.DeepEqual(tc.wit, wit) {
				t.Fatalf("\n\n%s\n", cmp.Diff(tc.wit, wit))
			}
		})
	}
}

func Test_Task_With(t *testing.T) {
	testCases := []struct {
		tas *Task
		key []string
		wit *Task
	}{
		// case 0
		{
			tas: &Task{},
			key: []string{},
			wit: &Task{},
		},
		// case 1
		{
			tas: &Task{},
			key: []string{
				"foo",
				metadata.ID,
			},
			wit: nil,
		},
		// case 2
		{
			tas: &Task{
				Obj: TaskObj{
					Metadata: map[string]string{
						metadata.ID: "1",
					},
				},
			},
			key: []string{
				"foo",
			},
			wit: nil,
		},
		// case 3
		{
			tas: &Task{
				Obj: TaskObj{
					Metadata: map[string]string{
						metadata.ID: "1",
					},
				},
			},
			key: []string{
				"foo", // foo not in tas
				metadata.ID,
			},
			wit: nil,
		},
		// case 4
		{
			tas: &Task{
				Obj: TaskObj{
					Metadata: map[string]string{
						metadata.ID: "1",
					},
				},
			},
			key: []string{
				metadata.ID,
			},
			wit: &Task{
				Obj: TaskObj{
					Metadata: map[string]string{
						metadata.ID: "1",
					},
				},
			},
		},
		// case 5
		{
			tas: &Task{
				Obj: TaskObj{
					Metadata: map[string]string{
						metadata.ID:      "1",
						metadata.Owner:   "a",
						metadata.Version: "5",
					},
				},
			},
			key: []string{
				metadata.ID,
			},
			wit: &Task{
				Obj: TaskObj{
					Metadata: map[string]string{
						metadata.ID: "1",
					},
				},
			},
		},
		// case 6
		{
			tas: &Task{
				Obj: TaskObj{
					Metadata: map[string]string{
						metadata.ID:      "1",
						metadata.Owner:   "a",
						metadata.Version: "5",
					},
				},
			},
			key: []string{
				metadata.ID,
				metadata.Owner,
			},
			wit: &Task{
				Obj: TaskObj{
					Metadata: map[string]string{
						metadata.ID:    "1",
						metadata.Owner: "a",
					},
				},
			},
		},
	}

	for i, tc := range testCases {
		t.Run(fmt.Sprintf("%03d", i), func(t *testing.T) {
			wit := tc.tas.With(tc.key...)

			if !reflect.DeepEqual(tc.wit, wit) {
				t.Fatalf("\n\n%s\n", cmp.Diff(tc.wit, wit))
			}
		})
	}
}
