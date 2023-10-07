package ticker

import (
	"math"
	"time"
)

func (t *Ticker) TickM1() time.Time {
	if t.qnt == 1 {
		return t.tickM1(t.qnt, t.uni)
	}

	if t.qnt >= 2 {
		return t.tickM1(t.qnt, t.uni[:len(t.uni)-1])
	}

	return time.Time{}
}

func (t *Ticker) tickM1(qnt int, uni string) time.Time {
	if uni == "minute" {
		var min int
		{
			min = modQnt(t.tim.Minute(), qnt)
		}

		var tm1 time.Time
		{
			tm1 = t.truHou().Add(time.Duration(min) * time.Minute)
		}

		if tm1.Equal(t.tim) {
			return tm1.Add(time.Duration(-qnt) * time.Minute)
		}

		return tm1
	}

	if uni == "hour" {
		var hou int
		{
			hou = modQnt(t.tim.Hour(), qnt)
		}

		var tm1 time.Time
		{
			tm1 = t.truDay().Add(time.Duration(hou) * time.Hour)
		}

		if tm1.Equal(t.tim) {
			return tm1.Add(time.Duration(-qnt) * time.Hour)
		}

		return tm1
	}

	if uni == "day" {
		var day int
		{
			day = modQnt(t.uniDay()-1, qnt)
		}

		var tm1 time.Time
		{
			tm1 = t.truUni().AddDate(0, 0, day)
		}

		if tm1.Equal(t.tim) {
			return tm1.AddDate(0, 0, -qnt)
		}

		return tm1
	}

	if uni == "week" {
		var day int
		{
			day = modQnt((t.uniDay()-4)-1, qnt*7)
		}

		// The basis for weekly schedules is the first unix time Monday, the 5th of
		// January 1970. Since we use unix seconds as fixed starting point for days
		// and weeks, and since the 1st of January 1970 was a Thursday, we need to
		// add 4 days to the truncated unix time in order to arrive at our starting
		// point of the first unix time Monday.
		var tm1 time.Time
		{
			tm1 = t.truUni().AddDate(0, 0, day+4)
		}

		if tm1.Equal(t.tim) {
			return tm1.AddDate(0, 0, -qnt*7)
		}

		return tm1
	}

	if uni == "month" {
		var mon int
		{
			mon = modQnt(int(t.tim.Month()-1), qnt)
		}

		var tm1 time.Time
		{
			tm1 = t.truYea().AddDate(0, mon, 0)
		}

		if tm1.Equal(t.tim) {
			return tm1.AddDate(0, -qnt, 0)
		}

		return tm1
	}

	return time.Time{}
}

func modQnt(x int, y int) int {
	m, _ := math.Modf(float64(x) / float64(y))
	return int(m) * y
}
