//go:build redis

package engine

import (
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/xh3b4sd/budget/v3/pkg/breaker"
	"github.com/xh3b4sd/logger"
	"github.com/xh3b4sd/redigo"
	"github.com/xh3b4sd/redigo/pkg/client"
	"github.com/xh3b4sd/rescue/task"
	"github.com/xh3b4sd/rescue/timer"
	"github.com/xh3b4sd/tracer"
)

func Test_Engine_Lifecycle_Cron_Failure(t *testing.T) {
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

	var tim *timer.Timer
	{
		tim = timer.New()
	}

	var eon *Engine
	{
		eon = New(Config{
			Logger: logger.Fake(),
			Redigo: red,
			Timer:  tim,
			Worker: "eon",
		})
	}

	var etw *Engine
	{
		etw = New(Config{
			Logger: logger.Fake(),
			Redigo: red,
			Timer:  tim,
			Worker: "etw",
		})
	}

	{
		tim.Setter(func() time.Time {
			return musTim("2023-09-28T14:23:24.161982Z")
		})
	}

	{
		exi, err := eon.Exists(&task.Task{Cron: &task.Cron{task.Aevery: "*"}})
		if err != nil {
			t.Fatal(err)
		}

		if exi {
			t.Fatal("expected", false, "got", true)
		}
	}

	{
		exi, err := eon.Exists(&task.Task{Cron: &task.Cron{task.Aevery: "hour"}})
		if err != nil {
			t.Fatal(err)
		}

		if exi {
			t.Fatal("expected", false, "got", true)
		}
	}

	var tas *task.Task
	{
		tas = &task.Task{
			Cron: &task.Cron{
				task.Aevery: "hour",
			},
			Meta: &task.Meta{
				"test.api.io/key": "foo",
			},
		}
	}

	{
		err = eon.Create(tas)
		if err != nil {
			t.Fatal(err)
		}
	}

	{
		exi, err := eon.Exists(&task.Task{Cron: &task.Cron{task.Aevery: "*"}})
		if err != nil {
			t.Fatal(err)
		}

		if !exi {
			t.Fatal("expected", true, "got", false)
		}
	}

	{
		exi, err := eon.Exists(&task.Task{Cron: &task.Cron{task.Aevery: "hour"}})
		if err != nil {
			t.Fatal(err)
		}

		if !exi {
			t.Fatal("expected", true, "got", false)
		}
	}

	var lis []*task.Task
	{
		lis, err = eon.Lister(&task.Task{Cron: tas.Cron.All(task.Aevery)})
		if err != nil {
			t.Fatal(err)
		}
	}

	{
		if len(lis) != 1 {
			t.Fatal("expected 1 task listed")
		}
		if lis[0].Cron.Get().Aevery() != "hour" {
			t.Fatal("scheduling failed")
		}
		if !lis[0].Cron.Get().TickM1().Equal(musTim("2023-09-28T14:00:00.000000Z")) {
			t.Fatal("scheduling failed")
		}
		if !lis[0].Cron.Get().TickP1().Equal(musTim("2023-09-28T15:00:00.000000Z")) {
			t.Fatal("scheduling failed")
		}
		if lis[0].Meta.Get("test.api.io/key") != "foo" {
			t.Fatal("scheduling failed")
		}
	}

	{
		err = eon.Ticker()
		if err != nil {
			t.Fatal("expected", nil, "got", err)
		}
		err = eon.Ticker()
		if err != nil {
			t.Fatal("expected", nil, "got", err)
		}
	}

	{
		lis, err = eon.Lister(&task.Task{Cron: tas.Cron.All(task.Aevery)})
		if err != nil {
			t.Fatal("expected", nil, "got", err)
		}
	}

	{
		if len(lis) != 1 {
			t.Fatal("expected", 1, "got", len(lis))
		}
	}

	{
		exi, err := eon.Exists(&task.Task{Root: &task.Root{task.Object: "*"}})
		if err != nil {
			t.Fatal(err)
		}

		if exi {
			t.Fatal("expected", false, "got", true)
		}
	}

	{
		tim.Setter(func() time.Time {
			return musTim("2023-09-28T15:00:00.161982Z")
		})
	}

	{
		err = eon.Ticker()
		if err != nil {
			t.Fatal("expected", nil, "got", err)
		}
	}

	{
		exi, err := eon.Exists(&task.Task{Root: &task.Root{task.Object: "*"}})
		if err != nil {
			t.Fatal(err)
		}

		if !exi {
			t.Fatal("expected", true, "got", false)
		}
	}

	{
		lis, err = eon.Lister(&task.Task{Cron: tas.Cron.All(task.Aevery)})
		if err != nil {
			t.Fatal("expected", nil, "got", err)
		}
	}

	{
		if len(lis) != 1 {
			t.Fatal("expected 1 task listed")
		}
	}

	{
		lis, err = eon.Lister(All())
		if err != nil {
			t.Fatal("expected", nil, "got", err)
		}
	}

	{
		if len(lis) != 2 {
			t.Fatal("expected", 2, "got", len(lis))
		}
		if lis[0].Meta.Get("test.api.io/key") != "foo" {
			t.Fatal("scheduling failed")
		}
		if lis[0].Cron.Get().Aevery() != "hour" {
			t.Fatal("expected", "hour", "got", lis[0].Cron.Get().Aevery())
		}
		if !lis[0].Cron.Get().TickM1().Equal(musTim("2023-09-28T14:00:00.000000Z")) {
			t.Fatal("scheduling failed")
		}
		if !lis[0].Cron.Get().TickP1().Equal(musTim("2023-09-28T16:00:00.000000Z")) {
			t.Fatal("scheduling failed")
		}
		if lis[0].Root != nil {
			t.Fatal("expected", nil, "got", lis[0].Root)
		}
		if lis[1].Meta.Get("test.api.io/key") != "foo" {
			t.Fatal("scheduling failed")
		}
		if lis[1].Core.Exi().Worker() {
			t.Fatal("scheduling failed")
		}
		if lis[1].Cron != nil {
			t.Fatal("scheduling failed")
		}
		if lis[1].Root.Get(task.Object) != lis[0].Core.Map().Object() {
			t.Fatal("scheduled task must define root for task template")
		}
	}

	{
		tas, err = etw.Search()
		if err != nil {
			t.Fatal("expected", nil, "got", err)
		}
	}

	{
		if tas.Meta.Get("test.api.io/key") != "foo" {
			t.Fatal("scheduling failed")
		}
		if tas.Cron != nil {
			t.Fatal("scheduling failed")
		}
		if tas.Root.Get(task.Object) != lis[0].Core.Map().Object() {
			t.Fatal("expected", lis[0].Core.Map().Object(), "got", tas.Root.Get(task.Object))
		}
	}

	{
		lis, err = eon.Lister(All())
		if err != nil {
			t.Fatal("expected", nil, "got", err)
		}
	}

	{
		if len(lis) != 2 {
			t.Fatal("expected", 2, "got", len(lis))
		}
		if lis[0].Meta.Get("test.api.io/key") != "foo" {
			t.Fatal("scheduling failed")
		}
		if lis[0].Cron.Get().Aevery() != "hour" {
			t.Fatal("scheduling failed")
		}
		if !lis[0].Cron.Get().TickM1().Equal(musTim("2023-09-28T14:00:00.000000Z")) {
			t.Fatal("scheduling failed")
		}
		if !lis[0].Cron.Get().TickP1().Equal(musTim("2023-09-28T16:00:00.000000Z")) {
			t.Fatal("scheduling failed")
		}
		if lis[0].Root != nil {
			t.Fatal("expected", nil, "got", lis[0].Root)
		}
		if lis[1].Meta.Get("test.api.io/key") != "foo" {
			t.Fatal("scheduling failed")
		}
		if !lis[1].Core.Exi().Worker() {
			t.Fatal("scheduling failed")
		}
		if lis[1].Cron != nil {
			t.Fatal("scheduling failed")
		}
		if lis[1].Root.Get(task.Object) != lis[0].Core.Map().Object() {
			t.Fatal("scheduled task must define root for task template")
		}
	}

	{
		tim.Setter(func() time.Time {
			return musTim("2023-09-28T16:00:00.161982Z")
		})
	}

	{
		err = eon.Ticker()
		if err != nil {
			t.Fatal(err)
		}
	}

	{
		lis, err = eon.Lister(All())
		if err != nil {
			t.Fatal("expected", nil, "got", err)
		}
	}

	{
		if len(lis) != 2 {
			t.Fatal("expected", 2, "got", len(lis))
		}
		if lis[0].Meta.Get("test.api.io/key") != "foo" {
			t.Fatal("scheduling failed")
		}
		if lis[0].Cron.Get().Aevery() != "hour" {
			t.Fatal("scheduling failed")
		}
		if !lis[0].Cron.Get().TickM1().Equal(musTim("2023-09-28T14:00:00.000000Z")) {
			t.Fatal("scheduling failed")
		}
		if !lis[0].Cron.Get().TickP1().Equal(musTim("2023-09-28T17:00:00.000000Z")) { // tick+1 moved forward
			t.Fatal("scheduling failed")
		}
		if lis[0].Root != nil {
			t.Fatal("expected", nil, "got", lis[0].Root)
		}
		if lis[1].Meta.Get("test.api.io/key") != "foo" {
			t.Fatal("scheduling failed")
		}
		if !lis[1].Core.Exi().Worker() {
			t.Fatal("scheduling failed")
		}
		if lis[1].Cron != nil {
			t.Fatal("scheduling failed")
		}
		if lis[1].Root.Get(task.Object) != lis[0].Core.Map().Object() {
			t.Fatal("scheduled task must define root for task template")
		}
	}

	{
		tim.Setter(func() time.Time {
			return musTim("2023-09-28T17:00:00.161982Z")
		})
	}

	{
		err = eon.Ticker()
		if err != nil {
			t.Fatal(err)
		}
	}

	{
		lis, err = eon.Lister(All())
		if err != nil {
			t.Fatal("expected", nil, "got", err)
		}
	}

	{
		if len(lis) != 2 {
			t.Fatal("expected", 2, "got", len(lis))
		}
		if lis[0].Meta.Get("test.api.io/key") != "foo" {
			t.Fatal("scheduling failed")
		}
		if lis[0].Cron.Get().Aevery() != "hour" {
			t.Fatal("scheduling failed")
		}
		if !lis[0].Cron.Get().TickM1().Equal(musTim("2023-09-28T14:00:00.000000Z")) {
			t.Fatal("scheduling failed")
		}
		if !lis[0].Cron.Get().TickP1().Equal(musTim("2023-09-28T18:00:00.000000Z")) { // tick+1 moved forward
			t.Fatal("scheduling failed")
		}
		if lis[0].Root != nil {
			t.Fatal("expected", nil, "got", lis[0].Root)
		}
		if lis[1].Meta.Get("test.api.io/key") != "foo" {
			t.Fatal("scheduling failed")
		}
		if !lis[1].Core.Exi().Worker() {
			t.Fatal("scheduling failed")
		}
		if lis[1].Cron != nil {
			t.Fatal("scheduling failed")
		}
		if lis[1].Root.Get(task.Object) != lis[0].Core.Map().Object() {
			t.Fatal("scheduled task must define root for task template")
		}
	}

	{
		err = eon.Delete(tas)
		if err != nil {
			t.Fatal(err)
		}
	}

	{
		err = eon.Ticker()
		if err != nil {
			t.Fatal(err)
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
			t.Fatal("expected 1 task listed")
		}
		if lis[0].Meta.Get("test.api.io/key") != "foo" {
			t.Fatal("scheduling failed")
		}
		if lis[0].Cron.Get().Aevery() != "hour" {
			t.Fatal("scheduling failed")
		}
		if !lis[0].Cron.Get().TickM1().Equal(musTim("2023-09-28T17:00:00.000000Z")) {
			t.Fatal("scheduling failed")
		}
		if !lis[0].Cron.Get().TickP1().Equal(musTim("2023-09-28T18:00:00.000000Z")) {
			t.Fatal("scheduling failed")
		}
		if lis[0].Root != nil {
			t.Fatal("task template must not define root")
		}
	}
}

