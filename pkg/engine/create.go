package engine

import (
	"fmt"
	"strings"
	"time"

	"github.com/xh3b4sd/tracer"

	"github.com/xh3b4sd/rescue/pkg/key"
	"github.com/xh3b4sd/rescue/pkg/task"
)

func (e *Engine) create(tsk *task.Task) error {
	var err error

	{
		if tsk == nil {
			return tracer.Maskf(invalidTaskError, "task must not be nil")
		}
		if len(tsk.Obj.Metadata) == 0 {
			return tracer.Maskf(invalidTaskError, "metadata must not be empty")
		}
		for k := range tsk.Obj.Metadata {
			if strings.HasPrefix(k, "task.rescue.io") {
				return tracer.Maskf(invalidTaskError, "metadata must not contain reserved scheme task.rescue.io")
			}
		}
	}

	// Creating tasks implies certain write operations on the task queue such as
	// adding a new task to a sorted set in redis. Due to such write operations
	// we need to ensure that only one process at a time can write information
	// back to the queue. Otherwise worker behaviour would be inconsistent and
	// the integrity of the queue could not be guaranteed.
	{
		err := e.redigo.Locker().Acquire()
		if err != nil {
			return tracer.Mask(err)
		}

		defer func() {
			err := e.redigo.Locker().Release()
			if err != nil {
				fmt.Println(err)
			}
		}()
	}

	var tid float64
	{
		tid = float64(time.Now().UTC().UnixNano())
	}

	{
		tsk.SetBackoff(0)
		tsk.SetExpire(0)
		tsk.SetID(tid)
		tsk.SetOwner("")
		tsk.SetVersion(1)
	}

	{
		k := key.Task
		v := task.ToString(tsk)
		s := tid

		err = e.redigo.Sorted().Create().Element(k, v, s)
		if err != nil {
			return tracer.Mask(err)
		}
	}

	return nil
}
