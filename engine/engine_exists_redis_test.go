//go:build redis

package engine

import (
	"testing"
	"time"

	"github.com/xh3b4sd/budget/v3/pkg/breaker"
	"github.com/xh3b4sd/logger"
	"github.com/xh3b4sd/redigo"
	"github.com/xh3b4sd/redigo/pkg/client"
	"github.com/xh3b4sd/rescue/task"
)

func Test_Engine_Exists(t *testing.T) {
	var err error

	var red redigo.Interface
	{
		c := client.Config{
			Kind: client.KindSingle,
			Locker: client.ConfigLocker{
				Budget: breaker.Default(),
			},
		}

		red, err = client.New(c)
		if err != nil {
			t.Fatal(err)
		}

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
