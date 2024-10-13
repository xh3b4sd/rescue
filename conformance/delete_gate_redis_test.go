//go:build redis

package conformance

import (
	"testing"

	"github.com/xh3b4sd/logger"
	"github.com/xh3b4sd/redigo"
	"github.com/xh3b4sd/rescue"
	"github.com/xh3b4sd/rescue/engine"
	"github.com/xh3b4sd/rescue/task"
)

func Test_Engine_Delete_Gate(t *testing.T) {
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

	var tas *task.Task
	{
		tas = &task.Task{
			Gate: &task.Gate{
				"test.api.io/k-0": task.Waiting,
				"test.api.io/k-1": task.Waiting,
			},
			Meta: &task.Meta{
				"test.api.io/key": "bar",
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

	// Ensure that deleting task templates does not work without bypassing.
	{
		err = eon.Delete(tas)
		if !engine.IsTaskOutdated(err) {
			t.Fatal("expected", "taskOutdatedError", "got", err)
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
