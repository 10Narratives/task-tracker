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

// NextDate calculates the next occurrence of a date based on a given repetition pattern.
func NextDate(now, date time.Time, repeat string) string {
	nextdate := ""
	switch repeat[0] {
	case 'd':
		days, _ := strconv.Atoi(repeat[2:])
		nextdate = shiftDaily(now, date, days).Format(DateLayout)
	}
	return nextdate
}
