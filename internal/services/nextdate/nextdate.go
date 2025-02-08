package nextdate

import (
	"strconv"
	"time"
)

const DateLayout = "20060102"

func shiftDaily(now, date time.Time, days int) time.Time {
	for now.After(date) {
		date = date.AddDate(0, 0, days)
	}
	return date
}

func shiftYearly(now, date time.Time) time.Time {
	for now.After(date) {
		date = date.AddDate(1, 0, 0)
	}
	return date
}

// NextDate calculates the next occurrence of a date based on a given repetition pattern.
func NextDate(now, date time.Time, repeat string) string {
	nextdate := ""
	switch repeat[0] {
	case 'd':
		days, _ := strconv.Atoi(repeat[2:])
		nextdate = shiftDaily(now, date, days).Format(DateLayout)
	case 'y':
		nextdate = shiftYearly(now, date).Format(DateLayout)
	}
	return nextdate
}
