//go:build redis

package conformance

import (
	"testing"
	"time"

	"github.com/xh3b4sd/logger"
	"github.com/xh3b4sd/redigo"
	"github.com/xh3b4sd/rescue"
	"github.com/xh3b4sd/rescue/engine"
	"github.com/xh3b4sd/rescue/task"
)

func Test_Engine_Extend(t *testing.T) {
	var err error

	var eon rescue.Interface
	{
		eon = engine.New(engine.Config{
			Expiry: time.Second,
			Logger: logger.Fake(),
			Redigo: prgAll(redigo.Default()),
			Worker: "eon",
		})
	}

	var etw rescue.Interface
	{
		etw = engine.New(engine.Config{
			Expiry: time.Second,
			Logger: logger.Fake(),
			Redigo: prgAll(redigo.Default()),
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
		lis, err = eon.Lister(engine.All())
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
		if !engine.IsTaskOutdated(err) {
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
		if !engine.IsTaskNotFound(err) {
			t.Fatal("expected", "taskNotFoundError", "got", err)
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
		if !engine.IsTaskNotFound(err) {
			t.Fatal("expected", "taskNotFoundError", "got", err)
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
		if !engine.IsTaskNotFound(err) {
			t.Fatal("expected", "taskNotFoundError", "got", err)
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
		if !engine.IsTaskNotFound(err) {
			t.Fatal("expected", "taskNotFoundError", "got", err)
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
		if !engine.IsTaskNotFound(err) {
			t.Fatal("expected", "taskNotFoundError", "got", err)
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
		if !engine.IsTaskNotFound(err) {
			t.Fatal("expected", "taskNotFoundError", "got", err)
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
		if !engine.IsTaskOutdated(err) {
			t.Fatal("task must be extended by owner")
		}
	}

	{
		err = etw.Delete(tas)
		if !engine.IsTaskOutdated(err) {
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
		if !engine.IsTaskNotFound(err) {
			t.Fatal("expected", "taskNotFoundError", "got", err)
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

		_, err = etw.Search()
		if !engine.IsTaskNotFound(err) {
			t.Fatal("expected", "taskNotFoundError", "got", err)
		}
	}
}
