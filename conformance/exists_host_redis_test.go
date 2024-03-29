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

func Test_Engine_Exists_Host(t *testing.T) {
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
			Logger: logger.Fake(),
			Redigo: red,
		})
	}

	{
		exi, err := eon.Exists(&task.Task{Node: &task.Node{task.Method: task.MthdAll}})
		if err != nil {
			t.Fatal(err)
		}

		if exi {
			t.Fatal("expected", false, "got", true)
		}
	}

	{
		exi, err := eon.Exists(&task.Task{Node: &task.Node{task.Method: task.MthdUni}})
		if err != nil {
			t.Fatal(err)
		}

		if exi {
			t.Fatal("expected", false, "got", true)
		}
	}

	{
		exi, err := eon.Exists(&task.Task{Node: &task.Node{task.Method: task.MthdUni, task.Worker: "etw"}})
		if err != nil {
			t.Fatal(err)
		}

		if exi {
			t.Fatal("expected", false, "got", true)
		}
	}

	{
		exi, err := eon.Exists(&task.Task{Node: &task.Node{task.Method: task.MthdAny}})
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
			Node: &task.Node{
				task.Method: task.MthdAll,
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

	{
		tas := &task.Task{
			Meta: &task.Meta{
				"test.api.io/key": "baz",
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

	{
		exi, err := eon.Exists(&task.Task{Node: &task.Node{task.Method: task.MthdAll}})
		if err != nil {
			t.Fatal(err)
		}

		if !exi {
			t.Fatal("expected", true, "got", false)
		}
	}

	{
		exi, err := eon.Exists(&task.Task{Node: &task.Node{task.Method: task.MthdUni}})
		if err != nil {
			t.Fatal(err)
		}

		if !exi {
			t.Fatal("expected", true, "got", false)
		}
	}

	{
		exi, err := eon.Exists(&task.Task{Node: &task.Node{task.Method: task.MthdUni, task.Worker: "etw"}})
		if err != nil {
			t.Fatal(err)
		}

		if !exi {
			t.Fatal("expected", true, "got", false)
		}
	}

	// There is no "any" task created, so it must not exist.
	{
		exi, err := eon.Exists(&task.Task{Node: &task.Node{task.Method: task.MthdAny}})
		if err != nil {
			t.Fatal(err)
		}

		if exi {
			t.Fatal("expected", false, "got", true)
		}
	}
}
