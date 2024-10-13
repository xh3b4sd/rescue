//go:build redis

package conformance

import (
	"reflect"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/xh3b4sd/logger"
	"github.com/xh3b4sd/redigo"
	"github.com/xh3b4sd/rescue"
	"github.com/xh3b4sd/rescue/engine"
	"github.com/xh3b4sd/rescue/task"
	"github.com/xh3b4sd/rescue/timer"
)

func Test_Engine_Delete_Sync(t *testing.T) {
	var err error

	var tim *timer.Timer
	{
		tim = timer.New()
	}

	{
		tim.Setter(func() time.Time {
			return musTim("2023-10-20T00:00:00Z")
		})
	}

	var eon rescue.Interface
	{
		eon = engine.New(engine.Config{
			Logger: logger.Fake(),
			Redigo: prgAll(redigo.Default()),
			Timer:  tim,
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

	{
		tim.Setter(func() time.Time {
			return musTim("2023-10-20T00:00:01Z")
		})
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

	{
		tim.Setter(func() time.Time {
			return musTim("2023-10-20T00:00:02Z")
		})
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
		tim.Setter(func() time.Time {
			return musTim("2023-10-20T00:00:03Z")
		})
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
		tim.Setter(func() time.Time {
			return musTim("2023-10-20T00:00:04Z")
		})
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

	{
		tim.Setter(func() time.Time {
			return musTim("2023-10-20T00:00:05Z")
		})
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
		tim.Setter(func() time.Time {
			return musTim("2023-10-20T00:00:06Z")
		})
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
