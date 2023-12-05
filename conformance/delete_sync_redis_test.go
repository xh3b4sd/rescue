//go:build redis

package conformance

import (
	"reflect"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/xh3b4sd/logger"
	"github.com/xh3b4sd/redigo"
	"github.com/xh3b4sd/rescue"
	"github.com/xh3b4sd/rescue/engine"
	"github.com/xh3b4sd/rescue/task"
)

func Test_Engine_Delete_Sync(t *testing.T) {
	var err error

	var eon rescue.Interface
	{
		eon = engine.New(engine.Config{
			Logger: logger.Fake(),
			Redigo: prgAll(redigo.Default()),
		})
	}

	var lis []*task.Task
	{
		lis, err = eon.Lister(engine.All())
		if err != nil {
			t.Fatal(err)
		}
	}

	{
		if len(lis) != 0 {
			t.Fatal("expected", 0, "got", len(lis))
		}
	}

	// Create one task with custom sync state.
	var tas *task.Task
	{
		tas = &task.Task{
			Meta: &task.Meta{
				"test.api.io/key": "foo",
			},
			Sync: &task.Sync{
				"test.api.io/one": "n/a",
			},
		}
	}

	{
		err = eon.Create(tas)
		if err != nil {
			t.Fatal(err)
		}
	}

	// Create another task with paging sync state.
	{
		tas = &task.Task{
			Meta: &task.Meta{
				"test.api.io/key": "bar",
			},
			Sync: &task.Sync{
				task.Paging: "0",
			},
		}
	}

	{
		err = eon.Create(tas)
		if err != nil {
			t.Fatal(err)
		}
	}

	{
		lis, err = eon.Lister(engine.All())
		if err != nil {
			t.Fatal(err)
		}
	}

	// We should have our two tasks stored now.
	{
		if len(lis) != 2 {
			t.Fatal("expected", 2, "got", len(lis))
		}
	}

	{
		tas, err = eon.Search()
		if err != nil {
			t.Fatal(err)
		}
	}

	{
		var exp *task.Task
		{
			exp = &task.Task{
				Core: tas.Core,
				Meta: &task.Meta{
					"test.api.io/key": "foo",
				},
				Node: &task.Node{
					task.Method: task.MthdAny,
				},
				Sync: &task.Sync{
					"test.api.io/one": "n/a",
				},
			}
		}

		{
			if !reflect.DeepEqual(tas, exp) {
				t.Fatalf("\n\n%s\n", cmp.Diff(exp, tas))
			}
		}
	}

	{
		err = eon.Delete(tas)
		if err != nil {
			t.Fatal(err)
		}
	}

	{
		tas, err = eon.Search()
		if err != nil {
			t.Fatal(err)
		}
	}

	{
		var exp *task.Task
		{
			exp = &task.Task{
				Core: tas.Core,
				Meta: &task.Meta{
					"test.api.io/key": "bar",
				},
				Node: &task.Node{
					task.Method: task.MthdAny,
				},
				Sync: &task.Sync{
					task.Paging: "0",
				},
			}
		}

		{
			if !reflect.DeepEqual(tas, exp) {
				t.Fatalf("\n\n%s\n", cmp.Diff(exp, tas))
			}
		}
	}

	// Set paging state to our paging task.
	{
		tas.Sync.Set(task.Paging, "233")
	}

	// Calling Delete on our paging task should expire the task together with the
	// synced paging state.
	{
		err = eon.Delete(tas)
		if err != nil {
			t.Fatal(err)
		}
	}

	{
		lis, err = eon.Lister(engine.All())
		if err != nil {
			t.Fatal(err)
		}
	}

	// Now, even after two calls of Engine.Delete we have one task remaining,
	// which is our paging task.
	{
		if len(lis) != 1 {
			t.Fatal("expected", 1, "got", len(lis))
		}
	}

	{
		tas, err = eon.Search()
		if err != nil {
			t.Fatal(err)
		}
	}

	{
		var exp *task.Task
		{
			exp = &task.Task{
				Core: tas.Core,
				Meta: &task.Meta{
					"test.api.io/key": "bar",
				},
				Node: &task.Node{
					task.Method: task.MthdAny,
				},
				Sync: &task.Sync{
					task.Paging: "233", // sync state kept after expiry via delete
				},
			}
		}

		{
			if !reflect.DeepEqual(tas, exp) {
				t.Fatalf("\n\n%s\n", cmp.Diff(exp, tas))
			}
		}
	}

	// At some point we can set our paging state back to zero and delete the task
	// entirely.
	{
		tas.Sync.Set(task.Paging, "0")
	}

	{
		err = eon.Delete(tas)
		if err != nil {
			t.Fatal(err)
		}
	}

	{
		lis, err = eon.Lister(engine.All())
		if err != nil {
			t.Fatal(err)
		}
	}

	{
		if len(lis) != 0 {
			t.Fatal("expected", 0, "got", len(lis))
		}
	}
}
