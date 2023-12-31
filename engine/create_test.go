package engine

import (
	"fmt"
	"testing"
	"time"

	"github.com/xh3b4sd/redigo"
	"github.com/xh3b4sd/redigo/locker"
	"github.com/xh3b4sd/rescue/task"
	"github.com/xh3b4sd/rescue/ticker"
)

func Test_Engine_Create_Core_Error(t *testing.T) {
	testCases := []struct {
		tas *task.Task
	}{
		// Case 000
		{
			tas: &task.Task{
				Core: &task.Core{
					task.Method: "foo",
				},
				Meta: &task.Meta{
					"foo": "bar",
				},
			},
		},
		// Case 001
		{
			tas: &task.Task{
				Core: &task.Core{
					task.Object: "bar",
				},
				Meta: &task.Meta{
					"foo": "bar",
				},
			},
		},
		// Case 002
		{
			tas: &task.Task{
				Core: &task.Core{
					task.Worker: "baz",
				},
				Meta: &task.Meta{
					"foo": "bar",
				},
			},
		},
		// Case 003
		{
			tas: &task.Task{
				Core: &task.Core{
					task.Method: task.MthdAll,
					task.Worker: "baz",
				},
				Meta: &task.Meta{
					"foo": "bar",
				},
			},
		},
		// Case 004
		{
			tas: &task.Task{
				Core: &task.Core{
					task.Method: task.MthdAny,
					task.Object: "bar",
				},
				Meta: &task.Meta{
					"foo": "bar",
				},
			},
		},
		// Case 005
		{
			tas: &task.Task{
				Core: &task.Core{
					task.Method: task.MthdAll,
				},
				Meta: &task.Meta{
					"foo": "bar",
				},
			},
		},
		// Case 006
		{
			tas: &task.Task{
				Core: &task.Core{
					task.Method: task.MthdAny,
				},
				Meta: &task.Meta{
					"foo": "bar",
				},
			},
		},
		// Case 007
		{
			tas: &task.Task{
				Core: &task.Core{
					task.Method: task.MthdUni,
					task.Worker: "bar",
				},
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
					Locker: &locker.Fake{},
				})
			}

			err := e.Create(tc.tas)
			if err == nil {
				t.Fatal("expected", "error", "got", nil)
			}
		})
	}
}

func Test_Engine_Create_Cron_Error(t *testing.T) {
	testCases := []struct {
		tas *task.Task
	}{
		// Case 000 ensures that @every and @exact cannot be defined together.
		{
			tas: &task.Task{
				Cron: &task.Cron{
					task.Aevery: "hour",
					task.Aexact: time.Now().UTC().Add(5 * time.Minute).Format(ticker.Layout),
				},
				Meta: &task.Meta{
					"foo": "bar",
				},
			},
		},
		// Case 001 ensures that @exact cannot be defined in the past.
		{
			tas: &task.Task{
				Cron: &task.Cron{
					task.Aexact: "2023-09-28T12:00:00",
				},
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
					Locker: &locker.Fake{},
				})
			}

			err := e.Create(tc.tas)
			if err == nil {
				t.Fatal("expected", "error", "got", nil)
			}
		})
	}
}

func Test_Engine_Create_Cron_No_Error(t *testing.T) {
	testCases := []struct {
		tas *task.Task
	}{
		// Case 000 ensures that @every can be defined.
		{
			tas: &task.Task{
				Cron: &task.Cron{
					task.Aevery: "hour",
				},
				Meta: &task.Meta{
					"foo": "bar",
				},
			},
		},
		// Case 001 ensures that @exact can be defined.
		{
			tas: &task.Task{
				Cron: &task.Cron{
					task.Aexact: time.Now().UTC().Add(5 * time.Minute).Format(ticker.Layout),
				},
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
					Locker: &locker.Fake{},
				})
			}

			err := e.Create(tc.tas)
			if err != nil {
				t.Fatal("expected", nil, "got", err)
			}
		})
	}
}

func Test_Engine_Create_Host_Error(t *testing.T) {
	testCases := []struct {
		tas *task.Task
	}{
		// Case 000
		{
			tas: &task.Task{
				Meta: &task.Meta{
					"foo": "bar",
				},
				Node: &task.Node{
					task.Method: "foo",
				},
			},
		},
		// Case 001
		{
			tas: &task.Task{
				Meta: &task.Meta{
					"foo": "bar",
				},
				Node: &task.Node{
					task.Object: "bar",
				},
			},
		},
		// Case 002
		{
			tas: &task.Task{
				Meta: &task.Meta{
					"foo": "bar",
				},
				Node: &task.Node{
					task.Worker: "baz",
				},
			},
		},
		// Case 003
		{
			tas: &task.Task{
				Meta: &task.Meta{
					"foo": "bar",
				},
				Node: &task.Node{
					task.Method: task.MthdAll,
					task.Worker: "baz",
				},
			},
		},
		// Case 004
		{
			tas: &task.Task{
				Meta: &task.Meta{
					"foo": "bar",
				},
				Node: &task.Node{
					task.Method: task.MthdAny,
					task.Object: "bar",
				},
			},
		},
		// Case 005
		{
			tas: &task.Task{
				Core: &task.Core{
					task.Method: task.MthdAny,
				},
				Meta: &task.Meta{
					"foo": "bar",
				},
			},
		},
		// Case 006
		{
			tas: &task.Task{
				Meta: &task.Meta{
					"foo": "bar",
				},
				Node: &task.Node{
					task.Method: task.MthdUni,
				},
			},
		},
		// Case 007
		{
			tas: &task.Task{
				Meta: &task.Meta{
					"foo": "bar",
				},
				Node: &task.Node{
					task.Method: task.MthdUni,
					task.Object: "foo",
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
					Locker: &locker.Fake{},
				})
			}

			err := e.Create(tc.tas)
			if err == nil {
				t.Fatal("expected", "error", "got", nil)
			}
		})
	}
}

