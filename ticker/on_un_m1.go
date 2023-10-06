package ticker

import "time"

func (t *Ticker) onUnM1(uni string) time.Time {
	if uni == "minute" {
		var tm1 time.Time
		{
			tm1 = t.truMin()
		}

		if tm1.Equal(t.tim) {
			return tm1.Add(-time.Minute)
		}

		return tm1
	}

	if uni == "hour" {
		var tm1 time.Time
		{
			tm1 = t.truHou()
		}

		if tm1.Equal(t.tim) {
			return tm1.Add(-time.Hour)
		}

		return tm1
	}

	if uni == "day" {
		var tm1 time.Time
		{
			tm1 = t.truDay()
		}

		if tm1.Equal(t.tim) {
			return tm1.AddDate(0, 0, -1)
		}

		return tm1
	}

	if uni == "week" {
		var tim time.Time
		{
			tim = t.truDay()
		}

		var wee time.Weekday
		if tim.Weekday() == time.Sunday {
			wee = 7
		} else {
			wee = tim.Weekday()
		}

		var rem int
		{
			rem = -int(wee - time.Monday)
		}

		var tm1 time.Time
		{
			tm1 = tim.AddDate(0, 0, rem)
		}

		if tm1.Equal(t.tim) {
			return tm1.AddDate(0, 0, -7)
		}

		return tm1
	}

	if uni == "month" {
		var tm1 time.Time
		{
			tm1 = t.truMon()
		}

		if tm1.Equal(t.tim) {
			return tm1.AddDate(0, -1, 0)
		}

		return tm1
	}

	return time.Time{}
}
