package ticker

import (
	"fmt"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
)

// TODO ensure multiple months are only valid if they are defined as a multiple
// of 12, so that they repeat monotonically equal in every year
func Test_Ticker_Invalid(t *testing.T) {
	testCases := []struct {
		fmt string
	}{
		// Case 000, ensures that empty strings produce a zero time.
		{
			fmt: "",
		},
		// Case 001, ensures that whitespace strings produce a zero time.
		{
			fmt: " ",
		},
		// Case 002, ensures that random strings produce a zero time.
		{
			fmt: "foo",
		},
		// Case 003, ensures that the unit second produces a zero time.
		{
			fmt: "second",
		},
		// Case 004, ensures that seconds without quantity produce a zero time.
		{
			fmt: "seconds",
		},
		// Case 005, ensures that seconds with quantity produce a zero time.
		{
			fmt: "20 seconds",
		},
		// Case 006, ensures that years without quantity produce a zero time.
		{
			fmt: "years",
		},
		// Case 007, ensures that years with quantity produce a zero time.
		{
			fmt: "8 years",
		},
		// Case 008, ensures that the single hour unit with quantity produces a zero
		// time.
		{
			fmt: "3 hour",
		},
		// Case 009, ensures that hours without quantity produce a zero time.
		{
			fmt: "hours",
		},
		// Case 010, ensures that non-integer quantities produce a zero time.
		{
			fmt: "eight hours",
		},
		// Case 011, ensures that days without quantity produce a zero time.
		{
			fmt: "days",
		},
		// Case 012, ensures that the single day unit with quantity produces a zero
		// time.
		{
			fmt: "5 day",
		},
		// Case 013, ensures that quantity and unit flipped produces a zero time.
		{
			fmt: "weeks 5",
		},
		// Case 014, ensures that quantity and quantity produces a zero time.
		{
			fmt: "18 4",
		},
	}

	for i, tc := range testCases {
		t.Run(fmt.Sprintf("%03d", i), func(t *testing.T) {
			var tic *Ticker
			{
				tic = New(tc.fmt)
			}

			var tm1 time.Time
			{
				tm1 = tic.TickM1()
			}

			if !tm1.IsZero() {
				t.Fatalf("%s tick-1\n\n%s\n", tc.fmt, cmp.Diff(time.Time{}, tm1))
			}

			var tp1 time.Time
			{
				tp1 = tic.TickP1()
			}

			if !tp1.IsZero() {
				t.Fatalf("%s tick+1\n\n%s\n", tc.fmt, cmp.Diff(time.Time{}, tp1))
			}
		})
	}
}

