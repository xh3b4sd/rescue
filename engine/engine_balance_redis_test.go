//go:build redis

package engine

import (
	"strconv"
	"testing"

	"github.com/xh3b4sd/logger"
	"github.com/xh3b4sd/redigo"
	"github.com/xh3b4sd/rescue/task"
)

func Test_Engine_Balance(t *testing.T) {
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

	var eth *Engine
	{
		eth = New(Config{
			Logger: logger.Fake(),
			Redigo: red,
			Worker: "eth",
		})
	}

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