func Test_Engine_Create_Host_No_Error(t *testing.T) {
	testCases := []struct {
		tas *task.Task
	}{
		// Case 000
		{
			tas: &task.Task{
				Meta: &task.Meta{
					"foo": "bar",
				},
				Node: &task.Node{
					task.Method: task.MthdAll,
				},
			},
		},
		// Case 001
		{
			tas: &task.Task{
				Meta: &task.Meta{
					"foo": "bar",
				},
				Node: &task.Node{
					task.Method: task.MthdAny,
				},
			},
		},
		// Case 002
		{
			tas: &task.Task{
				Meta: &task.Meta{
					"foo": "bar",
				},
				Node: &task.Node{
					task.Method: task.MthdUni,
					task.Worker: "foo",
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
					Locker: &locker.Fake{},
				})
			}

			err := e.Create(tc.tas)
			if err != nil {
				t.Fatal("expected", nil, "got", err)
			}
		})
	}
}

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
					Locker: &locker.Fake{},
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
		// Case 012 ensures that the reserved value "waiting" in Task.Gate cannot be
		// used together with Task.Cron.
		{
			tas: &task.Task{
				Cron: &task.Cron{
					task.Aevery: "hour",
				},
				Meta: &task.Meta{
					"foo": "bar",
				},
				Gate: &task.Gate{
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
					Locker: &locker.Fake{},
				})
			}

			err := e.Create(tc.tas)
			if err == nil {
				t.Fatal("expected", "error", "got", nil)
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
		// Case 001 ensures that task templates defining Task.Cron can schedule
		// trigger tasks with Task.Gate definitions containing the reserved value
		// "trigger".
		{
			tas: &task.Task{
				Cron: &task.Cron{
					task.Aevery: "hour",
				},
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
		// Case 002 ensures that trigger tasks can define a multiple trigger labels.
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
		// Case 003 ensures that task templates can define a single waiting label.
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
		// Case 004 ensures that task templates can define a multiple waiting
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
					Locker: &locker.Fake{},
				})
			}

			err := e.Create(tc.tas)
			if err != nil {
				t.Fatal("expected", nil, "got", err)
			}
		})
	}
}

func Test_Engine_Create_Sync_Error(t *testing.T) {
	testCases := []struct {
		tas *task.Task
	}{
		// Case 000
		{
			tas: &task.Task{
				Meta: &task.Meta{
					"foo": "bar",
				},
				Sync: &task.Sync{
					task.Paging: "0",
					task.Worker: "1234",
				},
			},
		},
		// Case 001
		{
			tas: &task.Task{
				Meta: &task.Meta{
					"foo": "bar",
				},
				Sync: &task.Sync{
					task.Object: "1234",
					task.Paging: "1",
				},
			},
		},
		// Case 002
		{
			tas: &task.Task{
				Meta: &task.Meta{
					"foo": "bar",
				},
				Sync: &task.Sync{
					task.Expiry: "1234",
					task.Object: "1234",
					task.Paging: "2837652",
				},
			},
		},
		// Case 003
		{
			tas: &task.Task{
				Meta: &task.Meta{
					"foo": "bar",
				},
				Sync: &task.Sync{
					task.Paging: "foo",
					task.Aevery: "hour",
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
					Locker: &locker.Fake{},
				})
			}

			err := e.Create(tc.tas)
			if err == nil {
				t.Fatal("expected", "error", "got", nil)
			}
		})
	}
}

func Test_Engine_Create_Sync_No_Error(t *testing.T) {
	testCases := []struct {
		tas *task.Task
	}{
		// Case 000
		{
			tas: &task.Task{
				Meta: &task.Meta{
					"foo": "bar",
				},
				Sync: &task.Sync{
					task.Paging: "0",
				},
			},
		},
		// Case 001
		{
			tas: &task.Task{
				Meta: &task.Meta{
					"foo": "bar",
				},
				Sync: &task.Sync{
					task.Paging: "1",
				},
			},
		},
		// Case 002
		{
			tas: &task.Task{
				Meta: &task.Meta{
					"foo": "bar",
				},
				Sync: &task.Sync{
					task.Paging: "2837652",
				},
			},
		},
		// Case 003
		{
			tas: &task.Task{
				Meta: &task.Meta{
					"foo": "bar",
				},
				Sync: &task.Sync{
					task.Paging: "foo",
				},
			},
		},
		// Case 004
		{
			tas: &task.Task{
				Meta: &task.Meta{
					"foo": "bar",
				},
				Sync: &task.Sync{
					task.Paging: "foo",
					"whatever":  "this",
				},
			},
		},
		// Case 005
		{
			tas: &task.Task{
				Meta: &task.Meta{
					"foo": "bar",
				},
				Sync: &task.Sync{
					"whatever": "this",
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
					Locker: &locker.Fake{},
				})
			}

			err := e.Create(tc.tas)
			if err != nil {
				t.Fatal("expected", nil, "got", err)
			}
		})
	}
}
