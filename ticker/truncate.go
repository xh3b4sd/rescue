package ticker

import "time"

func (t *Ticker) truMin() time.Time {
	return time.Date(t.tim.Year(), t.tim.Month(), t.tim.Day(), t.tim.Hour(), t.tim.Minute(), 0, 0, time.UTC)
}

func (t *Ticker) truHou() time.Time {
	return time.Date(t.tim.Year(), t.tim.Month(), t.tim.Day(), t.tim.Hour(), 0, 0, 0, time.UTC)
}

func (t *Ticker) truDay() time.Time {
	return time.Date(t.tim.Year(), t.tim.Month(), t.tim.Day(), 0, 0, 0, 0, time.UTC)
}

// truWee is to truncate the underlying time instance to the start of the first
// week of the year. That is, the first Monday of any given year.
func (t *Ticker) truWee() time.Time {
	wee := time.Date(t.tim.Year(), time.January, 1, 0, 0, 0, 0, time.UTC).Weekday()
	d := (8-int(wee))%7 + 1
	return time.Date(t.tim.Year(), time.January, d, 0, 0, 0, 0, time.UTC)
}

func (t *Ticker) truMon() time.Time {
	return time.Date(t.tim.Year(), t.tim.Month(), 1, 0, 0, 0, 0, time.UTC)
}

func (t *Ticker) truYea() time.Time {
	return time.Date(t.tim.Year(), time.January, 1, 0, 0, 0, 0, time.UTC)
}
