package nextdate

import (
	"strconv"
	"time"
)

const DateLayout = "20060102"

func shiftDaily(now, date time.Time, days int) time.Time {
	diffDays := int(now.Sub(date).Hours() / 24)

	if diffDays <= 0 {
		return date.AddDate(0, 0, days)
	}

	return date.AddDate(0, 0, diffDays+(days-diffDays%days))
}

func shiftYearly(now, date time.Time) time.Time {
	yearDiff := max(now.Year()-date.Year(), 0)
	if int(now.Month())-int(date.Month()) == 0 && (now.Day()-date.Day()) < 0 {
		return date.AddDate(yearDiff, 0, 0)
	} else if int(now.Month())-int(date.Month()) < 0 {
		return date.AddDate(yearDiff, 0, 0)
	}
	return date.AddDate(yearDiff+1, 0, 0)
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