func Test_Engine_Lifecycle_Cron_Resolve(t *testing.T) {
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

	var tim *timer.Timer
	{
		tim = timer.New()
	}

	var eon *Engine
	{
		eon = New(Config{
			Logger: logger.Fake(),
			Redigo: red,
			Timer:  tim,
			Worker: "eon",
		})
	}

	var etw *Engine
	{
		etw = New(Config{
			Logger: logger.Fake(),
			Redigo: red,
			Timer:  tim,
			Worker: "etw",
		})
	}

	{
		tim.Setter(func() time.Time {
			return musTim("2023-09-28T14:23:24.161982Z")
		})
	}

	{
		exi, err := eon.Exists(&task.Task{Cron: &task.Cron{task.Aevery: "*"}})
		if err != nil {
			t.Fatal(err)
		}

		if exi {
			t.Fatal("expected", false, "got", true)
		}
	}

	{
		exi, err := eon.Exists(&task.Task{Cron: &task.Cron{task.Aevery: "hour"}})
		if err != nil {
			t.Fatal(err)
		}

		if exi {
			t.Fatal("expected", false, "got", true)
		}
	}

	var tas *task.Task
	{
		tas = &task.Task{
			Cron: &task.Cron{
				task.Aevery: "hour",
			},
			Meta: &task.Meta{
				"test.api.io/key": "foo",
			},
		}
	}

	{
		err = eon.Create(tas)
		if err != nil {
			t.Fatal(err)
		}
	}

	{
		exi, err := eon.Exists(&task.Task{Cron: &task.Cron{task.Aevery: "*"}})
		if err != nil {
			t.Fatal(err)
		}

		if !exi {
			t.Fatal("expected", true, "got", false)
		}
	}

	{
		exi, err := eon.Exists(&task.Task{Cron: &task.Cron{task.Aevery: "hour"}})
		if err != nil {
			t.Fatal(err)
		}

		if !exi {
			t.Fatal("expected", true, "got", false)
		}
	}

	var lis []*task.Task
	{
		lis, err = eon.Lister(&task.Task{Cron: tas.Cron.All(task.Aevery)})
		if err != nil {
			t.Fatal(err)
		}
	}

	{
		if len(lis) != 1 {
			t.Fatal("expected 1 task listed")
		}
		if lis[0].Cron.Get().Aevery() != "hour" {
			t.Fatal("scheduling failed")
		}
		if !lis[0].Cron.Get().TickM1().Equal(musTim("2023-09-28T14:00:00.000000Z")) {
			t.Fatal("scheduling failed")
		}
		if !lis[0].Cron.Get().TickP1().Equal(musTim("2023-09-28T15:00:00.000000Z")) {
			t.Fatal("scheduling failed")
		}
		if lis[0].Meta.Get("test.api.io/key") != "foo" {
			t.Fatal("scheduling failed")
		}
	}

	{
		err = eon.Ticker()
		if err != nil {
			t.Fatal("expected", nil, "got", err)
		}
		err = eon.Ticker()
		if err != nil {
			t.Fatal("expected", nil, "got", err)
		}
	}

	{
		lis, err = eon.Lister(&task.Task{Cron: tas.Cron.All(task.Aevery)})
		if err != nil {
			t.Fatal("expected", nil, "got", err)
		}
	}

	{
		if len(lis) != 1 {
			t.Fatal("expected", 1, "got", len(lis))
		}
	}

	{
		exi, err := eon.Exists(&task.Task{Root: &task.Root{task.Object: "*"}})
		if err != nil {
			t.Fatal(err)
		}

		if exi {
			t.Fatal("expected", false, "got", true)
		}
	}

	{
		tim.Setter(func() time.Time {
			return musTim("2023-09-28T15:00:00.161982Z")
		})
	}

	{
		err = eon.Ticker()
		if err != nil {
			t.Fatal("expected", nil, "got", err)
		}
	}

	{
		exi, err := eon.Exists(&task.Task{Root: &task.Root{task.Object: "*"}})
		if err != nil {
			t.Fatal(err)
		}

		if !exi {
			t.Fatal("expected", true, "got", false)
		}
	}

	{
		lis, err = eon.Lister(&task.Task{Cron: tas.Cron.All(task.Aevery)})
		if err != nil {
			t.Fatal("expected", nil, "got", err)
		}
	}

	{
		if len(lis) != 1 {
			t.Fatal("expected 1 task listed")
		}
	}

	{
		lis, err = eon.Lister(All())
		if err != nil {
			t.Fatal("expected", nil, "got", err)
		}
	}

	{
		if len(lis) != 2 {
			t.Fatal("expected", 2, "got", len(lis))
		}
		if lis[0].Meta.Get("test.api.io/key") != "foo" {
			t.Fatal("scheduling failed")
		}
		if lis[0].Cron.Get().Aevery() != "hour" {
			t.Fatal("expected", "hour", "got", lis[0].Cron.Get().Aevery())
		}
		if !lis[0].Cron.Get().TickM1().Equal(musTim("2023-09-28T14:00:00.000000Z")) {
			t.Fatal("scheduling failed")
		}
		if !lis[0].Cron.Get().TickP1().Equal(musTim("2023-09-28T16:00:00.000000Z")) {
			t.Fatal("scheduling failed")
		}
		if lis[0].Root != nil {
			t.Fatal("expected", nil, "got", lis[0].Root)
		}
		if lis[1].Meta.Get("test.api.io/key") != "foo" {
			t.Fatal("scheduling failed")
		}
		if lis[1].Core.Exi().Worker() {
			t.Fatal("scheduling failed")
		}
		if lis[1].Cron != nil {
			t.Fatal("scheduling failed")
		}
		if lis[1].Root.Get(task.Object) != lis[0].Core.Map().Object() {
			t.Fatal("scheduled task must define root for task template")
		}
	}

	{
		tas, err = etw.Search()
		if err != nil {
			t.Fatal("expected", nil, "got", err)
		}
	}

	{
		if tas.Meta.Get("test.api.io/key") != "foo" {
			t.Fatal("scheduling failed")
		}
		if tas.Cron != nil {
			t.Fatal("scheduling failed")
		}
		if tas.Root.Get(task.Object) != lis[0].Core.Map().Object() {
			t.Fatal("expected", lis[0].Core.Map().Object(), "got", tas.Root.Get(task.Object))
		}
	}

	{
		lis, err = eon.Lister(All())
		if err != nil {
			t.Fatal("expected", nil, "got", err)
		}
	}

	{
		if len(lis) != 2 {
			t.Fatal("expected", 2, "got", len(lis))
		}
		if lis[0].Meta.Get("test.api.io/key") != "foo" {
			t.Fatal("scheduling failed")
		}
		if lis[0].Cron.Get().Aevery() != "hour" {
			t.Fatal("scheduling failed")
		}
		if !lis[0].Cron.Get().TickM1().Equal(musTim("2023-09-28T14:00:00.000000Z")) {
			t.Fatal("scheduling failed")
		}
		if !lis[0].Cron.Get().TickP1().Equal(musTim("2023-09-28T16:00:00.000000Z")) {
			t.Fatal("scheduling failed")
		}
		if lis[0].Root != nil {
			t.Fatal("expected", nil, "got", lis[0].Root)
		}
		if lis[1].Meta.Get("test.api.io/key") != "foo" {
			t.Fatal("scheduling failed")
		}
		if !lis[1].Core.Exi().Worker() {
			t.Fatal("scheduling failed")
		}
		if lis[1].Cron != nil {
			t.Fatal("scheduling failed")
		}
		if lis[1].Root.Get(task.Object) != lis[0].Core.Map().Object() {
			t.Fatal("scheduled task must define root for task template")
		}
	}

	{
		err = eon.Delete(tas)
		if err != nil {
			t.Fatal(err)
		}
	}

	{
		err = eon.Ticker()
		if err != nil {
			t.Fatal(err)
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
			t.Fatal("expected 1 task listed")
		}
		if lis[0].Meta.Get("test.api.io/key") != "foo" {
			t.Fatal("scheduling failed")
		}
		if lis[0].Cron.Get().Aevery() != "hour" {
			t.Fatal("scheduling failed")
		}
		if !lis[0].Cron.Get().TickM1().Equal(musTim("2023-09-28T15:00:00.000000Z")) {
			t.Fatal("scheduling failed")
		}
		if !lis[0].Cron.Get().TickP1().Equal(musTim("2023-09-28T16:00:00.000000Z")) {
			t.Fatal("scheduling failed")
		}
		if lis[0].Root != nil {
			t.Fatal("task template must not define root")
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

func musTim(str string) time.Time {
	tim, err := time.Parse("2006-01-02T15:04:05.999999Z", str)
	if err != nil {
		panic(err)
	}

	return tim
}
