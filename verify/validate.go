package verify

import (
	"strings"

	"github.com/xh3b4sd/rescue/task"
	"github.com/xh3b4sd/tracer"
)

func Empty(tas *task.Task) error {
	if tas == nil {
		return tracer.Maskf(invalidTaskError, "task must not be nil")
	}

	if len(tas.Meta) == 0 {
		return tracer.Maskf(invalidTaskError, "metadata must not be empty")
	}

	return nil
}

func Label(tas *task.Task) error {
	for k := range tas.Meta {
		if strings.HasPrefix(k, "task.rescue.io") {
			return tracer.Maskf(invalidTaskError, "metadata must not contain reserved scheme task.rescue.io")
		}
	}

	return nil
}
