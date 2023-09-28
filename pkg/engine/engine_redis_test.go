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
	"github.com/xh3b4sd/tracer"

	"github.com/xh3b4sd/rescue/pkg/metadata"
	"github.com/xh3b4sd/rescue/pkg/task"
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

			Owner: "eon",
		})
	}

	var etw *Engine
	{
		etw = New(Config{
			Logger: logger.Fake(),
			Redigo: red,

			Owner: "etw",
		})
	}

	var eth *Engine
	{
		eth = New(Config{
			Logger: logger.Fake(),
			Redigo: red,

			Owner: "eth",
		})
	}

	{
		for i := 0; i < 10; i++ {
			tas := &task.Task{
				Obj: task.TaskObj{
					Metadata: map[string]string{
						"test.rescue.io/num": strconv.Itoa(i),
					},
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
			Logger: logger.Fake(),
			Redigo: red,

			Owner: "eon",
			TTL:   1 * time.Millisecond,
		})
	}

	var etw *Engine
	{
		etw = New(Config{
			Logger: logger.Fake(),
			Redigo: red,

			Owner: "etw",
			TTL:   1 * time.Millisecond,
		})
	}

	{
		tas := &task.Task{
			Obj: task.TaskObj{
				Metadata: map[string]string{
					"test.rescue.io/key": "foo",
				},
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

		if tas.Obj.Metadata["test.rescue.io/key"] != "foo" {
			t.Fatal("scheduling failed")
		}
	}

	{
		err = eon.Expire()
		if err != nil {
			t.Fatal(err)
		}
	}

	{
		exi, err := eon.Exists(tas.With(metadata.ID, metadata.Owner))
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
		tas.SetPrivileged(true)
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
			Logger: logger.Fake(),
			Redigo: red,

			Owner: "eon",
			TTL:   500 * time.Millisecond,
		})
	}

	{
		tas := &task.Task{
			Obj: task.TaskObj{
				Metadata: map[string]string{
					"test.rescue.io/key": "foo",
				},
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

		if tas.Obj.Metadata["test.rescue.io/key"] != "foo" {
			t.Fatal("scheduling failed")
		}
	}

	{
		exi, err := eon.Exists(tas.With(metadata.ID))
		if err != nil {
			t.Fatal(err)
		}

		if !exi {
			t.Fatal("task must exist")
		}
	}

	{
		exi, err := eon.Exists(tas.With(metadata.ID, metadata.Owner))
		if err != nil {
			t.Fatal(err)
		}

		if !exi {
			t.Fatal("task must exist")
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
		exi, err := eon.Exists(tas.With(metadata.ID))
		if err != nil {
			t.Fatal(err)
		}

		if !exi {
			t.Fatal("task must exist")
		}
	}

	{
		exi, err := eon.Exists(tas.With(metadata.ID, metadata.Owner))
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

		if tas.Obj.Metadata["test.rescue.io/key"] != "foo" {
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
			Logger: logger.Fake(),
			Redigo: red,

			Owner: "eon",
			TTL:   time.Millisecond,
		})
	}

	var etw *Engine
	{
		etw = New(Config{
			Logger: logger.Fake(),
			Redigo: red,

			Owner: "etw",
			TTL:   time.Millisecond,
		})
	}

	{
		tas := &task.Task{
			Obj: task.TaskObj{
				Metadata: map[string]string{
					"test.rescue.io/key": "foo",
				},
			},
		}

		err = eon.Create(tas)
		if err != nil {
			t.Fatal(err)
		}
	}

	{
		tas := &task.Task{
			Obj: task.TaskObj{
				Metadata: map[string]string{
					"test.rescue.io/key": "bar",
				},
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

		s = append(s, tas.Obj.Metadata["test.rescue.io/key"])

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

		s = append(s, tas.Obj.Metadata["test.rescue.io/key"])

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
			Logger: logger.Fake(),
			Redigo: red,

			Owner: "eon",
			TTL:   time.Second,
		})
	}

	var etw *Engine
	{
		etw = New(Config{
			Logger: logger.Fake(),
			Redigo: red,

			Owner: "etw",
			TTL:   time.Second,
		})
	}

	{
		tas := &task.Task{
			Obj: task.TaskObj{
				Metadata: map[string]string{
					"test.rescue.io/key": "foo",
				},
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

		if tas.Obj.Metadata["test.rescue.io/key"] != "foo" {
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

		if tas.Obj.Metadata["test.rescue.io/key"] != "foo" {
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

func Test_Engine_Lifecycle(t *testing.T) {
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

			Owner: "eon",
		})
	}

	var etw *Engine
	{
		etw = New(Config{
			Logger: logger.Fake(),
			Redigo: red,

			Owner: "etw",
		})
	}

	{
		tas := &task.Task{
			Obj: task.TaskObj{
				Metadata: map[string]string{
					"test.rescue.io/key": "foo",
				},
			},
		}

		err = eon.Create(tas)
		if err != nil {
			t.Fatal(err)
		}
	}

	{
		tas := &task.Task{
			Obj: task.TaskObj{
				Metadata: map[string]string{
					"test.rescue.io/key": "bar",
				},
			},
		}

		err = etw.Create(tas)
		if err != nil {
			t.Fatal(err)
		}
	}

	{
		tas := &task.Task{
			Obj: task.TaskObj{
				Metadata: map[string]string{
					"test.rescue.io/key": "foo",
				},
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
			Obj: task.TaskObj{
				Metadata: map[string]string{
					"test.rescue.io/key": "bar",
				},
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

			s = append(s, tas.Obj.Metadata["test.rescue.io/key"])

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

			s = append(s, tas.Obj.Metadata["test.rescue.io/key"])

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
				Obj: task.TaskObj{
					Metadata: map[string]string{
						"test.rescue.io/key": "foo",
					},
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
				Obj: task.TaskObj{
					Metadata: map[string]string{
						"test.rescue.io/key": "bar",
					},
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
			Logger: logger.Fake(),
			Redigo: red,

			Owner: "eon",
			TTL:   500 * time.Millisecond,
		})
	}

	var etw *Engine
	{
		etw = New(Config{
			Logger: logger.Fake(),
			Redigo: red,

			Owner: "etw",
			TTL:   500 * time.Millisecond,
		})
	}

	{
		tas := &task.Task{
			Obj: task.TaskObj{
				Metadata: map[string]string{
					"test.rescue.io/key": "foo",
					"test.rescue.io/zer": "tru",
				},
			},
		}

		err = eon.Create(tas)
		if err != nil {
			t.Fatal(err)
		}
	}

	{
		tas := &task.Task{
			Obj: task.TaskObj{
				Metadata: map[string]string{
					"test.rescue.io/key": "foo",
					"test.rescue.io/zer": "tru",
					"test.rescue.io/sin": "baz",
				},
			},
		}

		err = eon.Create(tas)
		if err != nil {
			t.Fatal(err)
		}
	}

	{
		tas := &task.Task{
			Obj: task.TaskObj{
				Metadata: map[string]string{
					"test.rescue.io/key": "foo",
				},
			},
		}

		err = etw.Create(tas)
		if err != nil {
			t.Fatal(err)
		}
	}

	var lis task.Tasks
	{
		tas := &task.Task{
			Obj: task.TaskObj{
				Metadata: map[string]string{
					"test.rescue.io/key": "foo",
				},
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
		if len(lis[0].Pref("test").Obj.Metadata) != 2 {
			t.Fatal("expected 2 task labels")
		}
		if len(lis[1].Pref("test").Obj.Metadata) != 3 {
			t.Fatal("expected 3 task labels")
		}
		if len(lis[2].Pref("test").Obj.Metadata) != 1 {
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
			Obj: task.TaskObj{
				Metadata: map[string]string{
					"test.rescue.io/key": "foo",
				},
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
		if len(lis[0].Pref("test").Obj.Metadata) != 2 {
			t.Fatal("expected 2 task labels")
		}
		if len(lis[1].Pref("test").Obj.Metadata) != 3 {
			t.Fatal("expected 3 task labels")
		}
		if len(lis[2].Pref("test").Obj.Metadata) != 1 {
			t.Fatal("expected 1 task labels")
		}
	}

	{
		tas := &task.Task{
			Obj: task.TaskObj{
				Metadata: map[string]string{
					"test.rescue.io/zer": "tru",
				},
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
		if len(lis[0].Pref("test").Obj.Metadata) != 2 {
			t.Fatal("expected 2 task labels")
		}
		if len(lis[1].Pref("test").Obj.Metadata) != 3 {
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
			Obj: task.TaskObj{
				Metadata: map[string]string{
					"test.rescue.io/zer": "tru",
				},
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
		if len(lis[0].Pref("test").Obj.Metadata) != 2 {
			t.Fatal("expected 2 task labels")
		}
		if len(lis[1].Pref("test").Obj.Metadata) != 3 {
			t.Fatal("expected 3 task labels")
		}
	}

	{
		lis, err = eon.Lister(task.All())
		if err != nil {
			t.Fatal(err)
		}
	}

	{
		if len(lis) != 3 {
			t.Fatal("expected 3 tasks listed")
		}
		if len(lis[0].Pref("test").Obj.Metadata) != 2 {
			t.Fatal("expected 2 task labels")
		}
		if len(lis[1].Pref("test").Obj.Metadata) != 3 {
			t.Fatal("expected 3 task labels")
		}
		if len(lis[2].Pref("test").Obj.Metadata) != 1 {
			t.Fatal("expected 1 task labels")
		}
	}

	var tas *task.Task
	{
		tas, err = eon.Search()
		if err != nil {
			t.Fatal(err)
		}

		if tas.Obj.Metadata["test.rescue.io/key"] != "foo" {
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

		if tas.Obj.Metadata["test.rescue.io/key"] != "foo" {
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

		if tas.Obj.Metadata["test.rescue.io/key"] != "foo" {
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

func Test_Engine_Queue(t *testing.T) {
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

			Queue: "one",
			TTL:   500 * time.Millisecond,
		})
	}

	var etw *Engine
	{
		etw = New(Config{
			Logger: logger.Fake(),
			Redigo: red,

			Queue: "two",
			TTL:   500 * time.Millisecond,
		})
	}

	{
		tas := &task.Task{
			Obj: task.TaskObj{
				Metadata: map[string]string{
					"test.rescue.io/key": "foo",
					"test.rescue.io/zer": "tru",
				},
			},
		}

		err = eon.Create(tas)
		if err != nil {
			t.Fatal(err)
		}
	}

	{
		tas := &task.Task{
			Obj: task.TaskObj{
				Metadata: map[string]string{
					"test.rescue.io/key": "foo",
					"test.rescue.io/zer": "tru",
					"test.rescue.io/sin": "baz",
				},
			},
		}

		err = eon.Create(tas)
		if err != nil {
			t.Fatal(err)
		}
	}

	{
		tas := &task.Task{
			Obj: task.TaskObj{
				Metadata: map[string]string{
					"test.rescue.io/key": "foo",
				},
			},
		}

		err = etw.Create(tas)
		if err != nil {
			t.Fatal(err)
		}
	}

	var lis task.Tasks
	{
		tas := &task.Task{
			Obj: task.TaskObj{
				Metadata: map[string]string{
					"test.rescue.io/key": "foo",
				},
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
			Obj: task.TaskObj{
				Metadata: map[string]string{
					"test.rescue.io/key": "foo",
				},
			},
		}

		lis, err = etw.Lister(tas)
		if err != nil {
			t.Fatal(err)
		}
	}

	{
		if len(lis) != 1 {
			t.Fatal("expected 1 tasks listed")
		}
	}

	var tas *task.Task
	{
		tas, err = eon.Search()
		if err != nil {
			t.Fatal(err)
		}

		if tas.Obj.Metadata["test.rescue.io/key"] != "foo" {
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

		if tas.Obj.Metadata["test.rescue.io/key"] != "foo" {
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

		if tas.Obj.Metadata["test.rescue.io/key"] != "foo" {
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