func Test_Ticker_Only_Unit(t *testing.T) {
	testCases := []struct {
		fmt string
		tm1 time.Time
		now time.Time
		tp1 time.Time
	}{
		// Case 000
		{
			fmt: "minute",
			tm1: musTim("2023-09-28T14:22:00.000000Z"),
			now: musTim("2023-09-28T14:23:00.000000Z"),
			tp1: musTim("2023-09-28T14:24:00.000000Z"),
		},
		// Case 001
		{
			fmt: "minute",
			tm1: musTim("2023-09-28T14:23:00.000000Z"),
			now: musTim("2023-09-28T14:23:24.161982Z"),
			tp1: musTim("2023-09-28T14:24:00.000000Z"),
		},
		// Case 002
		{
			fmt: "minute",
			tm1: musTim("2023-09-28T23:59:00.000000Z"),
			now: musTim("2023-09-28T23:59:24.161982Z"),
			tp1: musTim("2023-09-29T00:00:00.000000Z"),
		},
		// Case 003
		{
			fmt: "hour",
			tm1: musTim("2023-09-28T13:00:00.000000Z"),
			now: musTim("2023-09-28T14:00:00.000000Z"),
			tp1: musTim("2023-09-28T15:00:00.000000Z"),
		},
		// Case 004
		{
			fmt: "hour",
			tm1: musTim("2023-09-28T14:00:00.000000Z"),
			now: musTim("2023-09-28T14:23:24.161982Z"),
			tp1: musTim("2023-09-28T15:00:00.000000Z"),
		},
		// Case 005
		{
			fmt: "hour",
			tm1: musTim("2023-09-28T23:00:00.000000Z"),
			now: musTim("2023-09-28T23:59:24.161982Z"),
			tp1: musTim("2023-09-29T00:00:00.000000Z"),
		},
		// Case 006
		{
			fmt: "day",
			tm1: musTim("2023-09-27T00:00:00.000000Z"),
			now: musTim("2023-09-28T00:00:00.000000Z"),
			tp1: musTim("2023-09-29T00:00:00.000000Z"),
		},
		// Case 007
		{
			fmt: "day",
			tm1: musTim("2023-09-30T00:00:00.000000Z"),
			now: musTim("2023-10-01T00:00:00.000000Z"),
			tp1: musTim("2023-10-02T00:00:00.000000Z"),
		},
		// Case 008
		{
			fmt: "day",
			tm1: musTim("2023-09-28T00:00:00.000000Z"),
			now: musTim("2023-09-28T14:23:24.161982Z"),
			tp1: musTim("2023-09-29T00:00:00.000000Z"),
		},
		// Case 009
		{
			fmt: "day",
			tm1: musTim("2023-09-30T00:00:00.000000Z"),
			now: musTim("2023-09-30T14:23:24.161982Z"),
			tp1: musTim("2023-10-01T00:00:00.000000Z"),
		},
		// Case 010
		{
			fmt: "week",
			tm1: musTim("2023-09-18T00:00:00.000000Z"), // Monday
			now: musTim("2023-09-25T00:00:00.000000Z"), // Monday
			tp1: musTim("2023-10-02T00:00:00.000000Z"), // Monday
		},
		// Case 011
		{
			fmt: "week",
			tm1: musTim("2023-09-25T00:00:00.000000Z"), // Monday
			now: musTim("2023-09-28T14:23:24.161982Z"), // Thursday
			tp1: musTim("2023-10-02T00:00:00.000000Z"), // Monday
		},
		// Case 012
		{
			fmt: "week",
			tm1: musTim("2023-09-25T00:00:00.000000Z"), // Monday
			now: musTim("2023-10-01T14:23:24.161982Z"), // Sunday
			tp1: musTim("2023-10-02T00:00:00.000000Z"), // Monday
		},
		// Case 013
		{
			fmt: "month",
			tm1: musTim("2023-08-01T00:00:00.000000Z"),
			now: musTim("2023-09-01T00:00:00.000000Z"),
			tp1: musTim("2023-10-01T00:00:00.000000Z"),
		},
		// Case 014
		{
			fmt: "month",
			tm1: musTim("2023-12-01T00:00:00.000000Z"),
			now: musTim("2024-01-01T00:00:00.000000Z"),
			tp1: musTim("2024-02-01T00:00:00.000000Z"),
		},
		// Case 015
		{
			fmt: "month",
			tm1: musTim("2023-09-01T00:00:00.000000Z"),
			now: musTim("2023-09-28T14:23:24.161982Z"),
			tp1: musTim("2023-10-01T00:00:00.000000Z"),
		},
		// Case 016
		{
			fmt: "month",
			tm1: musTim("2023-12-01T00:00:00.000000Z"),
			now: musTim("2023-12-28T14:23:24.161982Z"),
			tp1: musTim("2024-01-01T00:00:00.000000Z"),
		},
	}

	for i, tc := range testCases {
		t.Run(fmt.Sprintf("%03d", i), func(t *testing.T) {
			var tic *Ticker
			{
				tic = New(tc.fmt, tc.now)
			}

			var tm1 time.Time
			{
				tm1 = tic.TickM1()
			}

			if !tm1.Equal(tc.tm1) {
				t.Fatalf("%s tick-1\n\n%s\n", tc.fmt, cmp.Diff(tc.tm1, tm1))
			}

			var tp1 time.Time
			{
				tp1 = tic.TickP1()
			}

			if !tp1.Equal(tc.tp1) {
				t.Fatalf("%s tick+1\n\n%s\n", tc.fmt, cmp.Diff(tc.tp1, tp1))
			}
		})
	}
}

