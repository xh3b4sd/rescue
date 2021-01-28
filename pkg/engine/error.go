package engine

import (
	"errors"

	"github.com/xh3b4sd/tracer"
)

var searchFailedError = &tracer.Error{
	Kind: "searchFailedError",
	Desc: "When searching tasks, using a single score, the number of resulting tasks must be 1.",
}

func IsExecutionFailed(err error) bool {
	return errors.Is(err, searchFailedError)
}

var invalidConfigError = &tracer.Error{
	Kind: "invalidConfigError",
}

func IsInvalidConfig(err error) bool {
	return errors.Is(err, invalidConfigError)
}

var invalidTaskError = &tracer.Error{
	Kind: "invalidTaskError",
}

func IsInvalidTask(err error) bool {
	return errors.Is(err, invalidTaskError)
}

var noTaskError = &tracer.Error{
	Kind: "noTaskError",
	Desc: "When searching tasks, there might not be any task available. This is expected in case there is simply no work to be done. Searching for tasks should be repeated after a short period of time.",
}

func IsNoTask(err error) bool {
	return errors.Is(err, noTaskError)
}

var taskOutdatedError = &tracer.Error{
	Kind: "taskOutdatedError",
	Desc: "When deleting tasks, the tasks provided by workers must be consistent with the current state of the queue. One possibility is that the task provided got accidentially changed by the current worker. Another possibility is that the task expired meanwhile, which then resulted in the task being reset and picked up by another worker eventually.",
}

func IsTaskOutdatedled(err error) bool {
	return errors.Is(err, taskOutdatedError)
}
