//go:build redis

package engine

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/xh3b4sd/logger"
	"github.com/xh3b4sd/redigo"
	"github.com/xh3b4sd/rescue/task"
)

func Test_Engine_Create(t *testing.T) {
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

	var eon *Engine
	{
		eon = New(Config{
			Logger: logger.Fake(),
			Redigo: red,
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
		if !IsTaskMetaEmpty(err) {
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

	var lis []*task.Task
	{
		lis, err = eon.Lister(All())
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
			tas = lis[1]
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
			tas = lis[2]
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
		lis, err = eon.Lister(All())
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
			tas = lis[1]
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
		lis, err = eon.Lister(All())
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
		if !IsTaskNotFound(err) {
			t.Fatal("expected", taskNotFoundError, "got", err)
		}
	}
}

func Test_Engine_Create_Cron(t *testing.T) {
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

	var eon *Engine
	{
		eon = New(Config{
			Logger: logger.Fake(),
			Redigo: red,
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
			Cron: &task.Cron{
				task.Aevery: "hour",
			},
		}

		err = eon.Create(tas)
		if !IsTaskMetaEmpty(err) {
			t.Fatal(err)
		}
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
		lis, err = eon.Lister(All())
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
		err = eon.Delete(tas)
		if err != nil {
			t.Fatal(err)
		}
	}

	// The template for scheduled tasks cannot be returned as task for workers to
	// process. The only way to find them is through Engine.Lister.
	{
		_, err = eon.Search()
		if !IsTaskNotFound(err) {
			t.Fatal("expected", taskNotFoundError, "got", err)
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
		err = eon.Delete(tas)
		if !IsTaskOutdated(err) {
			t.Fatal("task must be deleted by owner")
		}
	}

	// Templates for scheduled tasks can only be deleted when bypassing the
	// internal ownership checks.
	{
		tas.Core.Set().Bypass(true)
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

	var eon *Engine
	{
		eon = New(Config{
			Logger: logger.Fake(),
			Redigo: red,
			Worker: "eon",
		})
	}

	var etw *Engine
	{
		etw = New(Config{
			Logger: logger.Fake(),
			Redigo: red,
			Worker: "etw",
		})
	}

	var eth *Engine
	{
		eth = New(Config{
			Logger: logger.Fake(),
			Redigo: red,
			Worker: "eth",
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
		if !IsTaskNotFound(err) {
			t.Fatal("expected", taskNotFoundError, "got", err)
		}
	}
}

func Test_Engine_Create_Root_First(t *testing.T) {
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

	var eon *Engine
	{
		eon = New(Config{
			Logger: logger.Fake(),
			Redigo: red,
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
		if !IsTaskMetaEmpty(err) {
			t.Fatal("expected task creation to fail without Task.Meta")
		}
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
		lis, err = eon.Lister(All())
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
		lis, err = eon.Lister(All())
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
		lis, err = eon.Lister(All())
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
		err = eon.Delete(tas)
		if err != nil {
			t.Fatal(err)
		}
	}

	{
		tas, err = eon.Search()
		if !IsTaskNotFound(err) {
			t.Fatal("expected", taskNotFoundError, "got", err)
		}
	}
}
