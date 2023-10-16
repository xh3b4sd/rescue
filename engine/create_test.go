package engine

import (
	"fmt"
	"testing"

	"github.com/xh3b4sd/redigo"
	"github.com/xh3b4sd/rescue/task"
)

func Test_Engine_Create_Meta_No_Error(t *testing.T) {
	testCases := []struct {
		tas *task.Task
	}{
		// Case 000
		{
			tas: &task.Task{
				Meta: &task.Meta{
					"foo": "bar",
				},
			},
		},
	}

	for i, tc := range testCases {
		t.Run(fmt.Sprintf("%03d", i), func(t *testing.T) {
			var e *Engine
			{
				e = New(Config{
					Redigo: redigo.Fake(),
				})
			}

			err := e.Create(tc.tas)
			if err != nil {
				t.Fatal("expected", nil, "got", err)
			}
		})
	}
}

func Test_Engine_Create_Gate_No_Error(t *testing.T) {
	testCases := []struct {
		tas *task.Task
	}{
		// Case 000 ensures that trigger tasks can define a single trigger label.
		{
			tas: &task.Task{
				Meta: &task.Meta{
					"foo": "bar",
				},
				Gate: &task.Gate{
					"foo": task.Trigger,
				},
			},
		},
		// Case 001 ensures that trigger tasks can define a multiple trigger labels.
		{
			tas: &task.Task{
				Meta: &task.Meta{
					"foo": "bar",
				},
				Gate: &task.Gate{
					"foo": task.Trigger,
					"bar": task.Trigger,
					"baz": task.Trigger,
				},
			},
		},
		// Case 002 ensures that task templates can define a single waiting label.
		{
			tas: &task.Task{
				Meta: &task.Meta{
					"foo": "bar",
				},
				Gate: &task.Gate{
					"baz": task.Waiting,
				},
			},
		},
		// Case 003 ensures that task templates can define a multiple waiting
		// labels.
		{
			tas: &task.Task{
				Meta: &task.Meta{
					"foo": "bar",
				},
				Gate: &task.Gate{
					"foo": task.Waiting,
					"bar": task.Waiting,
					"baz": task.Waiting,
				},
			},
		},
	}

	for i, tc := range testCases {
		t.Run(fmt.Sprintf("%03d", i), func(t *testing.T) {
			var e *Engine
			{
				e = New(Config{
					Redigo: redigo.Fake(),
				})
			}

			err := e.Create(tc.tas)
			if err != nil {
				t.Fatal("expected", nil, "got", err)
			}
		})
	}
}

func Test_Engine_Create_Gate_Error(t *testing.T) {
	testCases := []struct {
		tas *task.Task
	}{
		// Case 000
		{
			tas: &task.Task{
				Meta: &task.Meta{
					"foo": "bar",
				},
				Gate: &task.Gate{
					"foo": "bar",
				},
			},
		},
		// Case 001
		{
			tas: &task.Task{
				Meta: &task.Meta{
					"foo": "bar",
				},
				Gate: &task.Gate{
					"foo": "bar",
					"baz": "zap",
				},
			},
		},
		// Case 002
		{
			tas: &task.Task{
				Meta: &task.Meta{
					"foo": "bar",
				},
				Gate: &task.Gate{
					"foo": task.Deleted,
				},
			},
		},
		// Case 003
		{
			tas: &task.Task{
				Meta: &task.Meta{
					"foo": "bar",
				},
				Gate: &task.Gate{
					"foo": task.Deleted,
					"bar": task.Waiting,
				},
			},
		},
		// Case 004
		{
			tas: &task.Task{
				Meta: &task.Meta{
					"foo": "bar",
				},
				Gate: &task.Gate{
					"foo": task.Deleted,
					"baz": "zap",
					"bar": task.Waiting,
				},
			},
		},
		// Case 005
		{
			tas: &task.Task{
				Meta: &task.Meta{
					"foo": "bar",
				},
				Gate: &task.Gate{
					"foo": task.Deleted,
					"bar": task.Trigger,
				},
			},
		},
		// Case 006
		{
			tas: &task.Task{
				Meta: &task.Meta{
					"foo": "bar",
				},
				Gate: &task.Gate{
					"foo": task.Deleted,
					"baz": "zap",
					"bar": task.Trigger,
				},
			},
		},
		// Case 007
		{
			tas: &task.Task{
				Meta: &task.Meta{
					"foo": "bar",
				},
				Gate: &task.Gate{
					"foo": task.Trigger,
					"baz": task.Waiting,
				},
			},
		},
		// Case 008
		{
			tas: &task.Task{
				Meta: &task.Meta{
					"foo": "bar",
				},
				Gate: &task.Gate{
					"foo": task.Trigger,
					"bar": "zap",
					"baz": task.Waiting,
				},
			},
		},
		// Case 009
		{
			tas: &task.Task{
				Meta: &task.Meta{
					"foo": "bar",
				},
				Gate: &task.Gate{
					"foo": "bar",
					"baz": task.Waiting,
				},
			},
		},
		// Case 010
		{
			tas: &task.Task{
				Meta: &task.Meta{
					"foo": "bar",
				},
				Gate: &task.Gate{
					"foo": "bar",
					"baz": task.Deleted,
				},
			},
		},
		// Case 011
		{
			tas: &task.Task{
				Meta: &task.Meta{
					"foo": "bar",
				},
				Gate: &task.Gate{
					"foo": task.Trigger,
					"baz": "zap",
				},
			},
		},
	}

	for i, tc := range testCases {
		t.Run(fmt.Sprintf("%03d", i), func(t *testing.T) {
			var e *Engine
			{
				e = New(Config{
					Redigo: redigo.Fake(),
				})
			}

			err := e.Create(tc.tas)
			if err == nil {
				t.Fatal("expected", "error", "got", nil)
			}
		})
	}
}
