//go:build redis

package engine

import (
	"fmt"
	"strconv"
	"sync"
	"testing"
	"time"

	"github.com/xh3b4sd/budget/v3/pkg/breaker"
	"github.com/xh3b4sd/logger"
	"github.com/xh3b4sd/redigo"
	"github.com/xh3b4sd/redigo/pkg/client"
	"github.com/xh3b4sd/rescue/task"
	"github.com/xh3b4sd/tracer"
)

func Test_Engine_Balance(t *testing.T) {
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

	{
		for i := 0; i < 10; i++ {
			tas := &task.Task{
				Meta: &task.Meta{
					"test.api.io/num": strconv.Itoa(i),
				},
			}

			err = eon.Create(tas)
			if err != nil {
				t.Fatal(err)
			}
		}
	}

	var lon []*task.Task
	var ltw []*task.Task
	var lth []*task.Task

	// Make the system aware of each worker.
	{
		{
			tas, err := eon.Search()
			if err != nil {
				t.Fatal(err)
			}

			lon = append(lon, tas)
		}

		{
			tas, err := etw.Search()
			if err != nil {
				t.Fatal(err)
			}

			ltw = append(ltw, tas)
		}

		{
			tas, err := eth.Search()
			if err != nil {
				t.Fatal(err)
			}

			lth = append(lth, tas)
		}
	}

	// Let every worker consume as many tasks as they can.
	{
		for {
			tas, err := eon.Search()
			if IsTaskNotFound(err) {
				break
			} else if err != nil {
				t.Fatal(err)
			}

			lon = append(lon, tas)
		}

		for {
			tas, err := etw.Search()
			if IsTaskNotFound(err) {
				break
			} else if err != nil {
				t.Fatal(err)
			}

			ltw = append(ltw, tas)
		}

		for {
			tas, err := eth.Search()
			if IsTaskNotFound(err) {
				break
			} else if err != nil {
				t.Fatal(err)
			}

			lth = append(lth, tas)
		}
	}

	{
		if len(lon) != 4 {
			t.Fatal("worker one must claim 4 tasks")
		}
		if len(ltw) != 3 {
			t.Fatal("worker two must claim 3 tasks")
		}
		if len(lth) != 3 {
			t.Fatal("worker three must claim 3 tasks")
		}
	}
}

func Test_Engine_Create(t *testing.T) {
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

	{
		if len(lis) != 3 {
			t.Fatal("expected 3 tasks listed")
		}
		if lis[0].Meta.Get("test.api.io/key") != "foo" {
			t.Fatal("scheduling failed")
		}
		if lis[0].Root != nil {
			t.Fatal("scheduling failed")
		}
		if lis[1].Meta.Get("test.api.io/key") != "zap" {
			t.Fatal("scheduling failed")
		}
		if lis[1].Root.Get("test.api.io/key") != "foo" {
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
			t.Fatal("expected 1 tasks listed")
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
			t.Fatal("queue must be empty")
		}
	}
}

