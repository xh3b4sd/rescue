package timer

import "time"

type Timer struct {
	fac func() time.Time
}

func New() *Timer {
	return &Timer{
		fac: func() time.Time {
			return time.Now().UTC()
		},
	}
}

func (t *Timer) Create() time.Time {
	return t.fac()
}

func (t *Timer) Expire() time.Time {
	return t.fac()
}

func (t *Timer) Extend() time.Time {
	return t.fac()
}

func (t *Timer) Search() time.Time {
	return t.fac()
}

func (t *Timer) Ticker() time.Time {
	return t.fac()
}

func (t *Timer) Setter(fac func() time.Time) {
	t.fac = fac
}
