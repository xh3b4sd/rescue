package ticker

import "time"

func (t *Ticker) Duration() time.Duration {
	if t.qnt == 1 {
		return t.duration(t.qnt, t.uni)
	}

	if t.qnt >= 2 {
		return t.duration(t.qnt, t.uni[:len(t.uni)-1])
	}

	return 0
}

func (t *Ticker) duration(qnt int, uni string) time.Duration {
	if uni == "minute" {
		return time.Duration(qnt) * time.Minute
	}

	if uni == "hour" {
		return time.Duration(qnt) * time.Hour
	}

	if uni == "day" {
		return time.Duration(qnt) * time.Hour * 24
	}

	if uni == "week" {
		return time.Duration(qnt) * time.Hour * 24 * 7
	}

	if uni == "month" {
		return time.Duration(qnt) * time.Hour * 24 * 30
	}

	return 0
}
