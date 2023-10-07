package ticker

import (
	"strconv"
	"strings"
	"time"
)

type Ticker struct {
	qnt int
	tim time.Time
	uni string
}

func New(fmt string, tim ...time.Time) *Ticker {
	var fld []string
	{
		fld = strings.Fields(fmt)
	}

	var now time.Time
	if len(tim) == 1 && !tim[0].IsZero() {
		now = tim[0].UTC()
	} else {
		now = time.Now().UTC()
	}

	var qnt int
	var uni string

	switch len(fld) {
	case 1:
		qnt = 1
		uni = fld[0]
	case 2:
		qnt = musNum(fld[0])
		uni = fld[1]
	}

	var tic *Ticker
	{
		tic = &Ticker{
			qnt: qnt,
			tim: now,
			uni: uni,
		}
	}

	return tic
}

func (t *Ticker) TickM1() time.Time {
	if t.qnt == 1 {
		return t.tickM1(t.qnt, t.uni)
	}

	if t.qnt >= 2 {
		return t.tickM1(t.qnt, t.uni[:len(t.uni)-1])
	}

	return time.Time{}
}

func (t *Ticker) TickP1() time.Time {
	if t.qnt == 1 {
		return t.tickP1(t.qnt, t.uni)
	}

	if t.qnt >= 2 {
		return t.tickP1(t.qnt, t.uni[:len(t.uni)-1])
	}

	return time.Time{}
}

func musNum(str string) int {
	num, err := strconv.Atoi(str)
	if err != nil {
		return 0
	}

	return num
}
