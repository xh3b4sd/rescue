//go:build redis

package conformance

import (
	"fmt"
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

func Test_Engine_Create_Basic(t *testing.T) {
	var err error

	var eon rescue.Interface
	{
		eon = engine.New(engine.Config{
			Logger: logger.Fake(),
			Redigo: prgAll(redigo.Default()),
		})
	}

	// This is the root task.
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
			Root: &task.Root{
				"test.api.io/key": "rrr",
			},
		}

		err = eon.Create(tas)
		if !engine.IsTaskMetaEmpty(err) {
			t.Fatal("expected task creation to fail without Task.Meta")
		}
	}

	// This is the nested task that should be removed internally after calling
	// Engine.Search.
	{
		tas := &task.Task{
			Meta: &task.Meta{
				"test.api.io/key": "zap",
			},
			Root: &task.Root{
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
			Root: &task.Root{
				"test.api.io/key": "rrr",
			},
		}

		err = eon.Create(tas)
		if err != nil {
			t.Fatal(err)
		}
	}

	var lis task.Slicer
	{
		lis, err = eon.Lister(engine.All())
		if err != nil {
			t.Fatal(err)
		}
	}

	// We expect 3 tasks because 3 tasks got created, even if one task is
	// redundant. The redundant task will be cleaned up below when calling
	// Engine.Search.
	{
		if len(lis) != 3 {
			t.Fatal("expected", 3, "got", len(lis))
		}
	}

	{
		var tas *task.Task
		{
			tas = lis.TaskMeta(&task.Meta{"test.api.io/key": "foo"})[0]
		}

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

	{
		var tas *task.Task
		{
			tas = lis.TaskMeta(&task.Meta{"test.api.io/key": "zap"})[0]
		}

		var exp *task.Task
		{
			exp = &task.Task{
				Core: tas.Core,
				Meta: &task.Meta{
					"test.api.io/key": "zap",
				},
				Node: &task.Node{
					task.Method: task.MthdAny,
				},
				Root: &task.Root{
					"test.api.io/key": "foo",
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
			tas = lis.TaskMeta(&task.Meta{"test.api.io/key": "bar"})[0]
		}

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
				Root: &task.Root{
					"test.api.io/key": "rrr",
				},
			}
		}

		{
			if !reflect.DeepEqual(tas, exp) {
				t.Fatalf("\n\n%s\n", cmp.Diff(exp, tas))
			}
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
		var one *task.Task
		{
			one = &task.Task{
				Core: tas.Core,
				Meta: &task.Meta{
					"test.api.io/key": "foo",
				},
				Node: &task.Node{
					task.Method: task.MthdAny,
				},
			}
		}

		var two *task.Task
		{
			two = &task.Task{
				Core: tas.Core,
				Meta: &task.Meta{
					"test.api.io/key": "bar",
				},
				Node: &task.Node{
					task.Method: task.MthdAny,
				},
				Root: &task.Root{
					"test.api.io/key": "rrr",
				},
			}
		}

		{
			if !reflect.DeepEqual(tas, one) && !reflect.DeepEqual(tas, two) {
				t.Fatal("expected meta key foo or bar")
			}
		}
	}

	{
		lis, err = eon.Lister(engine.All())
		if err != nil {
			t.Fatal(err)
		}
	}

	// We expect 2 tasks even though 3 tasks got created. The user did not delete
	// any task. The system noticed one task was redundant due to the tree
	// structure defined by Task.Root. Calling Engine.Search cleaned up the
	// redundant task, leaving us with 2 valid tasks to process.
	{
		if len(lis) != 2 {
			t.Fatal("expected", 2, "got", len(lis))
		}
	}

	{
		var tas *task.Task
		{
			tas = lis.TaskMeta(&task.Meta{"test.api.io/key": "foo"})[0]
		}

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

	{
		var tas *task.Task
		{
			tas = lis.TaskMeta(&task.Meta{"test.api.io/key": "bar"})[0]
		}

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
				Root: &task.Root{
					"test.api.io/key": "rrr",
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

	var foo *task.Task
	{
		foo = &task.Task{
			Core: tas.Core,
			Meta: &task.Meta{
				"test.api.io/key": "foo",
			},
			Node: &task.Node{
				task.Method: task.MthdAny,
			},
		}
	}

	var bar *task.Task
	{
		bar = &task.Task{
			Core: tas.Core,
			Meta: &task.Meta{
				"test.api.io/key": "bar",
			},
			Node: &task.Node{
				task.Method: task.MthdAny,
			},
			Root: &task.Root{
				"test.api.io/key": "rrr",
			},
		}
	}

	var fnd string
	{
		{
			if reflect.DeepEqual(tas, foo) {
				fnd = "foo"
			}
			if reflect.DeepEqual(tas, bar) {
				fnd = "bar"
			}
		}

		{
			if fnd == "" {
				t.Fatal("expected meta key foo or bar")
			}
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

	{
		var tas *task.Task
		{
			tas = lis[0]
		}

		{
			if fnd == "foo" && !reflect.DeepEqual(tas, foo) {
				t.Fatalf("\n\n%s\n", cmp.Diff(foo, tas))
			}
			if fnd == "bar" && !reflect.DeepEqual(tas, bar) {
				t.Fatalf("\n\n%s\n", cmp.Diff(bar, tas))
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
		if !engine.IsTaskNotFound(err) {
			t.Fatal("expected", "taskNotFoundError", "got", err)
		}
	}
}

func Test_Engine_Create_Cron(t *testing.T) {
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

	{
		tas := &task.Task{
			Cron: &task.Cron{
				task.Aevery: "hour",
			},
		}

		err = eon.Create(tas)
		if !engine.IsTaskMetaEmpty(err) {
			t.Fatal(err)
		}
	}

	{
		tim.Setter(func() time.Time {
			return musTim("2023-10-20T00:00:03Z")
		})
	}

	{
		tas := &task.Task{
			Cron: &task.Cron{
				task.Aevery: "hour",
			},
			Meta: &task.Meta{
				"test.api.io/key": "bar",
			},
			Gate: &task.Gate{
				"test.api.io/k-0": task.Trigger,
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
		if len(lis) != 2 {
			t.Fatal("expected", 2, "got", len(lis))
		}

		if lis[0].Meta.Get("test.api.io/key") != "foo" {
			t.Fatal("scheduling failed")
		}
		if lis[0].Cron != nil {
			t.Fatal("expected", nil, "got", lis[0].Cron)
		}
		if lis[0].Gate != nil {
			t.Fatal("expected", nil, "got", lis[0].Gate)
		}

		if lis[1].Meta.Get("test.api.io/key") != "bar" {
			t.Fatal("expected", "bar", "got", lis[1].Meta.Get("test.api.io/key"))
		}
		if lis[1].Cron.Get().Aevery() != "hour" {
			t.Fatal("expected", "hour", "got", lis[1].Cron.Get().Aevery())
		}
		if lis[1].Gate.Get("test.api.io/k-0") != task.Trigger {
			t.Fatal("expected", task.Trigger, "got", lis[1].Gate.Get("test.api.io/k-0"))
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
		if tas.Meta.Get("test.api.io/key") != "foo" {
			t.Fatal("scheduling failed")
		}
		if tas.Root != nil {
			t.Fatal("scheduling failed")
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

	// The template for scheduled tasks cannot be returned as task for workers to
	// process. The only way to find them is through Engine.Lister.
	{
		_, err = eon.Search()
		if !engine.IsTaskNotFound(err) {
			t.Fatal("expected", "taskNotFoundError", "got", err)
		}
	}

	{
		tas = &task.Task{
			Cron: &task.Cron{
				task.Aevery: "hour",
			},
		}
	}

	{
		lis, err = eon.Lister(tas)
		if err != nil {
			t.Fatal(err)
		}
	}

	{
		if len(lis) != 1 {
			t.Fatal("expected 1 task listed")
		}
		if lis[0].Meta.Get("test.api.io/key") != "bar" {
			t.Fatal("scheduling failed")
		}
		if lis[0].Cron.Get().Aevery() != "hour" {
			t.Fatal("scheduling failed")
		}
	}

	{
		tas.Core = &task.Core{}
		tas.Core.Set().Object(lis[0].Core.Get().Object())
	}

	{
		tim.Setter(func() time.Time {
			return musTim("2023-10-20T00:00:05Z")
		})
	}

	{
		err = eon.Delete(tas)
		if !engine.IsTaskOutdated(err) {
			t.Fatal("task must be deleted by owner")
		}
	}

	// Templates for scheduled tasks can only be deleted when bypassing the
	// internal ownership checks.
	{
		tas.Core.Set().Bypass(true)
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
		lis, err = eon.Lister(&task.Task{Cron: tas.Cron.All(task.Aevery)})
		if err != nil {
			t.Fatal(err)
		}
	}

	{
		if len(lis) != 0 {
			t.Fatal("expected 0 tasks listed")
		}
	}
}

func Test_Engine_Create_Node_All(t *testing.T) {
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

	var eth rescue.Interface
	{
		eth = engine.New(engine.Config{
			Logger: logger.Fake(),
			Redigo: prgAll(redigo.Default()),
			Timer:  tim,
			Worker: "eth",
		})
	}

	{
		tim.Setter(func() time.Time {
			return musTim("2023-10-20T00:00:01Z")
		})
	}

	// The first task we create does not have a delivery method defined, so it
	// should default to "any".
	{
		tas := &task.Task{
			Meta: &task.Meta{
				"test.api.io/key": "fir",
			},
		}

		err = etw.Create(tas)
		if err != nil {
			t.Fatal(err)
		}
	}

	{
		tim.Setter(func() time.Time {
			return musTim("2023-10-20T00:00:02Z")
		})
	}

	// The second task we create defines the delivery method "all", so it should
	// be the first task every worker receives, even if it was not first in line
	// of creation.
	{
		tas := &task.Task{
			Meta: &task.Meta{
				"test.api.io/key": "sec",
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

	// Ensure engine one receives the task defining delivery method "all" first.
	{
		var exp *task.Task
		{
			exp = &task.Task{
				Core: tas.Core,
				Meta: &task.Meta{
					"test.api.io/key": "sec",
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
			if tas.Core.Exi().Worker() {
				t.Fatal("expected", false, "got", true)
			}
		}
	}

	{
		tas, err = etw.Search()
		if err != nil {
			t.Fatal(err)
		}
	}

	// Ensure engine two receives the task defining delivery method "all" first.
	{
		var exp *task.Task
		{
			exp = &task.Task{
				Core: tas.Core,
				Meta: &task.Meta{
					"test.api.io/key": "sec",
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
			if tas.Core.Exi().Worker() {
				t.Fatal("expected", false, "got", true)
			}
		}
	}

	{
		tas, err = eth.Search()
		if err != nil {
			t.Fatal(err)
		}
	}

	// Ensure engine three receives the task defining delivery method "all" first.
	{
		var exp *task.Task
		{
			exp = &task.Task{
				Core: tas.Core,
				Meta: &task.Meta{
					"test.api.io/key": "sec",
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
			if tas.Core.Exi().Worker() {
				t.Fatal("expected", false, "got", true)
			}
		}
	}

	{
		tas, err = etw.Search()
		if err != nil {
			t.Fatal(err)
		}
	}

	// Let engine two search for tasks again. Ensure it receives the task that we
	// created at first. This is the last task in the queue. Its delivery method
	// should default to "any".
	{
		var exp *task.Task
		{
			exp = &task.Task{
				Core: tas.Core,
				Meta: &task.Meta{
					"test.api.io/key": "fir",
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

	// Any further searches should result in no task being found.
	{
		_, err = eth.Search()
		if !engine.IsTaskNotFound(err) {
			t.Fatal("expected", "taskNotFoundError", "got", err)
		}
	}
}

func Test_Engine_Create_Root_First(t *testing.T) {
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

	{
		tim.Setter(func() time.Time {
			return musTim("2023-10-20T00:00:01Z")
		})
	}

	// This is the nested task that should be removed internally after calling
	// Engine.Search.
	{
		tas := &task.Task{
			Meta: &task.Meta{
				"test.api.io/key": "zap",
			},
			Root: &task.Root{
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
			Root: &task.Root{
				"test.api.io/key": "rrr",
			},
		}

		err = eon.Create(tas)
		if !engine.IsTaskMetaEmpty(err) {
			t.Fatal("expected task creation to fail without Task.Meta")
		}
	}

	{
		tim.Setter(func() time.Time {
			return musTim("2023-10-20T00:00:02Z")
		})
	}

	// This is the root task.
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
			return musTim("2023-10-20T00:00:03Z")
		})
	}

	{
		tas := &task.Task{
			Meta: &task.Meta{
				"test.api.io/key": "bar",
			},
			Root: &task.Root{
				"test.api.io/key": "rrr",
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
		if len(lis) != 3 {
			t.Fatal("expected 3 tasks listed")
		}
		if lis[0].Meta.Get("test.api.io/key") != "zap" {
			t.Fatal("scheduling failed")
		}
		if lis[0].Root.Get("test.api.io/key") != "foo" {
			t.Fatal("scheduling failed")
		}
		if lis[1].Meta.Get("test.api.io/key") != "foo" {
			t.Fatal("scheduling failed")
		}
		if lis[1].Root != nil {
			t.Fatal("scheduling failed")
		}
		if lis[2].Meta.Get("test.api.io/key") != "bar" {
			t.Fatal("scheduling failed")
		}
		if lis[2].Root.Get("test.api.io/key") != "rrr" {
			t.Fatal("scheduling failed")
		}
	}

	var tas *task.Task
	{
		tas, err = eon.Search()
		if err != nil {
			fmt.Printf("%#v\n", err)
			t.Fatal(err)
		}
	}

	{
		if tas.Meta.Get("test.api.io/key") != "foo" {
			t.Fatal("scheduling failed")
		}
		if tas.Root != nil {
			t.Fatal("scheduling failed")
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
			t.Fatal("expected 2 tasks listed")
		}
		if lis[0].Meta.Get("test.api.io/key") != "foo" {
			t.Fatal("scheduling failed")
		}
		if lis[0].Root != nil {
			t.Fatal("scheduling failed")
		}
		if lis[1].Meta.Get("test.api.io/key") != "bar" {
			t.Fatal("scheduling failed")
		}
		if lis[1].Root.Get("test.api.io/key") != "rrr" {
			t.Fatal("scheduling failed")
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
		if tas.Meta.Get("test.api.io/key") != "bar" {
			t.Fatal("scheduling failed")
		}
		if tas.Root.Get("test.api.io/key") != "rrr" {
			t.Fatal("scheduling failed")
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
			t.Fatal("expected 1 task listed")
		}
		if lis[0].Meta.Get("test.api.io/key") != "bar" {
			t.Fatal("scheduling failed")
		}
		if lis[0].Root.Get("test.api.io/key") != "rrr" {
			t.Fatal("scheduling failed")
		}
	}

	{
		tim.Setter(func() time.Time {
			return musTim("2023-10-20T00:00:05Z")
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
		if !engine.IsTaskNotFound(err) {
			t.Fatal("expected", "taskNotFoundError", "got", err)
		}
	}
}
