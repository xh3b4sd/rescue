//go:build redis

package engine

import (
	"slices"
	"testing"
	"time"

	"github.com/xh3b4sd/logger"
	"github.com/xh3b4sd/redigo"
	"github.com/xh3b4sd/rescue/task"
)

func Test_Engine_Exists(t *testing.T) {
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
			Expiry: 500 * time.Millisecond,
			Logger: logger.Fake(),
			Redigo: red,
			Worker: "eon",
		})
	}

	{
		exi, err := eon.Exists(&task.Task{Core: &task.Core{task.Object: "*"}})
		if err != nil {
			t.Fatal(err)
		}

		if exi {
			t.Fatal("task must not exist")
		}
	}

	{
		exi, err := eon.Exists(&task.Task{Meta: &task.Meta{"test.api.io/key": "foo"}})
		if err != nil {
			t.Fatal(err)
		}

		if exi {
			t.Fatal("task must not exist")
		}
	}

	{
		exi, err := eon.Exists(&task.Task{Root: &task.Root{"test.api.io/key": "rrr"}})
		if err != nil {
			t.Fatal(err)
		}

		if exi {
			t.Fatal("task must not exist")
		}
	}

	{
		exi, err := eon.Exists(&task.Task{Root: &task.Root{task.Object: "*"}})
		if err != nil {
			t.Fatal(err)
		}

		if exi {
			t.Fatal("task must not exist")
		}
	}

	{
		exi, err := eon.Exists(&task.Task{Meta: &task.Meta{"test.api.io/key": "foo"}, Root: &task.Root{"test.api.io/key": "rrr"}})
		if err != nil {
			t.Fatal(err)
		}

		if exi {
			t.Fatal("task must not exist")
		}
	}

	{
		_, err := eon.Exists(&task.Task{Root: &task.Root{task.Worker: "*"}})
		if !IsLabelReserved(err) {
			t.Fatal("expected", labelReservedError, "got", err)
		}
	}

	{
		tas := &task.Task{
			Meta: &task.Meta{
				"test.api.io/key": "foo",
				"test.api.io/obj": "one",
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
				"test.api.io/key": "foo",
				"test.api.io/obj": "two",
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

	{
		tas := &task.Task{
			Meta: &task.Meta{
				"test.api.io/key": "baz",
				"test.api.io/obj": "thr",
			},
		}

		err = eon.Create(tas)
		if err != nil {
			t.Fatal(err)
		}
	}

	{
		exi, err := eon.Exists(&task.Task{Meta: &task.Meta{"test.api.io/key": "foo"}})
		if err != nil {
			t.Fatal(err)
		}

		if !exi {
			t.Fatal("task must exist")
		}
	}

	{
		exi, err := eon.Exists(&task.Task{Root: &task.Root{"test.api.io/key": "rrr"}})
		if err != nil {
			t.Fatal(err)
		}

		if !exi {
			t.Fatal("task must exist")
		}
	}

	{
		exi, err := eon.Exists(&task.Task{Meta: &task.Meta{"test.api.io/key": "foo"}, Root: &task.Root{"test.api.io/key": "rrr"}})
		if err != nil {
			t.Fatal(err)
		}

		if !exi {
			t.Fatal("task must exist")
		}
	}

	var tas *task.Task
	{
		tas, err = eon.Search()
		if err != nil {
			t.Fatal(err)
		}

		if tas.Meta.Get("test.api.io/key") != "foo" {
			t.Fatal("scheduling failed")
		}
		if tas.Root != nil {
			t.Fatal("scheduling failed")
		}
	}

	{
		exi, err := eon.Exists(&task.Task{Core: tas.Core.All(task.Object)})
		if err != nil {
			t.Fatal(err)
		}

		if !exi {
			t.Fatal("task must exist")
		}
	}

	{
		exi, err := eon.Exists(&task.Task{Core: tas.Core.All(task.Object, task.Worker)})
		if err != nil {
			t.Fatal(err)
		}

		if !exi {
			t.Fatal("task must be found with owner")
		}
	}

	{
		time.Sleep(500 * time.Millisecond)
	}

	{
		err = eon.Expire()
		if err != nil {
			t.Fatal(err)
		}
	}

	{
		exi, err := eon.Exists(&task.Task{Core: tas.Core.All(task.Object)})
		if err != nil {
			t.Fatal(err)
		}

		if !exi {
			t.Fatal("task must exist")
		}
	}

	{
		exi, err := eon.Exists(&task.Task{Core: tas.Core.All(task.Object, task.Worker)})
		if err != nil {
			t.Fatal(err)
		}

		if exi {
			t.Fatal("task must not be found with owner")
		}
	}

	{
		err = eon.Delete(tas)
		if !IsTaskOutdated(err) {
			t.Fatal("expected", taskOutdatedError, "got", err)
		}
	}

	{
		tas, err = eon.Search()
		if err != nil {
			t.Fatal(err)
		}

		if tas.Meta.Get("test.api.io/key") != "foo" {
			t.Fatal("expected", "foo", "got", tas.Meta.Get("test.api.io/key"))
		}
		if tas.Root != nil {
			t.Fatal("expected", nil, "got", tas.Root)
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

		if tas.Meta.Get("test.api.io/key") != "foo" {
			t.Fatal("expected", "foo", "got", tas.Meta.Get("test.api.io/key"))
		}
		if tas.Root.Get("test.api.io/key") != "rrr" {
			t.Fatal("expected", "rrr", "got", tas.Root.Get("test.api.io/key"))
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

		if tas.Meta.Get("test.api.io/key") != "baz" {
			t.Fatal("expected", "baz", "got", tas.Meta.Get("test.api.io/key"))
		}
		if tas.Root != nil {
			t.Fatal("expected", nil, "got", tas.Root)
		}
	}

	{
		err = eon.Delete(tas)
		if err != nil {
			t.Fatal(err)
		}
	}

	{
		_, err = eon.Search()
		if !IsTaskNotFound(err) {
			t.Fatal("expected", taskNotFoundError, "got", err)
		}
	}
}

func Test_Engine_Exists_Gate(t *testing.T) {
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
		exi, err := eon.Exists(&task.Task{Gate: &task.Gate{"*": task.Deleted}})
		if err != nil {
			t.Fatal(err)
		}

		if exi {
			t.Fatal("expected", false, "got", true)
		}
	}

	{
		exi, err := eon.Exists(&task.Task{Gate: &task.Gate{"*": task.Trigger}})
		if err != nil {
			t.Fatal(err)
		}

		if exi {
			t.Fatal("expected", false, "got", true)
		}
	}

	{
		exi, err := eon.Exists(&task.Task{Gate: &task.Gate{"*": task.Waiting}})
		if err != nil {
			t.Fatal(err)
		}

		if exi {
			t.Fatal("expected", false, "got", true)
		}
	}

	{
		tas := &task.Task{
			Meta: &task.Meta{
				"test.api.io/key": "foo",
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

	{
		tas := &task.Task{
			Meta: &task.Meta{
				"test.api.io/key": "foo",
			},
			Gate: &task.Gate{
				"test.api.io/k-1": task.Trigger,
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
			Gate: &task.Gate{
				"test.api.io/k-0": task.Waiting,
				"test.api.io/k-1": task.Waiting,
			},
		}

		err = eon.Create(tas)
		if err != nil {
			t.Fatal(err)
		}
	}

	{
		exi, err := eon.Exists(&task.Task{Gate: &task.Gate{"*": task.Deleted}})
		if err != nil {
			t.Fatal(err)
		}

		if exi {
			t.Fatal("expected", false, "got", true)
		}
	}

	{
		exi, err := eon.Exists(&task.Task{Gate: &task.Gate{"*": task.Trigger}})
		if err != nil {
			t.Fatal(err)
		}

		if !exi {
			t.Fatal("expected", true, "got", false)
		}
	}

	{
		exi, err := eon.Exists(&task.Task{Gate: &task.Gate{"*": task.Waiting}})
		if err != nil {
			t.Fatal(err)
		}

		if !exi {
			t.Fatal("expected", true, "got", false)
		}
	}

	var tas *task.Task
	{
		tas, err = eon.Search()
		if err != nil {
			t.Fatal(err)
		}

		if tas.Meta.Get("test.api.io/key") != "foo" {
			t.Fatal("expected", "foo", "got", tas.Meta.Get("test.api.io/key"))
		}
	}

	{
		exi, err := eon.Exists(&task.Task{Gate: &task.Gate{"*": task.Deleted}})
		if err != nil {
			t.Fatal(err)
		}

		if exi {
			t.Fatal("expected", false, "got", true)
		}
	}

	{
		exi, err := eon.Exists(&task.Task{Gate: &task.Gate{"*": task.Trigger}})
		if err != nil {
			t.Fatal(err)
		}

		if !exi {
			t.Fatal("expected", true, "got", false)
		}
	}

	{
		exi, err := eon.Exists(&task.Task{Gate: &task.Gate{"*": task.Waiting}})
		if err != nil {
			t.Fatal(err)
		}

		if !exi {
			t.Fatal("expected", true, "got", false)
		}
	}

	{
		err = eon.Delete(tas)
		if err != nil {
			t.Fatal(err)
		}
	}

	{
		exi, err := eon.Exists(&task.Task{Gate: &task.Gate{"*": task.Deleted}})
		if err != nil {
			t.Fatal(err)
		}

		if !exi {
			t.Fatal("expected", true, "got", false)
		}
	}

	{
		exi, err := eon.Exists(&task.Task{Gate: &task.Gate{"*": task.Trigger}})
		if err != nil {
			t.Fatal(err)
		}

		if !exi {
			t.Fatal("expected", true, "got", false)
		}
	}

	{
		exi, err := eon.Exists(&task.Task{Gate: &task.Gate{"*": task.Waiting}})
		if err != nil {
			t.Fatal(err)
		}

		if !exi {
			t.Fatal("expected", true, "got", false)
		}
	}

	{
		tas, err = eon.Search()
		if err != nil {
			t.Fatal(err)
		}

		if tas.Meta.Get("test.api.io/key") != "foo" {
			t.Fatal("expected", "foo", "got", tas.Meta.Get("test.api.io/key"))
		}
	}

	{
		err = eon.Delete(tas)
		if err != nil {
			t.Fatal(err)
		}
	}

	// After the second call to Engine.Delete both trigger labels got set to
	// "deleted" and consequently set to "waiting" again in one go. So no reserved
	// "deleted" value should exist anymore at this point.
	{
		exi, err := eon.Exists(&task.Task{Gate: &task.Gate{"*": task.Deleted}})
		if err != nil {
			t.Fatal(err)
		}

		if exi {
			t.Fatal("expected", false, "got", true)
		}
	}

	{
		exi, err := eon.Exists(&task.Task{Gate: &task.Gate{"*": task.Trigger}})
		if err != nil {
			t.Fatal(err)
		}

		if exi {
			t.Fatal("expected", false, "got", true)
		}
	}

	// After the second call to Engine.Delete both trigger labels got set to
	// "deleted" and consequently set to "waiting" again in one go. So the
	// reserved value "waiting" should exist still at this point.
	{
		exi, err := eon.Exists(&task.Task{Gate: &task.Gate{"*": task.Waiting}})
		if err != nil {
			t.Fatal(err)
		}

		if !exi {
			t.Fatal("expected", true, "got", false)
		}
	}

	// After the second call to Engine.Delete the gating task template got
	// triggered to emit the new scheduled task.
	{
		tas, err = eon.Search()
		if err != nil {
			t.Fatal(err)
		}

		if tas.Meta.Get("test.api.io/key") != "bar" {
			t.Fatal("expected", "bar", "got", tas.Meta.Get("test.api.io/key"))
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

		if lis[0].Meta.Get("test.api.io/key") != "bar" {
			t.Fatal("expected", "bar", "got", lis[0].Meta.Get("test.api.io/key"))
		}
		if lis[0].Gate.Len() != 2 {
			t.Fatal("expected", 2, "got", lis[0].Gate.Len())
		}

		var key []string
		{
			key = lis[0].Gate.Key()
		}

		{
			slices.Sort(key)
		}

		if key[0] != "test.api.io/k-0" {
			t.Fatal("expected", "test.api.io/k-0", "got", key[0])
		}
		if key[1] != "test.api.io/k-1" {
			t.Fatal("expected", "test.api.io/k-1", "got", key[1])
		}
		if lis[0].Gate.Get(key[0]) != task.Waiting {
			t.Fatal("expected", task.Waiting, "got", lis[0].Gate.Get(key[0]))
		}
		if lis[0].Gate.Get(key[1]) != task.Waiting {
			t.Fatal("expected", task.Waiting, "got", lis[0].Gate.Get(key[1]))
		}

		if lis[0].Root != nil {
			t.Fatal("expected", nil, "got", lis[0].Root)
		}

		if lis[1].Meta.Get("test.api.io/key") != "bar" {
			t.Fatal("scheduling failed")
		}
		if lis[1].Gate != nil {
			t.Fatal("expected", nil, "got", lis[1].Gate)
		}
		if lis[1].Root.Len() != 1 {
			t.Fatal("expected", 1, "got", lis[1].Root.Len())
		}
		if lis[1].Root.Get(task.Object) != lis[0].Core.Map().Object() {
			t.Fatal("scheduled task must define root for task template")
		}
	}
}
