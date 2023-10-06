package ticker

import (
	"strconv"
	"strings"
	"time"
)

type Ticker struct {
	fld []string
	fmt string
	tim time.Time
}

func New(fmt string, tim ...time.Time) *Ticker {
	var now time.Time
	if len(tim) == 1 && !tim[0].IsZero() {
		now = tim[0].UTC()
	} else {
		now = time.Now().UTC()
	}

	var tic *Ticker
	{
		tic = &Ticker{
			fld: strings.Fields(fmt),
			fmt: fmt,
			tim: now,
		}
	}

	return tic
}

func (t *Ticker) TickM1() time.Time {
	if len(t.fld) == 0 {
		return time.Time{}
	}

	if len(t.fld) == 1 {
		return t.onUnM1(t.fld[0])
	}

	if len(t.fld) == 2 {
		return t.qnUnM1(musNum(t.fld[0]), t.fld[1])
	}

	return time.Time{}
}

func (t *Ticker) TickP1() time.Time {
	if len(t.fld) == 0 {
		return time.Time{}
	}

	if len(t.fld) == 1 {
		return t.onUnP1(t.fld[0])
	}

	if len(t.fld) == 2 {
		return t.qnUnP1(musNum(t.fld[0]), t.fld[1])
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
