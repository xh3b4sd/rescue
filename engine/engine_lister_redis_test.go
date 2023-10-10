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

func Test_Engine_Lister(t *testing.T) {
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

	var etw *Engine
	{
		etw = New(Config{
			Expiry: 500 * time.Millisecond,
			Logger: logger.Fake(),
			Redigo: red,
			Worker: "etw",
		})
	}

	{
		tas := &task.Task{
			Meta: &task.Meta{
				"test.api.io/key": "foo",
				"test.api.io/zer": "tru",
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
				"test.api.io/zer": "tru",
				"test.api.io/sin": "baz",
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
		}

		err = etw.Create(tas)
		if err != nil {
			t.Fatal(err)
		}
	}

	var lis []*task.Task
	{
		tas := &task.Task{
			Meta: &task.Meta{
				"test.api.io/key": "foo",
			},
		}

		lis, err = eon.Lister(tas)
		if err != nil {
			t.Fatal(err)
		}
	}

	{
		if len(lis) != 3 {
			t.Fatal("expected 3 tasks listed")
		}
		if len(*lis[0].Meta.All("test*")) != 2 {
			t.Fatal("expected", 2, "got", len(*lis[0].Meta.All("test*")))
		}
		if len(*lis[1].Meta.All("test*")) != 3 {
			t.Fatal("expected", 3, "got", len(*lis[1].Meta.All("test*")))
		}
		if len(*lis[2].Meta.All("test*")) != 1 {
			t.Fatal("expected", 1, "got", len(*lis[2].Meta.All("test*")))
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
		tas := &task.Task{
			Meta: &task.Meta{
				"test.api.io/key": "foo",
			},
		}

		lis, err = eon.Lister(tas)
		if err != nil {
			t.Fatal(err)
		}
	}

	{
		if len(lis) != 3 {
			t.Fatal("expected 3 tasks listed")
		}
		if len(*lis[0].Meta.All("test*")) != 2 {
			t.Fatal("expected", 2, "got", len(*lis[0].Meta.All("test*")))
		}
		if len(*lis[1].Meta.All("test*")) != 3 {
			t.Fatal("expected", 3, "got", len(*lis[1].Meta.All("test*")))
		}
		if len(*lis[2].Meta.All("test*")) != 1 {
			t.Fatal("expected", 1, "got", len(*lis[2].Meta.All("test*")))
		}
	}

	{
		tas := &task.Task{
			Meta: &task.Meta{
				"test.api.io/zer": "tru",
			},
		}

		lis, err = eon.Lister(tas)
		if err != nil {
			t.Fatal(err)
		}
	}

	{
		if len(lis) != 2 {
			t.Fatal("expected 2 tasks listed")
		}
		if len(*lis[0].Meta.All("test*")) != 2 {
			t.Fatal("expected", 2, "got", len(*lis[0].Meta.All("test*")))
		}
		if len(*lis[1].Meta.All("test*")) != 3 {
			t.Fatal("expected", 3, "got", len(*lis[1].Meta.All("test*")))
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
		tas := &task.Task{
			Meta: &task.Meta{
				"test.api.io/zer": "tru",
			},
		}

		lis, err = eon.Lister(tas)
		if err != nil {
			t.Fatal(err)
		}
	}

	{
		if len(lis) != 2 {
			t.Fatal("expected 2 tasks listed")
		}
		if len(*lis[0].Meta.All("test*")) != 2 {
			t.Fatal("expected", 2, "got", len(*lis[0].Meta.All("test*")))
		}
		if len(*lis[1].Meta.All("test*")) != 3 {
			t.Fatal("expected", 3, "got", len(*lis[1].Meta.All("test*")))
		}
	}

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
		if len(*lis[0].Meta.All("test*")) != 2 {
			t.Fatal("expected", 2, "got", len(*lis[0].Meta.All("test*")))
		}
		if len(*lis[1].Meta.All("test*")) != 3 {
			t.Fatal("expected", 3, "got", len(*lis[1].Meta.All("test*")))
		}
		if len(*lis[2].Meta.All("test*")) != 1 {
			t.Fatal("expected", 1, "got", len(*lis[2].Meta.All("test*")))
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
		if tas.Meta.Get("test.api.io/zer") != "tru" {
			t.Fatal("expected", "tru", "got", tas.Meta.Get("test.api.io/zer"))
		}
	}

	{
		err = eon.Delete(tas)
		if err != nil {
			t.Fatal(err)
		}
	}

	{
		tas, err = etw.Search()
		if err != nil {
			t.Fatal(err)
		}

		if tas.Meta.Get("test.api.io/key") != "foo" {
			t.Fatal("expected", "foo", "got", tas.Meta.Get("test.api.io/key"))
		}
		if tas.Meta.Get("test.api.io/zer") != "tru" {
			t.Fatal("expected", "tru", "got", tas.Meta.Get("test.api.io/zer"))
		}
		if tas.Meta.Get("test.api.io/sin") != "baz" {
			t.Fatal("expected", "baz", "got", tas.Meta.Get("test.api.io/sin"))
		}
	}

	{
		err = etw.Delete(tas)
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
		_, err = eon.Search()
		if !IsTaskNotFound(err) {
			t.Fatal("queue must be empty")
		}
	}
}
