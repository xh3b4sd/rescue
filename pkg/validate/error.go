package validate

import (
	"errors"

	"github.com/xh3b4sd/tracer"
)

var invalidTaskError = &tracer.Error{
	Kind: "invalidTaskError",
}

func IsInvalidTask(err error) bool {
	return errors.Is(err, invalidTaskError)
}
