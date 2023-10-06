package ticker

import (
	"math"
	"time"
)

func (t *Ticker) qnUnM1(qnt int, uni string) time.Time {
	if qnt <= 1 {
		return time.Time{}
	}

	if uni == "minutes" {
		var cur float64
		{
			cur = float64(t.tim.Minute())
		}

		var min time.Duration
		{
			mod, _ := math.Modf(cur / float64(qnt))
			min = time.Duration(mod * float64(qnt))
		}

		var tm1 time.Time
		{
			tm1 = t.truHou().Add(min * time.Minute)
		}

		if tm1.Equal(t.tim) {
			return tm1.Add(time.Duration(-qnt) * time.Minute)
		}

		return tm1
	}

	if uni == "hours" {
		var cur float64
		{
			cur = float64(t.tim.Hour())
		}

		var hou time.Duration
		{
			mod, _ := math.Modf(cur / float64(qnt))
			hou = time.Duration(mod * float64(qnt))
		}

		var tm1 time.Time
		{
			tm1 = t.truDay().Add(hou * time.Hour)
		}

		if tm1.Equal(t.tim) {
			return tm1.Add(time.Duration(-qnt) * time.Hour)
		}

		return tm1
	}

	if uni == "days" {
		var cur int
		{
			cur = t.tim.Day()
		}

		var day int
		{
			mod, _ := math.Modf(float64(cur) / float64(qnt))
			day = int(mod * float64(qnt))
		}

		if day == 0 {
			day = 1
		}

		var tm1 time.Time
		{
			tm1 = t.truMon().AddDate(0, 0, day-1)
		}

		if tm1.Equal(t.tim) {
			return tm1.AddDate(0, 0, -qnt)
		}

		return tm1
	}

	if uni == "weeks" {
		var cur int
		{
			_, cur = t.tim.ISOWeek()
		}

		var wee int
		{
			mod, _ := math.Modf(float64(cur) / float64(qnt))
			wee = int(mod * float64(qnt))
		}

		var tm1 time.Time
		{
			tm1 = t.truWee().AddDate(0, 0, wee*7)
		}

		if tm1.Equal(t.tim) || cur == wee {
			return tm1.AddDate(0, 0, -qnt*7)
		}

		return tm1
	}

	if uni == "months" {
		var cur int
		{
			cur = int(t.tim.Month())
		}

		var mon int
		{
			mod, _ := math.Modf(float64(cur) / float64(qnt))
			mon = int(mod * float64(qnt))
		}

		var tm1 time.Time
		{
			tm1 = t.truYea().AddDate(0, mon, 0)
		}

		if tm1.Equal(t.tim) || cur == mon {
			return tm1.AddDate(0, -qnt, 0)
		}

		return tm1
	}

	return time.Time{}
}
