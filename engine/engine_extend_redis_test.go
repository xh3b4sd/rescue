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

	var lis []*task.Task
	{
		lis, err = eon.Lister(All())
		if err != nil {
			t.Fatal(err)
		}
	}

	{
		if len(lis) != 1 {
			t.Fatal("expected 1 task listed")
		}
		if lis[0].Core.Exi().Worker() {
			t.Fatal("worker label must not be set")
		}
		if lis[0].Meta.Get("test.api.io/key") != "foo" {
			t.Fatal("scheduling failed")
		}
	}

	// Tasks must not be extended if both worker labels are empty, the one within
	// the queue, and the one given as argument to Engine.Extend.
	{
		err = etw.Extend(lis[0])
		if !IsTaskOutdated(err) {
			t.Fatal("task must be extended by owner")
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