func Test_Ticker_Quantity_And_Unit(t *testing.T) {
	testCases := []struct {
		fmt string
		tm1 time.Time
		now time.Time
		tp1 time.Time
	}{
		// Case 000
		{
			fmt: "5 minutes",
			tm1: musTim("2023-09-28T14:15:00.000000Z"),
			now: musTim("2023-09-28T14:20:00.000000Z"),
			tp1: musTim("2023-09-28T14:25:00.000000Z"),
		},
		// Case 001
		{
			fmt: "5 minutes",
			tm1: musTim("2023-09-27T23:55:00.000000Z"),
			now: musTim("2023-09-28T00:00:00.000000Z"),
			tp1: musTim("2023-09-28T00:05:00.000000Z"),
		},
		// Case 002
		{
			fmt: "5 minutes",
			tm1: musTim("2023-09-28T14:00:00.000000Z"),
			now: musTim("2023-09-28T14:03:24.161982Z"),
			tp1: musTim("2023-09-28T14:05:00.000000Z"),
		},
		// Case 003
		{
			fmt: "5 minutes",
			tm1: musTim("2023-09-28T14:20:00.000000Z"),
			now: musTim("2023-09-28T14:23:24.161982Z"),
			tp1: musTim("2023-09-28T14:25:00.000000Z"),
		},
		// Case 004
		{
			fmt: "5 minutes",
			tm1: musTim("2023-09-28T23:55:00.000000Z"),
			now: musTim("2023-09-28T23:59:24.161982Z"),
			tp1: musTim("2023-09-29T00:00:00.000000Z"),
		},
		// Case 005
		{
			fmt: "5 minutes",
			tm1: musTim("2023-09-28T14:00:00.000000Z"),
			now: musTim("2023-09-28T14:00:24.161982Z"),
			tp1: musTim("2023-09-28T14:05:00.000000Z"),
		},
		// Case 006
		{
			fmt: "6 hours",
			tm1: musTim("2023-09-28T06:00:00.000000Z"),
			now: musTim("2023-09-28T12:00:00.000000Z"),
			tp1: musTim("2023-09-28T18:00:00.000000Z"),
		},
		// Case 007
		{
			fmt: "6 hours",
			tm1: musTim("2023-09-27T18:00:00.000000Z"),
			now: musTim("2023-09-28T00:00:00.000000Z"),
			tp1: musTim("2023-09-28T06:00:00.000000Z"),
		},
		// Case 008
		{
			fmt: "6 hours",
			tm1: musTim("2023-09-28T00:00:00.000000Z"),
			now: musTim("2023-09-28T04:23:24.161982Z"),
			tp1: musTim("2023-09-28T06:00:00.000000Z"),
		},
		// Case 009
		{
			fmt: "6 hours",
			tm1: musTim("2023-09-28T12:00:00.000000Z"),
			now: musTim("2023-09-28T14:23:24.161982Z"),
			tp1: musTim("2023-09-28T18:00:00.000000Z"),
		},
		// Case 010
		{
			fmt: "6 hours",
			tm1: musTim("2023-09-28T18:00:00.000000Z"),
			now: musTim("2023-09-28T23:59:24.161982Z"),
			tp1: musTim("2023-09-29T00:00:00.000000Z"),
		},
		// Case 011
		{
			fmt: "6 hours",
			tm1: musTim("2023-09-28T00:00:00.000000Z"),
			now: musTim("2023-09-28T00:23:24.161982Z"),
			tp1: musTim("2023-09-28T06:00:00.000000Z"),
		},
		// Case 012
		{
			fmt: "3 days",
			tm1: musTim("2023-09-26T00:00:00.000000Z"),
			now: musTim("2023-09-27T00:00:00.000000Z"),
			tp1: musTim("2023-09-29T00:00:00.000000Z"),
		},
		// Case 013
		{
			fmt: "3 days",
			tm1: musTim("2023-09-26T00:00:00.000000Z"),
			now: musTim("2023-09-28T00:00:00.000000Z"),
			tp1: musTim("2023-09-29T00:00:00.000000Z"),
		},
		// Case 014
		{
			fmt: "3 days",
			tm1: musTim("2023-09-26T00:00:00.000000Z"),
			now: musTim("2023-09-29T00:00:00.000000Z"),
			tp1: musTim("2023-10-02T00:00:00.000000Z"),
		},
		// Case 015
		{
			fmt: "3 days",
			tm1: musTim("2023-09-29T00:00:00.000000Z"),
			now: musTim("2023-09-29T00:00:00.161982Z"),
			tp1: musTim("2023-10-02T00:00:00.000000Z"),
		},
		// Case 016
		{
			fmt: "3 days",
			tm1: musTim("2023-09-29T00:00:00.000000Z"),
			now: musTim("2023-09-29T14:23:24.161982Z"),
			tp1: musTim("2023-10-02T00:00:00.000000Z"),
		},
		// Case 017
		{
			fmt: "3 days",
			tm1: musTim("2023-09-29T00:00:00.000000Z"),
			now: musTim("2023-10-01T00:00:00.000000Z"),
			tp1: musTim("2023-10-02T00:00:00.000000Z"),
		},
		// Case 018
		{
			fmt: "3 days",
			tm1: musTim("2023-09-26T00:00:00.000000Z"),
			now: musTim("2023-09-28T14:23:24.161982Z"),
			tp1: musTim("2023-09-29T00:00:00.000000Z"),
		},
		// Case 019
		{
			fmt: "3 days",
			tm1: musTim("2023-09-26T00:00:00.000000Z"),
			now: musTim("2023-09-28T00:00:00.161982Z"),
			tp1: musTim("2023-09-29T00:00:00.000000Z"),
		},
		// Case 020
		{
			fmt: "3 days",
			tm1: musTim("2023-10-02T00:00:00.000000Z"),
			now: musTim("2023-10-02T23:59:24.161982Z"),
			tp1: musTim("2023-10-05T00:00:00.000000Z"),
		},
		// Case 021
		{
			fmt: "3 days",
			tm1: musTim("2023-09-29T00:00:00.000000Z"),
			now: musTim("2023-09-30T00:00:24.161982Z"),
			tp1: musTim("2023-10-02T00:00:00.000000Z"),
		},
		// Case 022
		{
			fmt: "3 days",
			tm1: musTim("2023-09-29T00:00:00.000000Z"),
			now: musTim("2023-10-01T23:59:24.161982Z"),
			tp1: musTim("2023-10-02T00:00:00.000000Z"),
		},
		// Case 023 shows how multiple day schedules properly carry over
		// continuously between the years. Using the first day of any given year or
		// the first ISO week as the basis of multi day schedules does not carry
		// over continuously between every year. Since we use the first unix time
		// day as basis, our calculations are correct regardless the year change.
		{
			fmt: "3 days",
			tm1: musTim("2022-12-30T00:00:00.000000Z"),
			now: musTim("2022-12-31T14:23:24.161982Z"),
			tp1: musTim("2023-01-02T00:00:00.000000Z"),
		},
		// Case 024 shows how multiple day schedules properly carry over
		// continuously between the years. Using the first day of any given year or
		// the first ISO week as the basis of multi day schedules does not carry
		// over continuously between every year. Since we use the first unix time
		// day as basis, our calculations are correct regardless the year change.
		{
			fmt: "3 days",
			tm1: musTim("2022-12-30T00:00:00.000000Z"),
			now: musTim("2023-01-01T00:00:00.000000Z"),
			tp1: musTim("2023-01-02T00:00:00.000000Z"),
		},
		// Case 025 shows how multiple day schedules properly carry over
		// continuously between the years. Using the first day of any given year or
		// the first ISO week as the basis of multi day schedules does not carry
		// over continuously between every year. Since we use the first unix time
		// day as basis, our calculations are correct regardless the year change.
		{
			fmt: "3 days",
			tm1: musTim("2022-12-30T00:00:00.000000Z"),
			now: musTim("2023-01-01T14:23:24.161982Z"),
			tp1: musTim("2023-01-02T00:00:00.000000Z"),
		},
		// Case 026
		{
			fmt: "7 days",
			tm1: musTim("2023-09-28T00:00:00.000000Z"),
			now: musTim("2023-09-28T14:23:24.161982Z"),
			tp1: musTim("2023-10-05T00:00:00.000000Z"),
		},
		// Case 027
		{
			fmt: "7 days",
			tm1: musTim("2023-09-28T00:00:00.000000Z"),
			now: musTim("2023-09-29T14:23:24.161982Z"),
			tp1: musTim("2023-10-05T00:00:00.000000Z"),
		},
		// Case 028
		{
			fmt: "2 weeks",
			tm1: musTim("2023-09-18T00:00:00.000000Z"), // week 38, Monday
			now: musTim("2023-09-28T14:23:24.161982Z"), // week 39, Thursday
			tp1: musTim("2023-10-02T00:00:00.000000Z"), // week 40, Monday
		},
		// Case 029
		{
			fmt: "2 weeks",
			tm1: musTim("2023-09-18T00:00:00.000000Z"), // week 38, Monday
			now: musTim("2023-09-24T14:23:24.161982Z"), // week 38, Sunday
			tp1: musTim("2023-10-02T00:00:00.000000Z"), // week 40, Monday
		},
		// Case 030
		{
			fmt: "2 weeks",
			tm1: musTim("2023-09-04T00:00:00.000000Z"), // week 36, Monday
			now: musTim("2023-09-18T00:00:00.000000Z"), // week 38, Monday
			tp1: musTim("2023-10-02T00:00:00.000000Z"), // week 40, Monday
		},
		// Case 031
		{
			fmt: "2 weeks",
			tm1: musTim("2023-09-18T00:00:00.000000Z"), // week 38, Monday
			now: musTim("2023-09-18T14:23:24.161982Z"), // week 38, Monday
			tp1: musTim("2023-10-02T00:00:00.000000Z"), // week 40, Monday
		},
		// Case 032
		{
			fmt: "2 weeks",
			tm1: musTim("2023-09-18T00:00:00.000000Z"), // week 38, Monday
			now: musTim("2023-09-18T00:00:00.161982Z"), // week 38, Monday
			tp1: musTim("2023-10-02T00:00:00.000000Z"), // week 40, Monday
		},
		// Case 033
		{
			fmt: "2 weeks",
			tm1: musTim("2022-12-12T00:00:00.000000Z"), // week 50, Monday (week of previous year)
			now: musTim("2022-12-19T00:00:00.000000Z"), // week 51, Monday (week of previous year)
			tp1: musTim("2022-12-26T00:00:00.000000Z"), // week 52, Monday (week of previous year)
		},
		// Case 034
		{
			fmt: "2 weeks",
			tm1: musTim("2022-12-12T00:00:00.000000Z"), // week 50, Monday (week of previous year)
			now: musTim("2022-12-19T14:23:24.161982Z"), // week 51, Monday (week of previous year)
			tp1: musTim("2022-12-26T00:00:00.000000Z"), // week 52, Monday (week of previous year)
		},
		// Case 035
		{
			fmt: "2 weeks",
			tm1: musTim("2022-12-12T00:00:00.000000Z"), // week 50, Monday (week of previous year)
			now: musTim("2022-12-26T00:00:00.000000Z"), // week 52, Monday (week of previous year)
			tp1: musTim("2023-01-09T00:00:00.000000Z"), // week  2, Monday
		},
		// Case 036
		{
			fmt: "2 weeks",
			tm1: musTim("2022-12-26T00:00:00.000000Z"), // week 52, Monday (week of previous year)
			now: musTim("2022-12-26T14:23:24.161982Z"), // week 52, Monday (week of previous year)
			tp1: musTim("2023-01-09T00:00:00.000000Z"), // week  2, Monday
		},
		// Case 037
		{
			fmt: "2 weeks",
			tm1: musTim("2022-12-26T00:00:00.000000Z"), // week 52, Monday (week of previous year)
			now: musTim("2023-01-01T14:23:24.161982Z"), // week 52, Sunday (week of previous year)
			tp1: musTim("2023-01-09T00:00:00.000000Z"), // week  2, Monday
		},
		// Case 038
		{
			fmt: "2 weeks",
			tm1: musTim("2022-12-26T00:00:00.000000Z"), // week 52, Monday (week of previous year)
			now: musTim("2023-01-08T14:23:24.161982Z"), // week  1, Sunday
			tp1: musTim("2023-01-09T00:00:00.000000Z"), // week  2, Monday
		},
		// Case 039
		{
			fmt: "3 months",
			tm1: musTim("2023-07-01T00:00:00.000000Z"),
			now: musTim("2023-09-28T14:23:24.161982Z"),
			tp1: musTim("2023-10-01T00:00:00.000000Z"),
		},
		// Case 040
		{
			fmt: "3 months",
			tm1: musTim("2023-01-01T00:00:00.000000Z"),
			now: musTim("2023-02-28T14:23:24.161982Z"),
			tp1: musTim("2023-04-01T00:00:00.000000Z"),
		},
		// Case 041
		{
			fmt: "3 months",
			tm1: musTim("2022-10-01T00:00:00.000000Z"),
			now: musTim("2023-01-01T00:00:00.000000Z"),
			tp1: musTim("2023-04-01T00:00:00.000000Z"),
		},
		// Case 042
		{
			fmt: "3 months",
			tm1: musTim("2023-07-01T00:00:00.000000Z"),
			now: musTim("2023-10-01T00:00:00.000000Z"),
			tp1: musTim("2024-01-01T00:00:00.000000Z"),
		},
		// Case 043
		{
			fmt: "3 months",
			tm1: musTim("2023-10-01T00:00:00.000000Z"),
			now: musTim("2023-10-01T14:23:24.161982Z"),
			tp1: musTim("2024-01-01T00:00:00.000000Z"),
		},
	}

	for i, tc := range testCases {
		t.Run(fmt.Sprintf("%03d", i), func(t *testing.T) {
			var tic *Ticker
			{
				tic = New(tc.fmt, tc.now)
			}

			var tm1 time.Time
			{
				tm1 = tic.TickM1()
			}

			if !tm1.Equal(tc.tm1) {
				t.Fatalf("%s tick-1\n\n%s\n", tc.fmt, cmp.Diff(tc.tm1, tm1))
			}

			var tp1 time.Time
			{
				tp1 = tic.TickP1()
			}

			if !tp1.Equal(tc.tp1) {
				t.Fatalf("%s tick+1\n\n%s\n", tc.fmt, cmp.Diff(tc.tp1, tp1))
			}
		})
	}
}

func musTim(str string) time.Time {
	tim, err := time.Parse("2006-01-02T15:04:05.999999Z", str)
	if err != nil {
		panic(err)
	}

	return tim
}
