//go:build redis

package engine

import (
	"testing"
	"time"

	"github.com/xh3b4sd/logger"
	"github.com/xh3b4sd/redigo"
	"github.com/xh3b4sd/rescue/task"
)

func Test_Engine_Delete(t *testing.T) {
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
			Expiry: 1 * time.Millisecond,
			Logger: logger.Fake(),
			Redigo: red,
			Worker: "eon",
		})
	}

	var etw *Engine
	{
		etw = New(Config{
			Expiry: 1 * time.Millisecond,
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

	var tas *task.Task
	{
		tas, err = eon.Search()
		if err != nil {
			t.Fatal(err)
		}

		if tas.Meta.Get("test.api.io/key") != "foo" {
			t.Fatal("scheduling failed")
		}
	}

	{
		time.Sleep(1 * time.Millisecond)
	}

	{
		err = eon.Expire()
		if err != nil {
			t.Fatal(err)
		}
	}

	{
		exi, err := eon.Exists(&task.Task{Core: tas.Core.All(task.Object, task.Worker)})
		if err != nil {
			t.Fatal(err)
		}

		if exi {
			t.Fatal("task must not be owned by eon")
		}
	}

	{
		err = etw.Delete(tas)
		if !IsTaskOutdated(err) {
			t.Fatal("task must be deleted by owner")
		}
	}

	{
		tas.Core.Set().Bypass(true)
	}

	{
		err = etw.Delete(tas)
		if err != nil {
			t.Fatal(err)
		}
	}

	{
		_, err = eon.Search()
		if !IsTaskNotFound(err) {
			t.Fatal("queue must be empty")
		}
	}
}

func Test_Engine_Delete_Gate(t *testing.T) {
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

	var lis []*task.Task
	{
		lis, err = eon.Lister(All())
		if err != nil {
			t.Fatal(err)
		}
	}

	{
		if len(lis) != 0 {
			t.Fatal("expected", 0, "got", len(lis))
		}
	}

	var tas *task.Task
	{
		tas = &task.Task{
			Meta: &task.Meta{
				"test.api.io/key": "bar",
			},
			Gate: &task.Gate{
				"test.api.io/k-0": task.Waiting,
				"test.api.io/k-1": task.Waiting,
			},
			Sync: &task.Sync{
				"test.api.io/zer": "n/a",
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
		err = eon.Delete(tas)
		if !IsTaskOutdated(err) {
			t.Fatal("task must be deleted by owner")
		}
	}

	// Templates for triggered tasks can only be deleted when bypassing the
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
		lis, err = eon.Lister(All())
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
