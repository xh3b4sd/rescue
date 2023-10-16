package engine

import (
	"errors"

	"github.com/xh3b4sd/tracer"
)

var labelReservedError = &tracer.Error{
	Kind: "labelReservedError",
}

func IsLabelReserved(err error) bool {
	return errors.Is(err, labelReservedError)
}

var labelValueError = &tracer.Error{
	Kind: "labelValueError",
}

func IsLabelValue(err error) bool {
	return errors.Is(err, labelValueError)
}

var taskCoreError = &tracer.Error{
	Kind: "taskCoreError",
}

func IsTaskCore(err error) bool {
	return errors.Is(err, taskCoreError)
}

var taskCronError = &tracer.Error{
	Kind: "taskCronError",
}

func IsTaskCron(err error) bool {
	return errors.Is(err, taskCronError)
}

var taskEmptyError = &tracer.Error{
	Kind: "taskEmptyError",
}

func IsTaskEmpty(err error) bool {
	return errors.Is(err, taskEmptyError)
}

var taskMetaEmptyError = &tracer.Error{
	Kind: "taskMetaEmptyError",
}

func IsTaskMetaEmpty(err error) bool {
	return errors.Is(err, taskMetaEmptyError)
}

var taskNotFoundError = &tracer.Error{
	Kind: "taskNotFoundError",
	Desc: "When searching tasks, there might not be any task available. This is expected in case there is simply no work to be done. Searching for tasks should be repeated after a short period of time.",
}

func IsTaskNotFound(err error) bool {
	return errors.Is(err, taskNotFoundError)
}

var taskNotRevokedError = &tracer.Error{
	Kind: "taskNotRevokedError",
	Desc: "When expiring tasks, there might be tasks that should be rescheduled. The deviations we find in the system must all be resolved in order to regain the desired balance of sharing the work load. This error indicates that the system is fundamentally broken.",
}

func IsTaskNotRevoked(err error) bool {
	return errors.Is(err, taskNotFoundError)
}

var taskOutdatedError = &tracer.Error{
	Kind: "taskOutdatedError",
	Desc: "When deleting tasks, the tasks provided by workers must be consistent with the current state of the queue. One possibility is that the task provided got accidentially changed by the current worker. Another possibility is that the task expired meanwhile, which then resulted in the task being reset and picked up by another worker eventually.",
}

func IsTaskOutdated(err error) bool {
	return errors.Is(err, taskOutdatedError)
}
