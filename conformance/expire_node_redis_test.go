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

func Test_Engine_Expire_Node_All(t *testing.T) {
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

	// The engines is configured with a particular time. This point in time will
	// be set inside each engine as the pointer for when they started processing
	// tasks.
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
			Worker: "eon",
		})
	}

	var etw rescue.Interface
	{
		etw = engine.New(engine.Config{
			Logger: logger.Fake(),
			Redigo: red,
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

	var tas *task.Task
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

	// Time advances by 31 more second. The two workers within the network
	// received the broadcasted task, but they did not complete it by calling
	// Engine.Delete.
	{
		tim.Setter(func() time.Time {
			return musTim("2023-10-20T00:00:32Z")
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

	// Engine one completes the task now.
	{
		err = eon.Delete(tas)
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
		_, err = etw.Search()
		if !engine.IsTaskNotFound(err) {
			t.Fatal("expected", "taskNotFoundError", "got", err)
		}
	}

	// Time advances by 31 more second. One worker completed the broadcasted task,
	// the other didn't.
	{
		tim.Setter(func() time.Time {
			return musTim("2023-10-20T00:01:03Z")
		})
	}

	// Since engine one completed its task and no new tasks got created it should
	// not find any more tasks.
	{
		_, err = eon.Search()
		if !engine.IsTaskNotFound(err) {
			t.Fatal("expected", "taskNotFoundError", "got", err)
		}
	}

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

	// Engine two finally completes the broadcasted task.
	{
		err = etw.Delete(tas)
		if err != nil {
			t.Fatal(err)
		}
	}

	// From this point forwad no worker in the network should find a task anymore.
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

	// All workers processed the broadcasted task. It will still remain in the
	// underyling queue for 1 week. Then it will get deleted eventually.
	{
		var tas *task.Task
		{
			tas = lis[0]
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
			}
		}

		{
			if !reflect.DeepEqual(tas, exp) {
				t.Fatalf("\n\n%s\n", cmp.Diff(exp, tas))
			}
		}
	}

	// Time advances some more.
	{
		tim.Setter(func() time.Time {
			return musTim("2023-10-20T00:05:00Z")
		})
	}

	// Just repeating all lookups should yield the very same results down below.
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

	// All workers processed the broadcasted task. It will still remain in the
	// underyling queue for 1 week. Then it will get deleted eventually.
	{
		var tas *task.Task
		{
			tas = lis[0]
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
			}
		}

		{
			if !reflect.DeepEqual(tas, exp) {
				t.Fatalf("\n\n%s\n", cmp.Diff(exp, tas))
			}
		}
	}
}

func Test_Engine_Expire_Node_Any(t *testing.T) {
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

	var eon rescue.Interface
	{
		eon = engine.New(engine.Config{
			Expiry: time.Millisecond,
			Logger: logger.Fake(),
			Redigo: red,
			Worker: "eon",
		})
	}

	var etw rescue.Interface
	{
		etw = engine.New(engine.Config{
			Expiry: time.Millisecond,
			Logger: logger.Fake(),
			Redigo: red,
			Worker: "etw",
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
		tas := &task.Task{
			Meta: &task.Meta{
				"test.api.io/key": "bar",
			},
		}

		err = etw.Create(tas)
		if err != nil {
			t.Fatal(err)
		}
	}

	var s []string

	{
		_, err = eon.Search()
		if err != nil {
			t.Fatal(err)
		}
	}

	{
		tas, err := etw.Search()
		if err != nil {
			t.Fatal(err)
		}

		s = append(s, tas.Meta.Get("test.api.io/key"))

		err = etw.Delete(tas)
		if err != nil {
			t.Fatal(err)
		}
	}

	// For engine one we simulate failure so that the acquired task can expire and
	// be rescheduled to engine two. For the simulation we call Expire which is
	// the responsibility of every worker to do periodically. It does not matter
	// which engine executes the expiration process.
	{
		err = eon.Expire()
		if err != nil {
			t.Fatal(err)
		}
	}

	{
		tas, err := etw.Search()
		if err != nil {
			t.Fatal(err)
		}

		s = append(s, tas.Meta.Get("test.api.io/key"))

		err = etw.Delete(tas)
		if err != nil {
			t.Fatal(err)
		}
	}

	{
		if len(s) != 2 {
			t.Fatal("length of s must be 2")
		}
		if s[0] == "foo" && s[1] != "bar" {
			t.Fatal("scheduling failed")
		}
		if s[1] == "foo" && s[0] != "bar" {
			t.Fatal("scheduling failed")
		}
	}

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
}

func Test_Engine_Expire_Node_Uni_Cleanup(t *testing.T) {
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

	// The engine is configured with a particular time. This point in time will be
	// set inside the worker process as the pointer for when it started processing
	// tasks.
	{
		tim.Setter(func() time.Time {
			return musTim("2023-10-20T00:00:00Z")
		})
	}

	var etw rescue.Interface
	{
		etw = engine.New(engine.Config{
			Logger: logger.Fake(),
			Redigo: red,
			Timer:  tim,
			Worker: "etw",
		})
	}

	// Worker two creates a task for worker one, which never shows up.
	{
		tas := &task.Task{
			Meta: &task.Meta{
				"test.api.io/key": "bar",
			},
			Node: &task.Node{
				task.Method: task.MthdUni,
				task.Worker: "eon",
			},
		}

		err = etw.Create(tas)
		if err != nil {
			t.Fatal(err)
		}
	}

	// Engine two can never receive any tasks.
	{
		_, err = etw.Search()
		if !engine.IsTaskNotFound(err) {
			t.Fatal("expected", "taskNotFoundError", "got", err)
		}
	}

	var lis []*task.Task
	{
		lis, err = etw.Lister(engine.All())
		if err != nil {
			t.Fatal(err)
		}
	}

	// There should be one task in the queue, since one got just created.
	{
		if len(lis) != 1 {
			t.Fatal("expected", 1, "got", len(lis))
		}
	}

	// Time advances 2 minutes. The task should still exist.
	{
		tim.Setter(func() time.Time {
			return musTim("2023-10-20T00:02:00Z")
		})
	}

	// Engine two can never receive any tasks.
	{
		_, err = etw.Search()
		if !engine.IsTaskNotFound(err) {
			t.Fatal("expected", "taskNotFoundError", "got", err)
		}
	}

	// Engine two can never receive any tasks, but calling Engine.Expire purges
	// any lingering task, regardless which engine executes it.
	{
		err = etw.Expire()
		if err != nil {
			t.Fatal(err)
		}
	}

	{
		lis, err = etw.Lister(engine.All())
		if err != nil {
			t.Fatal(err)
		}
	}

	{
		if len(lis) != 1 {
			t.Fatal("expected", 1, "got", len(lis))
		}
	}

	// Time advances 7 days. The task should be gone now.
	{
		tim.Setter(func() time.Time {
			return musTim("2023-10-27T00:02:00Z")
		})
	}

	// Engine two can never receive any tasks, but calling Engine.Expire purges
	// any lingering task, regardless which engine executes it.
	{
		err = etw.Expire()
		if err != nil {
			t.Fatal(err)
		}
	}

	{
		lis, err = etw.Lister(engine.All())
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

func Test_Engine_Expire_Node_Uni_Lifecycle(t *testing.T) {
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

	// The engine is configured with a particular time. This point in time will be
	// set inside the worker process as the pointer for when it started processing
	// tasks.
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
			Worker: "eon",
		})
	}

	var etw rescue.Interface
	{
		etw = engine.New(engine.Config{
			Logger: logger.Fake(),
			Redigo: red,
			Timer:  tim,
			Worker: "etw",
		})
	}

	// Time advances by 1 minute. So the first task "foo" got created at minute
	// one.
	{
		tim.Setter(func() time.Time {
			return musTim("2023-10-20T00:01:00Z")
		})
	}

	{
		tas := &task.Task{
			Meta: &task.Meta{
				"test.api.io/key": "foo",
			},
			Node: &task.Node{
				task.Method: task.MthdAny,
			},
		}

		err = eon.Create(tas)
		if err != nil {
			t.Fatal(err)
		}
	}

	// Time advances by 1 more minute. So the second task "bar" got created at
	// minute two.
	{
		tim.Setter(func() time.Time {
			return musTim("2023-10-20T00:02:00Z")
		})
	}

	{
		tas := &task.Task{
			Meta: &task.Meta{
				"test.api.io/key": "bar",
			},
			Node: &task.Node{
				task.Method: task.MthdUni,
				task.Worker: "eon",
			},
		}

		err = eon.Create(tas)
		if err != nil {
			t.Fatal(err)
		}
	}

	// Worker one looks for new tasks and finds task two first, because it is
	// directly addressed at worker "eon".
	var tas *task.Task
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
					task.Method: task.MthdUni,
					task.Worker: "eon",
				},
			}
		}

		{
			if !reflect.DeepEqual(tas, exp) {
				t.Fatalf("\n\n%s\n", cmp.Diff(exp, tas))
			}
		}
	}

	// For the test here we pretend task two gets stuck or fails for whatever
	// reason, it will expire within the underlying queue. After expiry we want to
	// see task two being picked up by worker "eon" again.

	// Worker two looks for a task now and finds task one.
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

	// Time advances some 10 seconds. So task one gets completed without issues.
	{
		tim.Setter(func() time.Time {
			return musTim("2023-10-20T00:02:10Z")
		})
	}

	{
		err = etw.Delete(tas)
		if err != nil {
			t.Fatal(err)
		}
	}

	// Worker "etw" cannot find another task anymore, because from its point of
	// view, the queue is empty now.
	{
		_, err = etw.Search()
		if !engine.IsTaskNotFound(err) {
			t.Fatal("expected", "taskNotFoundError", "got", err)
		}
	}

	// Now the failing task two can still not be received by worker one, because
	// of the "failing" task' expiry.
	{
		tas, err = eon.Search()
		if !engine.IsTaskNotFound(err) {
			t.Fatal("expected", "taskNotFoundError", "got", err)
		}
	}

	// Time advances by another 25 seconds. So task two expired within the
	// underlying queue and should be receivable again, particularly for worker
	// "eon".
	{
		tim.Setter(func() time.Time {
			return musTim("2023-10-20T00:02:35Z")
		})
	}

	// Without running Engine.Expire no task can be expired within the underlying
	// system. So even if the task's expiry is due, the system does not recognize
	// it yet. That means searching for a task will not yield any result.
	{
		_, err = eon.Search()
		if !engine.IsTaskNotFound(err) {
			t.Fatal("expected", "taskNotFoundError", "got", err)
		}
	}

	// Any worker may run Engine.Expire in order to expire tasks within the
	// underlying queue. Here worker "etw" is executing the expiration routine for
	// the task addressed directly at worker "eon" to be expired.
	{
		err = etw.Expire()
		if err != nil {
			t.Fatal(err)
		}
	}

	// After task two got properly expired within the system, it can now be
	// received by worker one.
	{
		tas, err = eon.Search()
		if err != nil {
			t.Fatal(err)
		}
	}

	// Engine two can never receive any more tasks from here.
	{
		_, err = etw.Search()
		if !engine.IsTaskNotFound(err) {
			t.Fatal("expected", "taskNotFoundError", "got", err)
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
					task.Method: task.MthdUni,
					task.Worker: "eon",
				},
			}
		}

		{
			if !reflect.DeepEqual(tas, exp) {
				t.Fatalf("\n\n%s\n", cmp.Diff(exp, tas))
			}
		}
	}

	// Time advances by some 5 more seconds. So this time around the worker "eon"
	// completed task two, which was "failing" earlier.
	{
		tim.Setter(func() time.Time {
			return musTim("2023-10-20T00:02:40Z")
		})
	}

	{
		err = eon.Delete(tas)
		if err != nil {
			t.Fatal(err)
		}
	}

	// Since worker one completed its directly assigned task there should not be
	// any more tasks to be received.
	{
		_, err = eon.Search()
		if !engine.IsTaskNotFound(err) {
			t.Fatal("expected", "taskNotFoundError", "got", err)
		}
	}

	// Engine two can never receive any more tasks from here.
	{
		_, err = etw.Search()
		if !engine.IsTaskNotFound(err) {
			t.Fatal("expected", "taskNotFoundError", "got", err)
		}
	}

	var lis []*task.Task
	{
		lis, err = eon.Lister(engine.All())
		if err != nil {
			t.Fatal(err)
		}
	}

	// There should be no task left in the queue, because all tasks got resolved
	// by the workers within the network.
	{
		if len(lis) != 0 {
			t.Fatal("expected", 0, "got", len(lis))
		}
	}
}
