package task

import "time"

// TODO unit test
// NewTimer returns a new timer configured with the duration between now and
// this task's expiry. Optionally, a single time instance can be given to
// overwrite the internal UTC based "now". If this task is already expired, then
// the result of NewTimer will fire immediately.
//
// Since timers allocate resources at runtime they need to be managed carefully.
// If a timer is not used anymore before it fired, then it should be drained
// manually in order for the garbage collector to remove the timer and its
// internally used channel.
//
//	https://pkg.go.dev/time#Timer.Stop
func NewTimer(tas *Task, tim ...time.Time) *time.Timer {
	var now time.Time
	if len(tim) == 1 && !tim[0].IsZero() {
		now = tim[0]
	} else {
		now = time.Now().UTC()
	}

	var exp time.Time
	{
		exp = tas.Get().Expiry()
	}

	var dur time.Duration
	if exp.After(now) {
		dur = time.Duration(exp.Sub(now))
	}

	return time.NewTimer(dur)
}
