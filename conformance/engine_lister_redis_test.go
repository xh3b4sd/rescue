//go:build redis

package conformance

import (
	"reflect"
	"slices"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/xh3b4sd/logger"
	"github.com/xh3b4sd/redigo"
	"github.com/xh3b4sd/rescue"
	"github.com/xh3b4sd/rescue/engine"
	"github.com/xh3b4sd/rescue/task"
)

func Test_Engine_Lister(t *testing.T) {
	var err error

	var eon rescue.Interface
	{
		eon = engine.New(engine.Config{
			Expiry: 500 * time.Millisecond,
			Logger: logger.Fake(),
			Redigo: prgAll(redigo.Default()),
			Worker: "eon",
		})
	}

	var etw rescue.Interface
	{
		etw = engine.New(engine.Config{
			Expiry: 500 * time.Millisecond,
			Logger: logger.Fake(),
			Redigo: prgAll(redigo.Default()),
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
		lis, err = eon.Lister(engine.All())
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
		if !engine.IsTaskNotFound(err) {
			t.Fatal("expected", "taskNotFoundError", "got", err)
		}
	}
}

func Test_Engine_Lister_Gate(t *testing.T) {
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

	{
		tas := &task.Task{
			Gate: &task.Gate{
				"test.api.io/k-0": task.Trigger,
			},
			Meta: &task.Meta{
				"test.api.io/key": "foo",
			},
			Sync: &task.Sync{
				"test.api.io/zer": "0",
			},
		}

		err = eon.Create(tas)
		if err != nil {
			t.Fatal(err)
		}
	}

	{
		tas := &task.Task{
			Gate: &task.Gate{
				"test.api.io/k-1": task.Trigger,
			},
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
		if len(lis) != 3 {
			t.Fatal("expected", 3, "got", len(lis))
		}

		if lis[0].Meta.Get("test.api.io/key") != "foo" {
			t.Fatal("expected", "foo", "got", lis[0].Meta.Get("test.api.io/key"))
		}
		if lis[0].Gate.Len() != 1 {
			t.Fatal("expected", 1, "got", lis[0].Gate.Len())
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
		if lis[0].Gate.Get(key[0]) != task.Trigger {
			t.Fatal("expected", task.Trigger, "got", lis[0].Gate.Get(key[0]))
		}
		if lis[0].Root != nil {
			t.Fatal("expected", nil, "got", lis[0].Root)
		}
		if lis[0].Sync.Len() != 1 {
			t.Fatal("expected", 1, "got", lis[0].Sync.Len())
		}
		if lis[0].Sync.Get("test.api.io/zer") != "0" {
			t.Fatal("expected", "0", "got", lis[0].Sync.Get("test.api.io/zer"))
		}

		if lis[1].Meta.Get("test.api.io/key") != "foo" {
			t.Fatal("expected", "foo", "got", lis[1].Meta.Get("test.api.io/key"))
		}
		if lis[1].Gate.Len() != 1 {
			t.Fatal("expected", 1, "got", lis[1].Gate.Len())
		}

		{
			key = lis[1].Gate.Key()
		}

		{
			slices.Sort(key)
		}

		if key[0] != "test.api.io/k-1" {
			t.Fatal("expected", "test.api.io/k-1", "got", key[0])
		}
		if lis[1].Gate.Get(key[0]) != task.Trigger {
			t.Fatal("expected", task.Trigger, "got", lis[1].Gate.Get(key[0]))
		}
		if lis[1].Root != nil {
			t.Fatal("expected", nil, "got", lis[1].Root)
		}
		if lis[1].Sync != nil {
			t.Fatal("expected", nil, "got", lis[1].Sync)
		}

		if lis[2].Meta.Get("test.api.io/key") != "bar" {
			t.Fatal("expected", "foo", "got", lis[2].Meta.Get("test.api.io/key"))
		}
		if lis[2].Gate.Len() != 2 {
			t.Fatal("expected", 2, "got", lis[2].Gate.Len())
		}

		{
			key = lis[2].Gate.Key()
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
		if lis[2].Gate.Get(key[0]) != task.Waiting {
			t.Fatal("expected", task.Waiting, "got", lis[2].Gate.Get(key[0]))
		}
		if lis[2].Gate.Get(key[1]) != task.Waiting {
			t.Fatal("expected", task.Waiting, "got", lis[2].Gate.Get(key[1]))
		}

		if lis[2].Root != nil {
			t.Fatal("expected", nil, "got", lis[2].Root)
		}
		if lis[2].Sync.Len() != 2 {
			t.Fatal("expected", 2, "got", lis[2].Sync.Len())
		}
		if lis[2].Sync.Get("test.api.io/zer") != "n/a" {
			t.Fatal("expected", "n/a", "got", lis[2].Sync.Get("test.api.io/zer"))
		}
		if lis[2].Sync.Get("test.api.io/one") != "n/a" {
			t.Fatal("expected", "n/a", "got", lis[2].Sync.Get("test.api.io/one"))
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
		if len(lis) != 2 {
			t.Fatal("expected", 2, "got", len(lis))
		}

		if lis[0].Meta.Get("test.api.io/key") != "foo" {
			t.Fatal("expected", "foo", "got", lis[0].Meta.Get("test.api.io/key"))
		}
		if lis[0].Gate.Len() != 1 {
			t.Fatal("expected", 1, "got", lis[0].Gate.Len())
		}

		var key []string
		{
			key = lis[0].Gate.Key()
		}

		{
			slices.Sort(key)
		}

		if key[0] != "test.api.io/k-1" {
			t.Fatal("expected", "test.api.io/k-1", "got", key[0])
		}
		if lis[0].Gate.Get(key[0]) != task.Trigger {
			t.Fatal("expected", task.Trigger, "got", lis[0].Gate.Get(key[0]))
		}
		if lis[0].Root != nil {
			t.Fatal("expected", nil, "got", lis[0].Root)
		}
		if lis[0].Sync != nil {
			t.Fatal("expected", nil, "got", lis[0].Sync)
		}

		if lis[1].Meta.Get("test.api.io/key") != "bar" {
			t.Fatal("expected", "foo", "got", lis[1].Meta.Get("test.api.io/key"))
		}
		if lis[1].Gate.Len() != 2 {
			t.Fatal("expected", 2, "got", lis[1].Gate.Len())
		}

		{
			key = lis[1].Gate.Key()
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
		if lis[1].Gate.Get(key[0]) != task.Deleted {
			t.Fatal("expected", task.Deleted, "got", lis[1].Gate.Get(key[0]))
		}
		if lis[1].Gate.Get(key[1]) != task.Waiting {
			t.Fatal("expected", task.Waiting, "got", lis[1].Gate.Get(key[1]))
		}

		if lis[1].Root != nil {
			t.Fatal("expected", nil, "got", lis[1].Root)
		}
		if lis[1].Sync.Len() != 2 {
			t.Fatal("expected", 2, "got", lis[1].Sync.Len())
		}
		if lis[1].Sync.Get("test.api.io/zer") != "0" {
			t.Fatal("expected", "0", "got", lis[1].Sync.Get("test.api.io/zer"))
		}
		if lis[1].Sync.Get("test.api.io/one") != "n/a" {
			t.Fatal("expected", "n/a", "got", lis[1].Sync.Get("test.api.io/one"))
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
	}

	// This is the gating task template defining task.Gate.
	{
		var tas *task.Task
		{
			tas = lis[0]
		}

		var exp *task.Task
		{
			exp = &task.Task{
				Core: tas.Core,
				Gate: &task.Gate{
					"test.api.io/k-0": task.Waiting,
					"test.api.io/k-1": task.Waiting,
				},
				Meta: &task.Meta{
					"test.api.io/key": "bar",
				},
				Node: &task.Node{
					task.Method: task.MthdAny,
				},
				Sync: &task.Sync{
					"test.api.io/zer": "0",
					"test.api.io/one": "n/a",
				},
			}
		}

		{
			if !reflect.DeepEqual(tas, exp) {
				t.Fatalf("\n\n%s\n", cmp.Diff(exp, tas))
			}
		}
	}

	// This is the triggered task emitted by the gating task template above.
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
					task.Object: lis[0].Core.Map().Object(),
				},
				Sync: &task.Sync{
					"test.api.io/zer": "0",
					"test.api.io/one": "n/a",
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

func Test_Engine_Lister_Gate_Node_All(t *testing.T) {
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

	{
		tas := &task.Task{
			Gate: &task.Gate{
				"test.api.io/k-0": task.Trigger,
			},
			Meta: &task.Meta{
				"test.api.io/key": "foo",
			},
			Sync: &task.Sync{
				"test.api.io/zer": "0",
			},
		}

		err = eon.Create(tas)
		if err != nil {
			t.Fatal(err)
		}
	}

	{
		tas := &task.Task{
			Gate: &task.Gate{
				"test.api.io/k-1": task.Trigger,
			},
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
			Gate: &task.Gate{
				"test.api.io/k-0": task.Waiting,
				"test.api.io/k-1": task.Waiting,
			},
			Meta: &task.Meta{
				"test.api.io/key": "bar",
			},
			Node: &task.Node{
				task.Method: task.MthdAll,
			},
			Sync: &task.Sync{
				"test.api.io/zer": "n/a",
				"test.api.io/one": "n/a",
			},
		}

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
				Gate: &task.Gate{
					"test.api.io/k-0": task.Trigger,
				},
				Meta: &task.Meta{
					"test.api.io/key": "foo",
				},
				Node: &task.Node{
					task.Method: task.MthdAny,
				},
				Sync: &task.Sync{
					"test.api.io/zer": "0",
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
				Gate: &task.Gate{
					"test.api.io/k-1": task.Trigger,
				},
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
			tas = lis[2]
		}

		var exp *task.Task
		{
			exp = &task.Task{
				Core: tas.Core,
				Gate: &task.Gate{
					"test.api.io/k-0": task.Waiting,
					"test.api.io/k-1": task.Waiting,
				},
				Meta: &task.Meta{
					"test.api.io/key": "bar",
				},
				Node: &task.Node{
					task.Method: task.MthdAll,
				},
				Sync: &task.Sync{
					"test.api.io/zer": "n/a",
					"test.api.io/one": "n/a",
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
				Gate: &task.Gate{
					"test.api.io/k-0": task.Trigger,
				},
				Meta: &task.Meta{
					"test.api.io/key": "foo",
				},
				Node: &task.Node{
					task.Method: task.MthdAny,
				},
				Sync: &task.Sync{
					"test.api.io/zer": "0",
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
		lis, err = eon.Lister(engine.All())
		if err != nil {
			t.Fatal(err)
		}
	}

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
				Gate: &task.Gate{
					"test.api.io/k-1": task.Trigger,
				},
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

	// Ensure that any state in Task.Gate and Task.Sync is properly propagated,
	// since one of the associated trigger tasks got resolved.
	{
		var tas *task.Task
		{
			tas = lis[1]
		}

		var exp *task.Task
		{
			exp = &task.Task{
				Core: tas.Core,
				Gate: &task.Gate{
					"test.api.io/k-0": task.Deleted,
					"test.api.io/k-1": task.Waiting,
				},
				Meta: &task.Meta{
					"test.api.io/key": "bar",
				},
				Node: &task.Node{
					task.Method: task.MthdAll,
				},
				Sync: &task.Sync{
					"test.api.io/zer": "0",
					"test.api.io/one": "n/a",
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
				Gate: &task.Gate{
					"test.api.io/k-1": task.Trigger,
				},
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

	// After the second call to Engine.Delete the gating task template got
	// triggered to emit the new scheduled task.
	{
		var exp *task.Task
		{
			exp = &task.Task{
				Core: tas.Core,
				Meta: &task.Meta{
					"test.api.io/key": "bar",
				},
				Node: &task.Node{
					task.Method: task.MthdAll,
				},
				Root: &task.Root{
					task.Object: lis[1].Core.Map().Object(), // 2nd task in lis is the template
				},
				Sync: &task.Sync{
					"test.api.io/zer": "0",
					"test.api.io/one": "n/a",
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
		lis, err = eon.Lister(engine.All())
		if err != nil {
			t.Fatal(err)
		}
	}

	{
		if len(lis) != 2 {
			t.Fatal("expected", 2, "got", len(lis))
		}
	}

	// This is the gating task template defining task.Gate. Since both of its
	// trigger tasks got just deleted, its Task.Gate definition flipped to
	// "waiting" again for each trigger task label.
	{
		var tas *task.Task
		{
			tas = lis[0]
		}

		var exp *task.Task
		{
			exp = &task.Task{
				Core: tas.Core,
				Gate: &task.Gate{
					"test.api.io/k-0": task.Waiting,
					"test.api.io/k-1": task.Waiting,
				},
				Meta: &task.Meta{
					"test.api.io/key": "bar",
				},
				Node: &task.Node{
					task.Method: task.MthdAll,
				},
				Sync: &task.Sync{
					"test.api.io/zer": "0",
					"test.api.io/one": "n/a",
				},
			}
		}

		{
			if !reflect.DeepEqual(tas, exp) {
				t.Fatalf("\n\n%s\n", cmp.Diff(exp, tas))
			}
		}
	}

	// This is the triggered task emitted by the gating task template above.
	// Ensure it has the template's object ID in Task.Root, all state as defined
	// in Task.Sync and the task delivery method "all".
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
					task.Method: task.MthdAll,
				},
				Root: &task.Root{
					task.Object: lis[0].Core.Map().Object(), // 1st task in lis is the template
				},
				Sync: &task.Sync{
					"test.api.io/zer": "0",
					"test.api.io/one": "n/a",
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
