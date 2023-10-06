package ticker

import "time"

func (t *Ticker) onUnP1(uni string) time.Time {
	if uni == "minute" {
		var tp1 time.Time
		{
			tp1 = t.onUnM1(uni).Add(time.Minute)
		}

		if tp1.Equal(t.tim) {
			return tp1.Add(time.Minute)
		}

		return tp1
	}

	if uni == "hour" {
		var tp1 time.Time
		{
			tp1 = t.onUnM1(uni).Add(time.Hour)
		}

		if tp1.Equal(t.tim) {
			return tp1.Add(time.Hour)
		}

		return tp1
	}

	if uni == "day" {
		var tp1 time.Time
		{
			tp1 = t.onUnM1(uni).AddDate(0, 0, 1)
		}

		if tp1.Equal(t.tim) {
			return tp1.AddDate(0, 0, 1)
		}

		return tp1
	}

	if uni == "week" {
		var tp1 time.Time
		{
			tp1 = t.onUnM1(uni).AddDate(0, 0, 7)
		}

		if tp1.Equal(t.tim) {
			return tp1.AddDate(0, 0, 7)
		}

		return tp1
	}

	if uni == "month" {
		var tp1 time.Time
		{
			tp1 = t.onUnM1(uni).AddDate(0, 1, 0)
		}

		if tp1.Equal(t.tim) {
			return tp1.AddDate(0, 1, 0)
		}

		return tp1
	}

	return time.Time{}
}
