//go:build redis

package engine

import (
	"testing"
	"time"

	"github.com/xh3b4sd/logger"
	"github.com/xh3b4sd/redigo"
	"github.com/xh3b4sd/rescue/task"
)

func Test_Engine_Expire(t *testing.T) {
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
			Expiry: time.Millisecond,
			Logger: logger.Fake(),
			Redigo: red,
			Worker: "eon",
		})
	}

	var etw *Engine
	{
		etw = New(Config{
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

	// For engine one we simulate failure so that the acquired task can
	// expire and be rescheduled to engine two. For the simulation we
	// call Expire which is the responsibility of every worker to do
	// periodically. It does not matter which engine executes the
	// expiration process.
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
		if !IsTaskNotFound(err) {
			t.Fatal("queue must be empty")
		}

		_, err = etw.Search()
		if !IsTaskNotFound(err) {
			t.Fatal("queue must be empty")
		}
	}
}
