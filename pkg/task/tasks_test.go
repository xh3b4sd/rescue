package task

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func Test_Tasks_With(t *testing.T) {
	testCases := []struct {
		tas []*Task
		pre []string
		wit []*Task
	}{
		// case 0
		{
			tas: []*Task{},
			pre: []string{},
			wit: nil,
		},
		// case 1
		{
			tas: []*Task{
				{
					Obj: TaskObj{
						Metadata: map[string]string{
							"test.domain.io/id":      "1",
							"test.domain.io/owner":   "a",
							"test.domain.io/version": "5",
						},
					},
				},
			},
			pre: []string{},
			wit: nil,
		},
		// case 2
		{
			tas: []*Task{
				{
					Obj: TaskObj{
						Metadata: map[string]string{
							"test.domain.io/id":      "1",
							"test.domain.io/owner":   "a",
							"test.domain.io/version": "5",
							"some.bucket.io/key":     "foo",
						},
					},
				},
				{
					Obj: TaskObj{
						Metadata: map[string]string{
							"test.domain.io/id":      "2",
							"test.domain.io/owner":   "b",
							"test.domain.io/version": "3",
							"this.bucket.io/key":     "bar",
						},
					},
				},
			},
			pre: []string{
				"test",
			},
			wit: []*Task{
				{
					Obj: TaskObj{
						Metadata: map[string]string{
							"test.domain.io/id":      "1",
							"test.domain.io/owner":   "a",
							"test.domain.io/version": "5",
							"some.bucket.io/key":     "foo",
						},
					},
				},
				{
					Obj: TaskObj{
						Metadata: map[string]string{
							"test.domain.io/id":      "2",
							"test.domain.io/owner":   "b",
							"test.domain.io/version": "3",
							"this.bucket.io/key":     "bar",
						},
					},
				},
			},
		},
		// case 3
		{
			tas: []*Task{
				{
					Obj: TaskObj{
						Metadata: map[string]string{
							"test.domain.io/id":      "1",
							"test.domain.io/owner":   "a",
							"test.domain.io/version": "5",
							"some.bucket.io/key":     "foo",
						},
					},
				},
				{
					Obj: TaskObj{
						Metadata: map[string]string{
							"test.domain.io/id":      "2",
							"test.domain.io/owner":   "b",
							"test.domain.io/version": "3",
							"this.bucket.io/key":     "bar",
						},
					},
				},
			},
			pre: []string{
				"this.buc",
			},
			wit: []*Task{
				{
					Obj: TaskObj{
						Metadata: map[string]string{
							"test.domain.io/id":      "2",
							"test.domain.io/owner":   "b",
							"test.domain.io/version": "3",
							"this.bucket.io/key":     "bar",
						},
					},
				},
			},
		},
		// case 4
		{
			tas: []*Task{
				{
					Obj: TaskObj{
						Metadata: map[string]string{
							"test.domain.io/id":      "1",
							"test.domain.io/owner":   "a",
							"test.domain.io/version": "5",
							"some.bucket.io/key":     "foo",
						},
					},
				},
				{
					Obj: TaskObj{
						Metadata: map[string]string{
							"test.domain.io/id":      "2",
							"test.domain.io/owner":   "b",
							"test.domain.io/version": "3",
							"this.bucket.io/key":     "bar",
						},
					},
				},
			},
			pre: []string{
				"some.bucket.io",
				"this.buck",
			},
			wit: []*Task{
				{
					Obj: TaskObj{
						Metadata: map[string]string{
							"test.domain.io/id":      "1",
							"test.domain.io/owner":   "a",
							"test.domain.io/version": "5",
							"some.bucket.io/key":     "foo",
						},
					},
				},
				{
					Obj: TaskObj{
						Metadata: map[string]string{
							"test.domain.io/id":      "2",
							"test.domain.io/owner":   "b",
							"test.domain.io/version": "3",
							"this.bucket.io/key":     "bar",
						},
					},
				},
			},
		},
		// case 5
		{
			tas: []*Task{
				{
					Obj: TaskObj{
						Metadata: map[string]string{
							"test.domain.io/id":      "1",
							"test.domain.io/owner":   "a",
							"test.domain.io/version": "5",
							"some.bucket.io/key":     "foo",
						},
					},
				},
				{
					Obj: TaskObj{
						Metadata: map[string]string{
							"test.domain.io/id":      "2",
							"test.domain.io/owner":   "b",
							"test.domain.io/version": "3",
							"this.bucket.io/key":     "bar",
						},
					},
				},
			},
			pre: []string{
				"some.others.io",
				"this.buck",
			},
			wit: []*Task{
				{
					Obj: TaskObj{
						Metadata: map[string]string{
							"test.domain.io/id":      "2",
							"test.domain.io/owner":   "b",
							"test.domain.io/version": "3",
							"this.bucket.io/key":     "bar",
						},
					},
				},
			},
		},
	}

	for i, tc := range testCases {
		t.Run(fmt.Sprintf("%03d", i), func(t *testing.T) {
			wit := Tasks(tc.tas).With(tc.pre...)

			if !reflect.DeepEqual(tc.wit, wit) {
				t.Fatalf("\n\n%s\n", cmp.Diff(tc.wit, wit))
			}
		})
	}
}
