package task

import "time"

type Getter interface {
	Bypass() bool
	Cycles() int64
	Expiry() time.Time
	Object() int64
	Worker() string
}

type Interface interface {
	// All returns a task containing all metadata matching the given metadata
	// keys. If any of the given keys does not match, nil is returned. That means
	// that the returned task will be nil, unless the complete list of the given
	// metadata keys matches against this task.
	All(...string) *Task

	// Any returns a task containing any metadata matching the given metadata
	// keys. If any of the given keys does not match, it is simply ignored. That
	// means that the returned task might be nil if not a single metadata key
	// matches this task's metadata. If some of the given keys match, a task with
	// the matching metadata is returned.
	Any(...string) *Task

	// Get provides a getter interface for reading task metadata.
	Get() Getter

	// Has expresses whether this task contains all of the given metadata.
	Has(map[string]string) bool

	// Prg provides a purger interface for purging task metadata.
	Prg() Setter

	// Set provides a setter interface for writing task metadata.
	Set() Setter
}

type Purger interface {
	Bypass()
	Cycles()
	Expiry()
	Object()
	Worker()
}

type Setter interface {
	Bypass(bool)
	Cycles(int64)
	Expiry(time.Time)
	Object(int64)
	Worker(string)
}
