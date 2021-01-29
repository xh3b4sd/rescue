// +build redis

package engine

import (
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/xh3b4sd/logger/fake"
	"github.com/xh3b4sd/redigo"
	"github.com/xh3b4sd/redigo/pkg/client"
	"github.com/xh3b4sd/tracer"

	"github.com/xh3b4sd/rescue/pkg/task"
)

func Test_Engine_Lifecycle(t *testing.T) {
	var err error

	var red redigo.Interface
	{
		c := client.Config{
			Kind: client.KindSingle,
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
		c := Config{
			Logger: fake.New(),
			Redigo: red,

			Owner: "eon",
		}

		eon, err = New(c)
		if err != nil {
			t.Fatal(err)
		}
	}

	var etw *Engine
	{
		c := Config{
			Logger: fake.New(),
			Redigo: red,

			Owner: "etw",
		}

		etw, err = New(c)
		if err != nil {
			t.Fatal(err)
		}
	}

	{
		tsk := &task.Task{
			Obj: task.TaskObj{
				Metadata: map[string]string{
					"test.rescue.io/key": "foo",
				},
			},
		}

		err = eon.Create(tsk)
		if err != nil {
			t.Fatal(err)
		}
	}

	{
		tsk := &task.Task{
			Obj: task.TaskObj{
				Metadata: map[string]string{
					"test.rescue.io/key": "bar",
				},
			},
		}

		err = etw.Create(tsk)
		if err != nil {
			t.Fatal(err)
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

			tsk, err := eon.Search()
			if err != nil {
				erc <- tracer.Mask(err)
				return
			}

			s = append(s, tsk.Obj.Metadata["test.rescue.io/key"])

			err = eon.Delete(tsk)
			if err != nil {
				erc <- tracer.Mask(err)
				return
			}
		}()

		go func() {
			defer w.Done()

			tsk, err := eon.Search()
			if err != nil {
				erc <- tracer.Mask(err)
				return
			}

			s = append(s, tsk.Obj.Metadata["test.rescue.io/key"])

			err = etw.Delete(tsk)
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
			_, err = eon.Search()
			if !IsNoTask(err) {
				erc <- fmt.Errorf("queue must be empty")
				return
			}

			_, err = etw.Search()
			if !IsNoTask(err) {
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

func Test_Engine_Expire(t *testing.T) {
	var err error

	var red redigo.Interface
	{
		c := client.Config{
			Kind: client.KindSingle,
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
		c := Config{
			Logger: fake.New(),
			Redigo: red,

			Owner:  "eon",
			Expire: time.Millisecond,
		}

		eon, err = New(c)
		if err != nil {
			t.Fatal(err)
		}
	}

	var etw *Engine
	{
		c := Config{
			Logger: fake.New(),
			Redigo: red,

			Owner:  "etw",
			Expire: time.Millisecond,
		}

		etw, err = New(c)
		if err != nil {
			t.Fatal(err)
		}
	}

	{
		tsk := &task.Task{
			Obj: task.TaskObj{
				Metadata: map[string]string{
					"test.rescue.io/key": "foo",
				},
			},
		}

		err = eon.Create(tsk)
		if err != nil {
			t.Fatal(err)
		}
	}

	{
		tsk := &task.Task{
			Obj: task.TaskObj{
				Metadata: map[string]string{
					"test.rescue.io/key": "bar",
				},
			},
		}

		err = etw.Create(tsk)
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
		tsk, err := etw.Search()
		if err != nil {
			t.Fatal(err)
		}

		s = append(s, tsk.Obj.Metadata["test.rescue.io/key"])

		err = etw.Delete(tsk)
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
		tsk, err := etw.Search()
		if err != nil {
			t.Fatal(err)
		}

		s = append(s, tsk.Obj.Metadata["test.rescue.io/key"])

		err = etw.Delete(tsk)
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
		if !IsNoTask(err) {
			t.Fatal("queue must be empty")
		}

		_, err = etw.Search()
		if !IsNoTask(err) {
			t.Fatal("queue must be empty")
		}
	}
}
