package engine

import (
	"fmt"

	"github.com/xh3b4sd/tracer"

	"github.com/xh3b4sd/rescue/pkg/metadata"
	"github.com/xh3b4sd/rescue/pkg/task"
)

func (e *Engine) exists(tsk *task.Task) (bool, error) {
	var err error

	// Checking for existing tasks implies certain read operations on the task
	// queue. For consistency reasons we need to ensure that only one process at
	// a time can read information from the queue. Otherwise worker behaviour
	// would be inconsistent and the integrity of the queue could not be
	// guaranteed.
	{
		err := e.redigo.Locker().Acquire()
		if err != nil {
			return false, tracer.Mask(err)
		}

		defer func() {
			err := e.redigo.Locker().Release()
			if err != nil {
				fmt.Println(err)
			}
		}()
	}

	var tks []*task.Task
	{
		tks, err = e.searchAll()
		if err != nil {
			return false, tracer.Mask(err)
		}

		e.metric.Task.Queued.Set(float64(len(tks)))
	}

	{
		for _, t := range tks {
			// When checking for metadata the task fetched from the queue must
			// be given first since it contains all the metadata of the task
			// itself. The task given to Engine.Exists contains only the
			// relevant subset of metadata we want to match against.
			if metadata.Contains(t.Obj.Metadata, tsk.Obj.Metadata) {
				return true, nil
			}
		}
	}

	return false, nil
}
