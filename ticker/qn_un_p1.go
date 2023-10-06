package ticker

import "time"

func (t *Ticker) qnUnP1(qnt int, uni string) time.Time {
	if qnt <= 1 {
		return time.Time{}
	}

	if uni == "minutes" {
		var tp1 time.Time
		{
			tp1 = t.qnUnM1(qnt, uni).Add(time.Duration(qnt) * time.Minute)
		}

		if tp1.Equal(t.tim) {
			return tp1.Add(time.Duration(qnt) * time.Minute)
		}

		return tp1
	}

	if uni == "hours" {
		var tp1 time.Time
		{
			tp1 = t.qnUnM1(qnt, uni).Add(time.Duration(qnt) * time.Hour)
		}

		if tp1.Equal(t.tim) {
			return tp1.Add(time.Duration(qnt) * time.Hour)
		}

		return tp1
	}

	if uni == "days" {
		var tp1 time.Time
		{
			tp1 = t.qnUnM1(qnt, uni).AddDate(0, 0, qnt)
		}

		if tp1.Equal(t.tim) {
			return tp1.AddDate(0, 0, qnt)
		}

		return tp1
	}

	if uni == "weeks" {
		var tp1 time.Time
		{
			tp1 = t.qnUnM1(qnt, uni).AddDate(0, 0, qnt*7)
		}

		if tp1.Equal(t.tim) {
			return tp1.AddDate(0, 0, qnt*7)
		}

		return tp1
	}

	if uni == "months" {
		var tp1 time.Time
		{
			tp1 = t.qnUnM1(qnt, uni).AddDate(0, qnt, 0)
		}

		if tp1.Equal(t.tim) {
			return tp1.AddDate(0, qnt, 0)
		}

		return tp1
	}

	return time.Time{}
}
