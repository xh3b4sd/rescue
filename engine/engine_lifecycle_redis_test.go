//go:build redis

package engine

import (
	"fmt"
	"reflect"
	"sync"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/xh3b4sd/logger"
	"github.com/xh3b4sd/redigo"
	"github.com/xh3b4sd/rescue/task"
	"github.com/xh3b4sd/rescue/timer"
	"github.com/xh3b4sd/tracer"
)

func Test_Engine_Lifecycle_Cron_3Days(t *testing.T) {
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

	var tim *timer.Timer
	{
		tim = timer.New()
	}

	// The engine and the ticker instances are configured with a time for the end
	// of the year.
	{
		tim.Setter(func() time.Time {
			return musTim("2022-12-31T14:23:24.161982Z")
		})
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

	{
		exi, err := eon.Exists(&task.Task{Cron: &task.Cron{task.Aevery: "*"}})
		if err != nil {
			t.Fatal("expected", true, "got", false)
		}

		if exi {
			t.Fatal("expected", false, "got", true)
		}
	}

	// We set the task template to schedule every 3 days, which from the point of
	// view of "now" would imply the next tick to be the 2nd of January.
	{
		exi, err := eon.Exists(&task.Task{Cron: &task.Cron{task.Aevery: "3 days"}})
		if err != nil {
			t.Fatal("expected", true, "got", false)
		}

		if exi {
			t.Fatal("expected", false, "got", true)
		}
	}

	var tas *task.Task
	{
		tas = &task.Task{
			Core: &task.Core{
				task.Method: task.MthdAll,
			},
			Cron: &task.Cron{
				task.Aevery: "3 days",
			},
			Meta: &task.Meta{
				"test.api.io/key": "foo",
			},
		}
	}

	{
		err = eon.Create(tas)
		if err != nil {
			t.Fatal("expected", true, "got", false)
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
		exi, err := eon.Exists(&task.Task{Cron: &task.Cron{task.Aevery: "3 days"}})
		if err != nil {
			t.Fatal(err)
		}

		if !exi {
			t.Fatal("expected", true, "got", false)
		}
	}

	// Search for all task templates defining Task.Cron.
	var lis []*task.Task
	{
		lis, err = eon.Lister(&task.Task{Cron: tas.Cron.All(task.Aevery)})
		if err != nil {
			t.Fatal(err)
		}
	}

	{
		if len(lis) != 1 {
			t.Fatal("expected", 1, "got", len(lis))
		}
	}

	// Ensure the task template defining Task.Cron contains the correct scheduling
	// information. Also ensure that the delivery method is set to "all".
	{
		var tas *task.Task
		{
			tas = lis[0]
		}

		var exp *task.Task
		{
			exp = &task.Task{
				Core: tas.Core,
				Cron: &task.Cron{
					task.Aevery: "3 days",
					task.TickM1: "2022-12-30T00:00:00Z",
					task.TickP1: "2023-01-02T00:00:00Z",
				},
				Meta: &task.Meta{
					"test.api.io/key": "foo",
				},
			}
		}

		{
			if !reflect.DeepEqual(tas, exp) {
				t.Fatalf("\n\n%s\n", cmp.Diff(exp, tas))
			}
			if tas.Core.Get().Method() != task.MthdAll {
				t.Fatal("expected", task.MthdAll, "got", tas.Core.Get().Method())
			}
		}
	}

	// The next ticker execution happens 1st of January. Here we will test for
	// schedule continuation. Based on the previous year's ticker calculation, the
	// current interval would end at the 2nd of January. Without carrying over
	// properly, the new year's ticker on its own would calculate the end of the
	// current interval to be on the 4th of January. And we do not want that.
	{
		tim.Setter(func() time.Time {
			return musTim("2023-01-01T14:23:24.161982Z")
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
			t.Fatal(err)
		}
	}

	{
		if len(lis) != 1 {
			t.Fatal("expected", 1, "got", len(lis))
		}
	}

	// After moving forward into the next year, the schedule continues properly
	// and did so far not change.
	{
		var tas *task.Task
		{
			tas = lis[0]
		}

		var exp *task.Task
		{
			exp = &task.Task{
				Core: tas.Core,
				Cron: &task.Cron{
					task.Aevery: "3 days",
					task.TickM1: "2022-12-30T00:00:00Z",
					task.TickP1: "2023-01-02T00:00:00Z",
				},
				Meta: &task.Meta{
					"test.api.io/key": "foo",
				},
			}
		}

		{
			if !reflect.DeepEqual(tas, exp) {
				t.Fatalf("\n\n%s\n", cmp.Diff(exp, tas))
			}
			if tas.Core.Get().Method() != task.MthdAll {
				t.Fatal("expected", task.MthdAll, "got", tas.Core.Get().Method())
			}
		}
	}

	// Move past tick+1 in order to emit the new scheduled task.
	{
		tim.Setter(func() time.Time {
			return musTim("2023-01-02T00:00:00.161982Z")
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
				Cron: &task.Cron{
					task.Aevery: "3 days",
					task.TickM1: "2023-01-02T00:00:00Z",
					task.TickP1: "2023-01-05T00:00:00Z",
				},
				Meta: &task.Meta{
					"test.api.io/key": "foo",
				},
			}
		}

		{
			if !reflect.DeepEqual(tas, exp) {
				t.Fatalf("\n\n%s\n", cmp.Diff(exp, tas))
			}
			if tas.Core.Get().Method() != task.MthdAll {
				t.Fatal("expected", task.MthdAll, "got", tas.Core.Get().Method())
			}
		}
	}

	// Ensure the scheduled task is delivered to all workers within the network as
	// defined by the delivery method "all".
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
					"test.api.io/key": "foo",
				},
				Root: &task.Root{
					task.Object: lis[0].Core.Map().Object(),
				},
			}
		}

		{
			if !reflect.DeepEqual(tas, exp) {
				t.Fatalf("\n\n%s\n", cmp.Diff(exp, tas))
			}
			if tas.Core.Get().Method() != task.MthdAll {
				t.Fatal("expected", task.MthdAll, "got", tas.Core.Get().Method())
			}
		}
	}
}

