//go:build redis

package ticker

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

func Test_Engine_Ticker_Cron_All(t *testing.T) {
	var err error

	var red redigo.Interface
	{
		red = redigo.Default()
	}

	{
		err = red.Purge()
		if err != nil {
			t.Fatal(err)
		}
	}

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
			Redigo: red,
			Timer:  tim,
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
			Node: &task.Node{
				task.Method: task.MthdAll,
			},
		}

		err = eon.Create(tas)
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
	// exist and have its ticks defined according to the current time.
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
					task.Method: task.MthdAll,
				},
			}
		}

		{
			if !reflect.DeepEqual(tas, exp) {
				t.Fatalf("\n\n%s\n", cmp.Diff(exp, tas))
			}
		}
	}

	// Time advances 1 hour.
	{
		tim.Setter(func() time.Time {
			return musTim("2023-10-20T01:00:01Z")
		})
	}

	{
		err = eon.Ticker()
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
		if len(lis) != 2 {
			t.Fatal("expected", 2, "got", len(lis))
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
					task.Aevery: "hour",
					task.TickM1: "2023-10-20T01:00:00Z",
					task.TickP1: "2023-10-20T02:00:00Z",
				},
				Meta: &task.Meta{
					"test.api.io/key": "foo",
				},
				Node: &task.Node{
					task.Method: task.MthdAll,
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
		var tas *task.Task
		{
			tas = lis[1]
		}

		var exp *task.Task
		{
			exp = &task.Task{
				Core: tas.Core,
				Meta: &task.Meta{
					"test.api.io/key": "foo",
				},
				Node: &task.Node{
					task.Method: task.MthdAll,
				},
				Root: &task.Root{
					task.Object: lis[0].Core.Map().Object(),
				},
			}
		}

		{
			if !reflect.DeepEqual(tas, exp) {
				t.Fatalf("\n\n%s\n", cmp.Diff(exp, tas))
			}
		}
	}

	// We remember the object ID of the first task that got scheduled during the
	// test execution. In the next step we will compare it against the scheduled
	// task that we expect to be scheduled after time advanced beyond the upcoming
	// interval.
	var fir string
	{
		fir = lis[1].Core.Map().Object()
	}

	// Time advances 1 more hour.
	{
		tim.Setter(func() time.Time {
			return musTim("2023-10-20T02:00:01Z")
		})
	}

	{
		err = eon.Ticker()
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

	// We still only expect 2 tasks. That is, the task template, and the scheduled
	// task emitted by its task template. The task template emits tasks with
	// delivery method "all". That implies tasks are locally processed and expired
	// on each worker node separately without acknowledging or removing those
	// special tasks. If we do not handle those special tasks properly, they would
	// simply pile up until the retention period hits. Then we would receive 3
	// tasks instead of 2. So what we want to ensure here, is that scheduled tasks
	// get replaced, if they would be kept otherwise.
	{
		if len(lis) != 2 {
			t.Fatal("expected", 2, "got", len(lis))
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
					task.Aevery: "hour",
					task.TickM1: "2023-10-20T02:00:00Z",
					task.TickP1: "2023-10-20T03:00:00Z",
				},
				Meta: &task.Meta{
					"test.api.io/key": "foo",
				},
				Node: &task.Node{
					task.Method: task.MthdAll,
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
		var tas *task.Task
		{
			tas = lis[1]
		}

		var exp *task.Task
		{
			exp = &task.Task{
				Core: tas.Core,
				Meta: &task.Meta{
					"test.api.io/key": "foo",
				},
				Node: &task.Node{
					task.Method: task.MthdAll,
				},
				Root: &task.Root{
					task.Object: lis[0].Core.Map().Object(),
				},
			}
		}

		{
			if !reflect.DeepEqual(tas, exp) {
				t.Fatalf("\n\n%s\n", cmp.Diff(exp, tas))
			}
		}
	}

	// We track the object ID of the second task that we expect to be scheduled,
	// so we can compare the first and the second object ID below.
	var sec string
	{
		sec = lis[1].Core.Map().Object()
	}

	// If the first and the second ID we tracked are equal, then no new task got
	// emitted by the task template. The task template is supposed to emit a new
	// task after each interval got crossed. So if the second scheduled task is in
	// fact the first scheduled task, then the processing of scheduled tasks
	// defining delivery method "all" is broken.
	if fir == sec {
		t.Fatal("expected", false, "got", true)
	}
}

func Test_Engine_Ticker_Cron_Uni(t *testing.T) {
	var err error

	var red redigo.Interface
	{
		red = redigo.Default()
	}

	{
		err = red.Purge()
		if err != nil {
			t.Fatal(err)
		}
	}

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
			Redigo: red,
			Timer:  tim,
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
			Node: &task.Node{
				task.Method: task.MthdUni,
				task.Worker: "etw",
			},
		}

		err = eon.Create(tas)
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
	// exist and have its ticks defined according to the current time.
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
					task.Method: task.MthdUni,
					task.Worker: "etw",
				},
			}
		}

		{
			if !reflect.DeepEqual(tas, exp) {
				t.Fatalf("\n\n%s\n", cmp.Diff(exp, tas))
			}
		}
	}

	// Time advances 1 hour.
	{
		tim.Setter(func() time.Time {
			return musTim("2023-10-20T01:00:01Z")
		})
	}

	{
		err = eon.Ticker()
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
		if len(lis) != 2 {
			t.Fatal("expected", 2, "got", len(lis))
		}
	}

	// Note that in this test the addressed worker "etw" does not process its
	// designated task. And so tick-m will stay unchanged, indicating the initial
	// or last completed scheduling time.
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
					task.TickP1: "2023-10-20T02:00:00Z",
				},
				Meta: &task.Meta{
					"test.api.io/key": "foo",
				},
				Node: &task.Node{
					task.Method: task.MthdUni,
					task.Worker: "etw",
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
		var tas *task.Task
		{
			tas = lis[1]
		}

		var exp *task.Task
		{
			exp = &task.Task{
				Core: tas.Core,
				Meta: &task.Meta{
					"test.api.io/key": "foo",
				},
				Node: &task.Node{
					task.Method: task.MthdUni,
					task.Worker: "etw",
				},
				Root: &task.Root{
					task.Object: lis[0].Core.Map().Object(),
				},
			}
		}

		{
			if !reflect.DeepEqual(tas, exp) {
				t.Fatalf("\n\n%s\n", cmp.Diff(exp, tas))
			}
		}
	}

	// We remember the object ID of the first task that got scheduled during the
	// test execution. In the next step we will compare it against the scheduled
	// task that we expect to be scheduled after time advanced beyond the upcoming
	// interval.
	var fir string
	{
		fir = lis[1].Core.Map().Object()
	}

	// Time advances 1 more hour.
	{
		tim.Setter(func() time.Time {
			return musTim("2023-10-20T02:00:01Z")
		})
	}

	{
		err = eon.Ticker()
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

	// We still only expect 2 tasks. That is, the task template, and the scheduled
	// task emitted by its task template. The task template emits tasks with
	// delivery method "all". That implies tasks are locally processed and expired
	// on each worker node separately without acknowledging or removing those
	// special tasks. If we do not handle those special tasks properly, they would
	// simply pile up until the retention period hits. Then we would receive 3
	// tasks instead of 2. So what we want to ensure here, is that scheduled tasks
	// get replaced, if they would be kept otherwise.
	{
		if len(lis) != 2 {
			t.Fatal("expected", 2, "got", len(lis))
		}
	}

	// Note that in this test the addressed worker "etw" does not process its
	// designated task. And so tick-m will stay unchanged, indicating the initial
	// or last completed scheduling time.
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
					task.TickP1: "2023-10-20T03:00:00Z",
				},
				Meta: &task.Meta{
					"test.api.io/key": "foo",
				},
				Node: &task.Node{
					task.Method: task.MthdUni,
					task.Worker: "etw",
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
		var tas *task.Task
		{
			tas = lis[1]
		}

		var exp *task.Task
		{
			exp = &task.Task{
				Core: tas.Core,
				Meta: &task.Meta{
					"test.api.io/key": "foo",
				},
				Node: &task.Node{
					task.Method: task.MthdUni,
					task.Worker: "etw",
				},
				Root: &task.Root{
					task.Object: lis[0].Core.Map().Object(),
				},
			}
		}

		{
			if !reflect.DeepEqual(tas, exp) {
				t.Fatalf("\n\n%s\n", cmp.Diff(exp, tas))
			}
		}
	}

	// We track the object ID of the second task that we expect to be scheduled,
	// so we can compare the first and the second object ID below.
	var sec string
	{
		sec = lis[1].Core.Map().Object()
	}

	// If the first and the second ID we tracked are equal, then no new task got
	// emitted by the task template. The task template is supposed to emit a new
	// task after each interval got crossed. So if the second scheduled task is in
	// fact the first scheduled task, then the processing of scheduled tasks
	// defining delivery method "all" is broken.
	if fir == sec {
		t.Fatal("expected", false, "got", true)
	}
}
