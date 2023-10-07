package ticker

import "time"

func (t *Ticker) uniDay() int {
	return int(t.tim.Unix()/60/60/24) + 1
}

func (t *Ticker) truUni() time.Time {
	return time.Date(1970, time.January, 1, 0, 0, 0, 0, time.UTC)
}

func (t *Ticker) truHou() time.Time {
	return time.Date(t.tim.Year(), t.tim.Month(), t.tim.Day(), t.tim.Hour(), 0, 0, 0, time.UTC)
}

func (t *Ticker) truDay() time.Time {
	return time.Date(t.tim.Year(), t.tim.Month(), t.tim.Day(), 0, 0, 0, 0, time.UTC)
}

func (t *Ticker) truYea() time.Time {
	return time.Date(t.tim.Year(), time.January, 1, 0, 0, 0, 0, time.UTC)
}
