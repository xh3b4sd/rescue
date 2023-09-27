package validate

import (
	"strings"

	"github.com/xh3b4sd/tracer"

	"github.com/xh3b4sd/rescue/pkg/task"
)

func Empty(tas *task.Task) error {
	if tas == nil {
		return tracer.Maskf(invalidTaskError, "task must not be nil")
	}

	if len(tas.Obj.Metadata) == 0 {
		return tracer.Maskf(invalidTaskError, "metadata must not be empty")
	}

	return nil
}

func Label(tas *task.Task) error {
	for k := range tas.Obj.Metadata {
		if strings.HasPrefix(k, "task.rescue.io") {
			return tracer.Maskf(invalidTaskError, "metadata must not contain reserved scheme task.rescue.io")
		}
	}

	return nil
}
