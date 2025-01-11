package authorizenet

import "time"

func Now() time.Time {
	current_time := time.Now().UTC()
	return current_time
}

func LastWeek() time.Time {
	t := time.Now().UTC().AddDate(0, 0, -7)
	return t
}

func LastMonth() time.Time {
	t := time.Now().UTC().AddDate(0, -1, 0)
	return t
}

func LastYear() time.Time {
	t := time.Now().UTC().AddDate(-1, 0, 0)
	return t
}

func CurrentDate() string {
	currentTime := time.Now().UTC()
	return currentTime.Format("2006-01-02")
}

func IntervalMonthly() Interval {
	return Interval{Length: "1", Unit: "months"}
}

func IntervalQuarterly() Interval {
	return Interval{Length: "3", Unit: "months"}
}

func IntervalWeekly() Interval {
	return Interval{Length: "7", Unit: "days"}
}

func IntervalDays(amount string) Interval {
	return Interval{Length: amount, Unit: "days"}
}

func IntervalMonths(amount string) Interval {
	return Interval{Length: amount, Unit: "months"}
}

func IntervalYearly() Interval {
	return Interval{Length: "365", Unit: "days"}
}
