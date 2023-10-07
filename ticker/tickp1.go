package ticker

import "time"

func (t *Ticker) TickP1() time.Time {
	if t.qnt == 1 {
		return t.tickP1(t.qnt, t.uni)
	}

	if t.qnt >= 2 {
		return t.tickP1(t.qnt, t.uni[:len(t.uni)-1])
	}

	return time.Time{}
}

func (t *Ticker) tickP1(qnt int, uni string) time.Time {
	if uni == "minute" {
		var tp1 time.Time
		{
			tp1 = t.tickM1(qnt, uni).Add(time.Duration(qnt) * time.Minute)
		}

		if tp1.Equal(t.tim) {
			return tp1.Add(time.Duration(qnt) * time.Minute)
		}

		return tp1
	}

	if uni == "hour" {
		var tp1 time.Time
		{
			tp1 = t.tickM1(qnt, uni).Add(time.Duration(qnt) * time.Hour)
		}

		if tp1.Equal(t.tim) {
			return tp1.Add(time.Duration(qnt) * time.Hour)
		}

		return tp1
	}

	if uni == "day" {
		var tp1 time.Time
		{
			tp1 = t.tickM1(qnt, uni).AddDate(0, 0, qnt)
		}

		if tp1.Equal(t.tim) {
			return tp1.AddDate(0, 0, qnt)
		}

		return tp1
	}

	if uni == "week" {
		var tp1 time.Time
		{
			tp1 = t.tickM1(qnt, uni).AddDate(0, 0, qnt*7)
		}

		if tp1.Equal(t.tim) {
			return tp1.AddDate(0, 0, qnt*7)
		}

		return tp1
	}

	if uni == "month" {
		var tp1 time.Time
		{
			tp1 = t.tickM1(qnt, uni).AddDate(0, qnt, 0)
		}

		if tp1.Equal(t.tim) {
			return tp1.AddDate(0, qnt, 0)
		}

		return tp1
	}

	return time.Time{}
}
