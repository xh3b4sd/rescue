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

func Test_Engine_Expire_Cron_Adefer(t *testing.T) {
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
			Worker: "eon",
		})
	}

	var etw rescue.Interface
	{
		etw = engine.New(engine.Config{
			Logger: logger.Fake(),
			Redigo: prgAll(redigo.Default()),
			Timer:  tim,
			Worker: "etw",
		})
	}

	{
		tim.Setter(func() time.Time {
			return musTim("2023-10-20T00:00:01Z")
		})
	}

	{
		tas := &task.Task{
			Meta: &task.Meta{
				"test.api.io/key": "foo",
			},
		}

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

	var tas *task.Task
	{
		tas, err = eon.Search()
		if err != nil {
			t.Fatal(err)
		}
	}

	if tas.Cron == nil {
		tas.Cron = &task.Cron{}
	}

	{
		tas.Cron.Set().Adefer("2 hours")
	}

	// plus 1 minute
	{
		tim.Setter(func() time.Time {
			return musTim("2023-10-20T00:01:00Z")
		})
	}

	// In order to persist the @defer statement, the task must be updated inside
	// of Engine.Delete.
	{
		err = etw.Delete(tas)
		if err != nil {
			t.Fatal(err)
		}
	}

	var lis []*task.Task
	{
		lis, err = eon.Lister(engine.All())
		if err != nil {
			t.Fatal(err)
		}
	}

	{
		if len(lis) != 1 {
			t.Fatal("expected", 1, "got", len(lis))
		}
	}

	{
		var tas *task.Task
		{
			tas = lis[0]
		}

		var exp *task.Task
		{
			exp = &task.Task{
				Core: tas.Core,
				Cron: &task.Cron{
					task.Adefer: "2 hours",
					task.TickP1: "2023-10-20T02:01:00Z", // plus 2 hours
				},
				Meta: &task.Meta{
					"test.api.io/key": "foo",
				},
				Node: &task.Node{
					task.Method: task.MthdAny,
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
		_, err = eon.Search()
		if !engine.IsTaskNotFound(err) {
			t.Fatal("expected", "taskNotFoundError", "got", err)
		}
	}

	{
		tim.Setter(func() time.Time {
			return musTim("2023-10-20T01:00:00Z")
		})
	}

	{
		_, err = eon.Search()
		if !engine.IsTaskNotFound(err) {
			t.Fatal("expected", "taskNotFoundError", "got", err)
		}
	}

	{
		tim.Setter(func() time.Time {
			return musTim("2023-10-20T02:00:59Z")
		})
	}

	{
		lis, err = eon.Lister(engine.All())
		if err != nil {
			t.Fatal(err)
		}
	}

	{
		if len(lis) != 1 {
			t.Fatal("expected", 1, "got", len(lis))
		}
	}

	{
		var tas *task.Task
		{
			tas = lis[0]
		}

		var exp *task.Task
		{
			exp = &task.Task{
				Core: tas.Core,
				Cron: &task.Cron{
					task.Adefer: "2 hours",
					task.TickP1: "2023-10-20T02:01:00Z",
				},
				Meta: &task.Meta{
					"test.api.io/key": "foo",
				},
				Node: &task.Node{
					task.Method: task.MthdAny,
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
		_, err = eon.Search()
		if !engine.IsTaskNotFound(err) {
			t.Fatal("expected", "taskNotFoundError", "got", err)
		}
	}

	{
		tim.Setter(func() time.Time {
			return musTim("2023-10-20T02:01:00Z") // at defer
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
				Cron: &task.Cron{
					task.Adefer: "2 hours",
					task.TickP1: "2023-10-20T02:01:00Z",
				},
				Meta: &task.Meta{
					"test.api.io/key": "foo",
				},
				Node: &task.Node{
					task.Method: task.MthdAny,
				},
			}
		}

		{
			if !reflect.DeepEqual(tas, exp) {
				t.Fatalf("\n\n%s\n", cmp.Diff(exp, tas))
			}
		}
	}

	// do some work
	{
		tim.Setter(func() time.Time {
			return musTim("2023-10-20T02:01:07Z")
		})
	}

	{
		err = etw.Delete(tas)
		if err != nil {
			t.Fatal(err)
		}
	}

	{
		_, err = eon.Search()
		if !engine.IsTaskNotFound(err) {
			t.Fatal("expected", "taskNotFoundError", "got", err)
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

func Test_Engine_Expire_Cron_Aevery(t *testing.T) {
	var err error

	var tim *timer.Timer
	{
		tim = timer.New()
	}

	// The engines are configured with a particular time. This point in time will
	// be set inside each worker process as the pointer for when they started
	// processing tasks.
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
			Worker: "eon",
		})
	}

	var etw rescue.Interface
	{
		etw = engine.New(engine.Config{
			Logger: logger.Fake(),
			Redigo: prgAll(redigo.Default()),
			Timer:  tim,
			Worker: "etw",
		})
	}

	// Time advances by 1 second. So one second after the workers started
	// participating in the network, the task below got created.
	{
		tim.Setter(func() time.Time {
			return musTim("2023-10-20T00:00:01Z")
		})
	}

	{
		tas := &task.Task{
			Cron: &task.Cron{
				task.Aevery: "hour",
			},
			Meta: &task.Meta{
				"test.api.io/key": "foo",
			},
		}

		err = eon.Create(tas)
		if err != nil {
			t.Fatal(err)
		}
	}

	// Time advances by 1 more second. So one second after the task got created
	// above, both participating workers search for the task.
	{
		tim.Setter(func() time.Time {
			return musTim("2023-10-20T00:00:02Z")
		})
	}

	// Calling Engine.Expire purges any lingering task, regardless which engine
	// executes it.
	{
		err = etw.Expire()
		if err != nil {
			t.Fatal(err)
		}
	}

	var lis []*task.Task
	{
		lis, err = eon.Lister(engine.All())
		if err != nil {
			t.Fatal(err)
		}
	}

	{
		if len(lis) != 1 {
			t.Fatal("expected", 1, "got", len(lis))
		}
	}

	// Shortly after task creation, the task template defining Task.Cron should
	// still exist, regardless the call to Engine.Expire.
	{
		var tas *task.Task
		{
			tas = lis[0]
		}

		var exp *task.Task
		{
			exp = &task.Task{
				Core: tas.Core,
				Cron: &task.Cron{
					task.Aevery: "hour",
					task.TickM1: "2023-10-20T00:00:00Z",
					task.TickP1: "2023-10-20T01:00:00Z",
				},
				Meta: &task.Meta{
					"test.api.io/key": "foo",
				},
				Node: &task.Node{
					task.Method: task.MthdAny,
				},
			}
		}

		{
			if !reflect.DeepEqual(tas, exp) {
				t.Fatalf("\n\n%s\n", cmp.Diff(exp, tas))
			}
		}
	}

	// Time advances 7 days.
	{
		tim.Setter(func() time.Time {
			return musTim("2023-10-27T00:00:02Z")
		})
	}

	// Calling Engine.Expire purges any lingering task, regardless which engine
	// executes it.
	{
		err = etw.Expire()
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
		if len(lis) != 1 {
			t.Fatal("expected", 1, "got", len(lis))
		}
	}

	// More than one week after task creation, the task template defining
	// Task.Cron should still exist, regardless the call to Engine.Expire. Note
	// that since we do not call Engine.Ticker the ticks in Task.Cron do not
	// advance for this test, which is also not subject of this very test case, so
	// we ignore it.
	{
		var tas *task.Task
		{
			tas = lis[0]
		}

		var exp *task.Task
		{
			exp = &task.Task{
				Core: tas.Core,
				Cron: &task.Cron{
					task.Aevery: "hour",
					task.TickM1: "2023-10-20T00:00:00Z",
					task.TickP1: "2023-10-20T01:00:00Z",
				},
				Meta: &task.Meta{
					"test.api.io/key": "foo",
				},
				Node: &task.Node{
					task.Method: task.MthdAny,
				},
			}
		}

		{
			if !reflect.DeepEqual(tas, exp) {
				t.Fatalf("\n\n%s\n", cmp.Diff(exp, tas))
			}
		}
	}
}

func Test_Engine_Expire_Cron_Aexact(t *testing.T) {
	var err error

	var tim *timer.Timer
	{
		tim = timer.New()
	}

	// The engines are configured with a particular time. This point in time will
	// be set inside each worker process as the pointer for when they started
	// processing tasks.
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
			Worker: "eon",
		})
	}

	var etw rescue.Interface
	{
		etw = engine.New(engine.Config{
			Logger: logger.Fake(),
			Redigo: prgAll(redigo.Default()),
			Timer:  tim,
			Worker: "etw",
		})
	}

	// Time advances by 1 second. So one second after the workers started
	// participating in the network, the task below got created.
	{
		tim.Setter(func() time.Time {
			return musTim("2023-10-20T00:00:01Z")
		})
	}

	// Define the task to be executed in 5 days, not before 15:30.
	{
		tas := &task.Task{
			Cron: &task.Cron{
				task.Aexact: "2023-10-25T15:30:00Z",
			},
			Meta: &task.Meta{
				"test.api.io/key": "foo",
			},
		}

		err = eon.Create(tas)
		if err != nil {
			t.Fatal(err)
		}
	}

	// Time advances by 1 more second. So one second after the task got created
	// above, both participating workers search for the task.
	{
		tim.Setter(func() time.Time {
			return musTim("2023-10-20T00:00:02Z")
		})
	}

	// Calling Engine.Expire purges any lingering task, regardless which engine
	// executes it.
	{
		err = etw.Expire()
		if err != nil {
			t.Fatal(err)
		}
	}

	var lis []*task.Task
	{
		lis, err = eon.Lister(engine.All())
		if err != nil {
			t.Fatal(err)
		}
	}

	{
		if len(lis) != 1 {
			t.Fatal("expected", 1, "got", len(lis))
		}
	}

	// Shortly after task creation, the task template defining Task.Cron should
	// still exist, regardless the call to Engine.Expire.
	{
		var tas *task.Task
		{
			tas = lis[0]
		}

		var exp *task.Task
		{
			exp = &task.Task{
				Core: tas.Core,
				Cron: &task.Cron{
					task.Aexact: "2023-10-25T15:30:00Z",
				},
				Meta: &task.Meta{
					"test.api.io/key": "foo",
				},
				Node: &task.Node{
					task.Method: task.MthdAny,
				},
			}
		}

		{
			if !reflect.DeepEqual(tas, exp) {
				t.Fatalf("\n\n%s\n", cmp.Diff(exp, tas))
			}
		}
	}

	// None of the engines should be able to search for the task just yet.
	{
		_, err = eon.Search()
		if !engine.IsTaskNotFound(err) {
			t.Fatal("expected", "taskNotFoundError", "got", err)
		}
	}

	{
		_, err = etw.Search()
		if !engine.IsTaskNotFound(err) {
			t.Fatal("expected", "taskNotFoundError", "got", err)
		}
	}

	// Time advances 5 days.
	{
		tim.Setter(func() time.Time {
			return musTim("2023-10-25T00:00:00Z")
		})
	}

	// Calling Engine.Expire purges any lingering task, regardless which engine
	// executes it.
	{
		err = etw.Expire()
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
		if len(lis) != 1 {
			t.Fatal("expected", 1, "got", len(lis))
		}
	}

	// More than one week after task creation, the task template defining
	// Task.Cron should still exist, regardless the call to Engine.Expire. Note
	// that since we do not call Engine.Ticker the ticks in Task.Cron do not
	// advance for this test, which is also not subject of this very test case, so
	// we ignore it.
	{
		var tas *task.Task
		{
			tas = lis[0]
		}

		var exp *task.Task
		{
			exp = &task.Task{
				Core: tas.Core,
				Cron: &task.Cron{
					task.Aexact: "2023-10-25T15:30:00Z",
				},
				Meta: &task.Meta{
					"test.api.io/key": "foo",
				},
				Node: &task.Node{
					task.Method: task.MthdAny,
				},
			}
		}

		{
			if !reflect.DeepEqual(tas, exp) {
				t.Fatalf("\n\n%s\n", cmp.Diff(exp, tas))
			}
		}
	}

	// None of the engines should be able to search for the task just yet.
	{
		_, err = eon.Search()
		if !engine.IsTaskNotFound(err) {
			t.Fatal("expected", "taskNotFoundError", "got", err)
		}
	}

	{
		_, err = etw.Search()
		if !engine.IsTaskNotFound(err) {
			t.Fatal("expected", "taskNotFoundError", "got", err)
		}
	}

	// Time advances into the afternoon.
	{
		tim.Setter(func() time.Time {
			return musTim("2023-10-25T15:30:00Z")
		})
	}

	// The scheduled task should be found and executed now.
	var tas *task.Task
	{
		tas, err = etw.Search()
		if err != nil {
			t.Fatal(err)
		}
	}

	{
		var exp *task.Task
		{
			exp = &task.Task{
				Core: tas.Core,
				Cron: &task.Cron{
					task.Aexact: "2023-10-25T15:30:00Z",
				},
				Meta: &task.Meta{
					"test.api.io/key": "foo",
				},
				Node: &task.Node{
					task.Method: task.MthdAny,
				},
			}
		}

		{
			if !reflect.DeepEqual(tas, exp) {
				t.Fatalf("\n\n%s\n", cmp.Diff(exp, tas))
			}
		}
	}

	// After execution the task should be able to be deleted.
	{
		err = etw.Delete(tas)
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