func Test_Engine_Delete(t *testing.T) {
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
		exi, err := eon.Exists(&task.Task{Meta: &task.Meta{"test.api.io/key": "foo"}, Root: &task.Root{"test.api.io/key": "rrr"}})
		if err != nil {
			t.Fatal(err)
		}

		if exi {
			t.Fatal("task must not exist")
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
			t.Fatal("task must be deleted by owner")
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

	{
		tas, err = eon.Search()
		if err != nil {
			t.Fatal(err)
		}

		if tas.Meta.Get("test.api.io/key") != "foo" {
			t.Fatal("scheduling failed")
		}
		if tas.Root.Get("test.api.io/key") != "rrr" {
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

		if tas.Meta.Get("test.api.io/key") != "baz" {
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

	{
		_, err = eon.Search()
		if !IsTaskNotFound(err) {
			t.Fatal("queue must be empty")
		}
	}
}

func Test_Engine_Expire(t *testing.T) {
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

func Test_Engine_Extend(t *testing.T) {
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
			Expiry: time.Second,
			Logger: logger.Fake(),
			Redigo: red,
			Worker: "eon",
		})
	}

	var etw *Engine
	{
		etw = New(Config{
			Expiry: time.Second,
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
		tas, err = etw.Search()
		if err != nil {
			t.Fatal(err)
		}

		if tas.Meta.Get("test.api.io/key") != "foo" {
			t.Fatal("scheduling failed")
		}
	}

	{
		_, err = eon.Search()
		if !IsTaskNotFound(err) {
			t.Fatal("queue must be empty")
		}
	}

	{
		time.Sleep(200 * time.Millisecond)
	}

	{
		err = eon.Expire()
		if err != nil {
			t.Fatal(err)
		}
		err = etw.Expire()
		if err != nil {
			t.Fatal(err)
		}
	}

	{
		_, err = eon.Search()
		if !IsTaskNotFound(err) {
			t.Fatal("queue must be empty")
		}

		err = etw.Extend(tas)
		if err != nil {
			t.Fatal(err)
		}
	}

	{
		time.Sleep(200 * time.Millisecond)
	}

	{
		err = eon.Expire()
		if err != nil {
			t.Fatal(err)
		}
		err = etw.Expire()
		if err != nil {
			t.Fatal(err)
		}
	}

	{
		_, err = eon.Search()
		if !IsTaskNotFound(err) {
			t.Fatal("queue must be empty")
		}

		err = etw.Extend(tas)
		if err != nil {
			t.Fatal(err)
		}
	}

	{
		time.Sleep(200 * time.Millisecond)
	}

	{
		err = eon.Expire()
		if err != nil {
			t.Fatal(err)
		}
		err = etw.Expire()
		if err != nil {
			t.Fatal(err)
		}
	}

	{
		_, err = eon.Search()
		if !IsTaskNotFound(err) {
			t.Fatal("queue must be empty")
		}

		err = etw.Extend(tas)
		if err != nil {
			t.Fatal(err)
		}
	}

	{
		time.Sleep(200 * time.Millisecond)
	}

	{
		err = eon.Expire()
		if err != nil {
			t.Fatal(err)
		}
		err = etw.Expire()
		if err != nil {
			t.Fatal(err)
		}
	}

	{
		_, err = eon.Search()
		if !IsTaskNotFound(err) {
			t.Fatal("queue must be empty")
		}

		err = etw.Extend(tas)
		if err != nil {
			t.Fatal(err)
		}
	}

	{
		time.Sleep(200 * time.Millisecond)
	}

	{
		err = eon.Expire()
		if err != nil {
			t.Fatal(err)
		}
		err = etw.Expire()
		if err != nil {
			t.Fatal(err)
		}
	}

	{
		_, err = eon.Search()
		if !IsTaskNotFound(err) {
			t.Fatal("queue must be empty")
		}

		err = etw.Extend(tas)
		if err != nil {
			t.Fatal(err)
		}
	}

	{
		time.Sleep(1 * time.Second)
	}

	{
		err = eon.Expire()
		if err != nil {
			t.Fatal(err)
		}
		err = etw.Expire()
		if err != nil {
			t.Fatal(err)
		}
	}

	{
		err = etw.Extend(tas)
		if !IsTaskOutdated(err) {
			t.Fatal("task must be extended by owner")
		}
	}

	{
		err = etw.Delete(tas)
		if !IsTaskOutdated(err) {
			t.Fatal("task must be deleted by owner")
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
		_, err = etw.Search()
		if !IsTaskNotFound(err) {
			t.Fatal("queue must be empty")
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

		_, err = etw.Search()
		if !IsTaskNotFound(err) {
			t.Fatal("queue must be empty")
		}
	}
}

func Test_Engine_Lifecycle_Race(t *testing.T) {
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

	{
		tas := &task.Task{
			Meta: &task.Meta{
				"test.api.io/key": "foo",
			},
		}

		exi, err := eon.Exists(tas)
		if err != nil {
			t.Fatal(err)
		}

		if !exi {
			t.Fatal("task must exist")
		}
	}

	{
		tas := &task.Task{
			Meta: &task.Meta{
				"test.api.io/key": "bar",
			},
		}

		exi, err := etw.Exists(tas)
		if err != nil {
			t.Fatal(err)
		}

		if !exi {
			t.Fatal("task must exist")
		}
	}

	erc := make(chan error, 1)

	go func() {
		defer close(erc)

		var s []string
		var w sync.WaitGroup

		w.Add(2)

		go func() {
			defer w.Done()

			tas, err := eon.Search()
			if err != nil {
				erc <- tracer.Mask(err)
				return
			}

			s = append(s, tas.Meta.Get("test.api.io/key"))

			err = eon.Delete(tas)
			if err != nil {
				erc <- tracer.Mask(err)
				return
			}
		}()

		go func() {
			defer w.Done()

			tas, err := eon.Search()
			if err != nil {
				erc <- tracer.Mask(err)
				return
			}

			s = append(s, tas.Meta.Get("test.api.io/key"))

			err = etw.Delete(tas)
			if err != nil {
				erc <- tracer.Mask(err)
				return
			}
		}()

		w.Wait()

		{
			if len(s) != 2 {
				erc <- fmt.Errorf("length of s must be 2")
				return
			}
			if s[0] == "foo" && s[1] != "bar" {
				erc <- fmt.Errorf("scheduling failed")
				return
			}
			if s[1] == "foo" && s[0] != "bar" {
				erc <- fmt.Errorf("scheduling failed")
				return
			}
		}

		{
			tas := &task.Task{
				Meta: &task.Meta{
					"test.api.io/key": "foo",
				},
			}

			exi, err := eon.Exists(tas)
			if err != nil {
				panic(err)
			}

			if exi {
				panic("task must not exist")
			}
		}

		{
			tas := &task.Task{
				Meta: &task.Meta{
					"test.api.io/key": "bar",
				},
			}

			exi, err := etw.Exists(tas)
			if err != nil {
				panic(err)
			}

			if exi {
				panic("task must not exist")
			}
		}

		{
			_, err = eon.Search()
			if !IsTaskNotFound(err) {
				erc <- fmt.Errorf("queue must be empty")
				return
			}

			_, err = etw.Search()
			if !IsTaskNotFound(err) {
				erc <- fmt.Errorf("queue must be empty")
				return
			}
		}
	}()

	{
		err = <-erc
		if err != nil {
			t.Fatal(err)
		}
	}
}

func Test_Engine_Lifecycle_Sync(t *testing.T) {
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
			Queue:  "one", // engines use different queues
		})
	}

	var etw *Engine
	{
		etw = New(Config{
			Expiry: 500 * time.Millisecond,
			Logger: logger.Fake(),
			Redigo: red,
			Queue:  "two", // engines use different queues
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
				"test.api.io/key": "zap",
			},
			Root: &task.Root{
				"test.api.io/key": "oth",
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
			Root: &task.Root{
				"test.api.io/key": "roo",
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

	{
		tas := &task.Task{
			Root: &task.Root{
				"test.api.io/key": "rrr",
			},
		}

		err = etw.Create(tas)
		if !IsTaskMetaEmpty(err) {
			t.Fatal("expected task creation to fail without Task.Meta")
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
		if len(lis) != 2 {
			t.Fatal("expected 2 tasks listed")
		}
	}

	{
		tas := &task.Task{
			Root: &task.Root{
				"test.api.io/key": "roo",
			},
		}

		lis, err = eon.Lister(tas)
		if err != nil {
			t.Fatal(err)
		}
	}

	{
		if len(lis) != 1 {
			t.Fatal("expected 1 task listed")
		}
	}

	{
		tas := &task.Task{
			Meta: &task.Meta{
				"test.api.io/key": "foo",
			},
		}

		lis, err = etw.Lister(tas)
		if err != nil {
			t.Fatal(err)
		}
	}

	{
		if len(lis) != 1 {
			t.Fatal("expected 1 task listed")
		}
	}

	{
		tas := &task.Task{
			Root: &task.Root{
				"test.api.io/key": "roo",
			},
		}

		lis, err = etw.Lister(tas)
		if err != nil {
			t.Fatal(err)
		}
	}

	{
		if len(lis) != 0 {
			t.Fatal("expected 0 tasks listed")
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
			t.Fatal("scheduling failed")
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

		if tas.Meta.Get("test.api.io/key") != "zap" {
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

	{
		_, err = etw.Search()
		if !IsTaskNotFound(err) {
			t.Fatal("queue must be empty")
		}
	}
}

func Test_Engine_Lister_Order(t *testing.T) {
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
			t.Fatal("expected 2 task labels")
		}
		if len(*lis[1].Meta.All("test*")) != 3 {
			t.Fatal("expected 3 task labels")
		}
		if len(*lis[2].Meta.All("test*")) != 1 {
			t.Fatal("expected 1 task labels")
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
			t.Fatal("expected 2 task labels")
		}
		if len(*lis[1].Meta.All("test*")) != 3 {
			t.Fatal("expected 3 task labels")
		}
		if len(*lis[2].Meta.All("test*")) != 1 {
			t.Fatal("expected 1 task labels")
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
			t.Fatal("expected 2 task labels")
		}
		if len(*lis[1].Meta.All("test*")) != 3 {
			t.Fatal("expected 3 task labels")
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
			t.Fatal("expected 2 task labels")
		}
		if len(*lis[1].Meta.All("test*")) != 3 {
			t.Fatal("expected 3 task labels")
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
			t.Fatal("expected 2 task labels")
		}
		if len(*lis[1].Meta.All("test*")) != 3 {
			t.Fatal("expected 3 task labels")
		}
		if len(*lis[2].Meta.All("test*")) != 1 {
			t.Fatal("expected 1 task labels")
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
			t.Fatal("scheduling failed")
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