func Test_Engine_Lifecycle_Cron_Failure(t *testing.T) {
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
			t.Fatal("expected", 1, "got", len(lis))
		}
		if lis[0].Cron.Get().Aevery() != "hour" {
			t.Fatal("expected", "hour", "got", lis[0].Cron.Get().Aevery())
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

	// Verify the raw string format to make sure the ticker layout is persisted as
	// expected.
	{
		if lis[0].Cron.Map().TickM1() != "2023-09-28T14:00:00Z" {
			t.Fatal("expected", "2023-09-28T14:00:00Z", "got", lis[0].Cron.Map().TickM1())
		}
		if lis[0].Cron.Map().TickP1() != "2023-09-28T15:00:00Z" {
			t.Fatal("expected", "2023-09-28T15:00:00Z", "got", lis[0].Cron.Map().TickP1())
		}
	}

	{
		err = eon.Ticker()
		if err != nil {
			t.Fatal(err)
		}
		err = eon.Ticker()
		if err != nil {
			t.Fatal(err)
		}
	}

	{
		lis, err = eon.Lister(&task.Task{Cron: tas.Cron.All(task.Aevery)})
		if err != nil {
			t.Fatal(err)
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
			t.Fatal(err)
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
			t.Fatal(err)
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
			t.Fatal(err)
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

	// Verify the raw string format to make sure the ticker layout is persisted as
	// expected.
	{
		if lis[0].Cron.Map().TickM1() != "2023-09-28T14:00:00Z" {
			t.Fatal("expected", "2023-09-28T14:00:00Z", "got", lis[0].Cron.Map().TickM1())
		}
		if lis[0].Cron.Map().TickP1() != "2023-09-28T16:00:00Z" {
			t.Fatal("expected", "2023-09-28T16:00:00Z", "got", lis[0].Cron.Map().TickP1())
		}
	}

	{
		tas, err = etw.Search()
		if err != nil {
			t.Fatal(err)
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
			t.Fatal(err)
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

	// Verify the raw string format to make sure the ticker layout is persisted as
	// expected.
	{
		if lis[0].Cron.Map().TickM1() != "2023-09-28T14:00:00Z" {
			t.Fatal("expected", "2023-09-28T14:00:00Z", "got", lis[0].Cron.Map().TickM1())
		}
		if lis[0].Cron.Map().TickP1() != "2023-09-28T16:00:00Z" {
			t.Fatal("expected", "2023-09-28T16:00:00Z", "got", lis[0].Cron.Map().TickP1())
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
			t.Fatal(err)
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

	// Verify the raw string format to make sure the ticker layout is persisted as
	// expected.
	{
		if lis[0].Cron.Map().TickM1() != "2023-09-28T14:00:00Z" {
			t.Fatal("expected", "2023-09-28T14:00:00Z", "got", lis[0].Cron.Map().TickM1())
		}
		if lis[0].Cron.Map().TickP1() != "2023-09-28T17:00:00Z" {
			t.Fatal("expected", "2023-09-28T17:00:00Z", "got", lis[0].Cron.Map().TickP1())
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
			t.Fatal(err)
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

	// Verify the raw string format to make sure the ticker layout is persisted as
	// expected.
	{
		if lis[0].Cron.Map().TickM1() != "2023-09-28T14:00:00Z" {
			t.Fatal("expected", "2023-09-28T14:00:00Z", "got", lis[0].Cron.Map().TickM1())
		}
		if lis[0].Cron.Map().TickP1() != "2023-09-28T18:00:00Z" {
			t.Fatal("expected", "2023-09-28T18:00:00Z", "got", lis[0].Cron.Map().TickP1())
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

	// Verify the raw string format to make sure the ticker layout is persisted as
	// expected.
	{
		if lis[0].Cron.Map().TickM1() != "2023-09-28T17:00:00Z" {
			t.Fatal("expected", "2023-09-28T17:00:00Z", "got", lis[0].Cron.Map().TickM1())
		}
		if lis[0].Cron.Map().TickP1() != "2023-09-28T18:00:00Z" {
			t.Fatal("expected", "2023-09-28T18:00:00Z", "got", lis[0].Cron.Map().TickP1())
		}
	}
}

func Test_Engine_Lifecycle_Cron_Resolve(t *testing.T) {
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
			Gate: &task.Gate{
				"test.api.io/k-1": task.Trigger,
			},
			Sync: &task.Sync{
				"test.api.io/lat": "initial",
				"test.api.io/foo": "should-not-change",
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
		if lis[0].Meta.Get("test.api.io/key") != "foo" {
			t.Fatal("scheduling failed")
		}
		if lis[0].Cron.Get().Aevery() != "hour" {
			t.Fatal("expected", "hour", "got", lis[0].Cron.Get().Aevery())
		}
		if !lis[0].Cron.Get().TickM1().Equal(musTim("2023-09-28T14:00:00.000000Z")) {
			t.Fatal("scheduling failed")
		}
		if !lis[0].Cron.Get().TickP1().Equal(musTim("2023-09-28T15:00:00.000000Z")) {
			t.Fatal("scheduling failed")
		}
		if lis[0].Root != nil {
			t.Fatal("expected", nil, "got", lis[0].Root)
		}
		if lis[0].Sync.Len() != 2 {
			t.Fatal("expected", 2, "got", lis[0].Sync.Len())
		}
		if lis[0].Sync.Get("test.api.io/lat") != "initial" {
			t.Fatal("expected", "initial", "got", lis[0].Sync.Get("test.api.io/lat"))
		}
		if lis[0].Sync.Get("test.api.io/foo") != "should-not-change" {
			t.Fatal("expected", "should-not-change", "got", lis[0].Sync.Get("test.api.io/lat"))
		}
	}

	{
		err = eon.Ticker()
		if err != nil {
			t.Fatal(err)
		}
		err = eon.Ticker()
		if err != nil {
			t.Fatal(err)
		}
	}

	{
		lis, err = eon.Lister(&task.Task{Cron: tas.Cron.All(task.Aevery)})
		if err != nil {
			t.Fatal(err)
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
			t.Fatal(err)
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
			t.Fatal(err)
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
			t.Fatal(err)
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
		if lis[0].Sync.Len() != 2 {
			t.Fatal("expected", 2, "got", lis[0].Sync.Len())
		}
		if lis[0].Sync.Get("test.api.io/lat") != "initial" {
			t.Fatal("expected", "initial", "got", lis[0].Sync.Get("test.api.io/lat"))
		}
		if lis[0].Sync.Get("test.api.io/foo") != "should-not-change" {
			t.Fatal("expected", "should-not-change", "got", lis[0].Sync.Get("test.api.io/lat"))
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
		if lis[1].Sync.Len() != 2 {
			t.Fatal("expected", 2, "got", lis[1].Sync.Len())
		}
		if lis[1].Sync.Get("test.api.io/lat") != "initial" {
			t.Fatal("expected", "initial", "got", lis[1].Sync.Get("test.api.io/lat"))
		}
		if lis[1].Sync.Get("test.api.io/foo") != "should-not-change" {
			t.Fatal("expected", "should-not-change", "got", lis[1].Sync.Get("test.api.io/lat"))
		}
	}

	{
		tas, err = etw.Search()
		if err != nil {
			t.Fatal(err)
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
				Cron: &task.Cron{
					task.Aevery: "hour",
					task.TickM1: "2023-09-28T14:00:00Z",
					task.TickP1: "2023-09-28T16:00:00Z",
				},
				Gate: &task.Gate{
					"test.api.io/k-1": task.Trigger,
				},
				Meta: &task.Meta{
					"test.api.io/key": "foo",
				},
				Sync: &task.Sync{
					"test.api.io/foo": "should-not-change",
					"test.api.io/lat": "initial",
				},
			}
		}

		{
			if !reflect.DeepEqual(tas, exp) {
				t.Fatalf("\n\n%s\n", cmp.Diff(exp, tas))
			}
			if tas.Core.Get().Method() != task.MthdAny {
				t.Fatal("expected", task.MthdAny, "got", tas.Core.Get().Method())
			}
		}
	}

	// Verify the scheduled task that should contain Task.Gate, Task.Meta,
	// Task.Root and Task.Sync according to the task template emitting it. Also
	// ensure that the default delivery method "any" is set.
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
				Root: &task.Root{
					task.Object: lis[0].Core.Map().Object(),
				},
				Sync: &task.Sync{
					"test.api.io/foo": "should-not-change",
					"test.api.io/lat": "initial",
				},
			}
		}

		{
			if !reflect.DeepEqual(tas, exp) {
				t.Fatalf("\n\n%s\n", cmp.Diff(exp, tas))
			}
			if tas.Core.Get().Method() != task.MthdAny {
				t.Fatal("expected", task.MthdAny, "got", tas.Core.Get().Method())
			}
		}
	}

	// We modify the data in Task.Sync to verify that the latest pointer of our
	// scheduled task gets propagated to our task template upon deleting the
	// scheduled task.
	{
		tas.Sync.Set("test.api.io/lat", "updated")
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
			t.Fatal("expected", 1, "got", len(lis))
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
				Cron: &task.Cron{
					task.Aevery: "hour",
					task.TickM1: "2023-09-28T15:00:00Z",
					task.TickP1: "2023-09-28T16:00:00Z",
				},
				Gate: &task.Gate{
					"test.api.io/k-1": task.Trigger,
				},
				Meta: &task.Meta{
					"test.api.io/key": "foo",
				},
				Sync: &task.Sync{
					"test.api.io/foo": "should-not-change",
					"test.api.io/lat": "updated",
				},
			}
		}

		{
			if !reflect.DeepEqual(tas, exp) {
				t.Fatalf("\n\n%s\n", cmp.Diff(exp, tas))
			}
			if tas.Core.Get().Method() != task.MthdAny {
				t.Fatal("expected", task.MthdAny, "got", tas.Core.Get().Method())
			}
		}
	}
}

func Test_Engine_Lifecycle_Method_All_Failure(t *testing.T) {
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

	var tim *timer.Timer
	{
		tim = timer.New()
	}

	// The engine is configured with a particular time. This point in time will be
	// set inside the worker process as the pointer for when it started processing
	// tasks.
	{
		tim.Setter(func() time.Time {
			return musTim("2023-10-20T00:00:00Z")
		})
	}

	var eon *Engine
	{
		eon = New(Config{
			Logger: logger.Fake(),
			Redigo: red,
			Timer:  tim,
		})
	}

	// Time advances by 1 minute. So the first task "foo" got created at minute
	// one.
	{
		tim.Setter(func() time.Time {
			return musTim("2023-10-20T00:01:00Z")
		})
	}

	{
		tas := &task.Task{
			Core: &task.Core{
				task.Method: task.MthdAll,
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

	// Time advances by 1 more minute. So the second task "bar" got created at
	// minute two.
	{
		tim.Setter(func() time.Time {
			return musTim("2023-10-20T00:02:00Z")
		})
	}

	{
		tas := &task.Task{
			Core: &task.Core{
				task.Method: task.MthdAll,
			},
			Meta: &task.Meta{
				"test.api.io/key": "bar",
			},
		}

		err = eon.Create(tas)
		if err != nil {
			t.Fatal(err)
		}
	}

	// The worker looks for new tasks and finds task one.
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
				Meta: &task.Meta{
					"test.api.io/key": "foo",
				},
			}
		}

		{
			if !reflect.DeepEqual(tas, exp) {
				t.Fatalf("\n\n%s\n", cmp.Diff(exp, tas))
			}
			if tas.Core.Get().Method() != task.MthdAll {
				t.Fatal("expected", task.MthdAll, "got", tas.Core.Get().Method())
			}
		}
	}

	// For the test here we pretend task one gets stuck or fails for whatever
	// reason, it will expire locally for the worker in memory.

	// The worker looks for another task and finds task two.
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
				Meta: &task.Meta{
					"test.api.io/key": "bar",
				},
			}
		}

		{
			if !reflect.DeepEqual(tas, exp) {
				t.Fatalf("\n\n%s\n", cmp.Diff(exp, tas))
			}
			if tas.Core.Get().Method() != task.MthdAll {
				t.Fatal("expected", task.MthdAll, "got", tas.Core.Get().Method())
			}
		}
	}

	// Time advances some 10 seconds. So task two gets completed without issues.
	{
		tim.Setter(func() time.Time {
			return musTim("2023-10-20T00:02:10Z")
		})
	}

	{
		fmt.Printf("delete %#v\n", 11111111)
		err = eon.Delete(tas)
		if err != nil {
			t.Fatal(err)
		}
	}

	var etw *Engine
	{
		etw = New(Config{
			Logger: logger.Fake(),
			Redigo: red,
			Timer:  tim,
		})
	}

	// A new worker, engine two, came online just now. It cannot find any task
	// available, since all created tasks got created before it came online.
	{
		_, err = etw.Search()
		if !IsTaskNotFound(err) {
			t.Fatal("expected", taskNotFoundError, "got", err)
		}
	}

	// Now the failing task one can still not be received because of its internal
	// expiry.
	{
		tas, err = eon.Search()
		if !IsTaskNotFound(err) {
			t.Fatal("expected", taskNotFoundError, "got", err)
		}
	}

	// Time advances by another 25 seconds. So task one expired in worker memory
	// and should be receivable again.
	{
		tim.Setter(func() time.Time {
			return musTim("2023-10-20T00:02:35Z")
		})
	}

	{
		tas, err = eon.Search()
		if err != nil {
			t.Fatal(err)
		}
	}

	// Engine two can never receive any task because no new tasks got created
	// since it came online.
	{
		_, err = etw.Search()
		if !IsTaskNotFound(err) {
			t.Fatal("expected", taskNotFoundError, "got", err)
		}
	}

	{
		var exp *task.Task
		{
			exp = &task.Task{
				Core: tas.Core,
				Meta: &task.Meta{
					"test.api.io/key": "foo",
				},
			}
		}

		{
			if !reflect.DeepEqual(tas, exp) {
				t.Fatalf("\n\n%s\n", cmp.Diff(exp, tas))
			}
			if tas.Core.Get().Method() != task.MthdAll {
				t.Fatal("expected", task.MthdAll, "got", tas.Core.Get().Method())
			}
		}
	}

	// Time advances by some 5 more seconds. So this time around the worker
	// completed task one, which was failing earlier.
	{
		tim.Setter(func() time.Time {
			return musTim("2023-10-20T00:02:40Z")
		})
	}

	{
		err = eon.Delete(tas)
		if err != nil {
			t.Fatal(err)
		}
	}

	// Since the worker completed all two broadcasted tasks there should not be
	// any more tasks to be received.
	{
		_, err = eon.Search()
		if !IsTaskNotFound(err) {
			t.Fatal("expected", taskNotFoundError, "got", err)
		}
	}

	// Engine two can never receive any task because no new tasks got created
	// since it came online.
	{
		_, err = etw.Search()
		if !IsTaskNotFound(err) {
			t.Fatal("expected", taskNotFoundError, "got", err)
		}
	}

	// Time advances by 6 hours. Even here no more tasks should be available for
	// the worker.
	{
		tim.Setter(func() time.Time {
			return musTim("2023-10-20T06:02:40Z")
		})
	}

	{
		_, err = eon.Search()
		if !IsTaskNotFound(err) {
			t.Fatal("expected", taskNotFoundError, "got", err)
		}
	}

	// Engine two can never receive any task because no new tasks got created
	// since it came online.
	{
		_, err = etw.Search()
		if !IsTaskNotFound(err) {
			t.Fatal("expected", taskNotFoundError, "got", err)
		}
	}

	// Time advances by 9 days. Even here no more tasks should be available for
	// the worker.
	{
		tim.Setter(func() time.Time {
			return musTim("2023-10-29T06:02:40Z")
		})
	}

	{
		_, err = eon.Search()
		if !IsTaskNotFound(err) {
			t.Fatal("expected", taskNotFoundError, "got", err)
		}
	}

	// Engine two can never receive any task because no new tasks got created
	// since it came online.
	{
		_, err = etw.Search()
		if !IsTaskNotFound(err) {
			t.Fatal("expected", taskNotFoundError, "got", err)
		}
	}
}

func Test_Engine_Lifecycle_Race(t *testing.T) {
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
				erc <- fmt.Errorf("%s %#v %s %#v", "expected", taskNotFoundError, "got", err)
				return
			}

			_, err = etw.Search()
			if !IsTaskNotFound(err) {
				erc <- fmt.Errorf("%s %#v %s %#v", "expected", taskNotFoundError, "got", err)
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
			t.Fatal("expected", taskNotFoundError, "got", err)
		}
	}

	{
		_, err = etw.Search()
		if !IsTaskNotFound(err) {
			t.Fatal("expected", taskNotFoundError, "got", err)
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
