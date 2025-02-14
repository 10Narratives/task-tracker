package nextdate

import (
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/10Narratives/task-tracker/internal/lib"
)

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

func shiftWeekly(base time.Time, weekdays []int) time.Time {
	baseWeekday := int(base.Weekday())
	if baseWeekday == 0 {
		baseWeekday = 7
	}
	for _, day := range weekdays {
		if day > baseWeekday {
			return base.AddDate(0, 0, 7-baseWeekday)
		}
	}
	return base.AddDate(0, 0, 7+weekdays[0]-baseWeekday)
}

func daysInMonth(month time.Month, year int) int {
	return time.Date(year, month+1, 0, 0, 0, 0, 0, time.UTC).Day()
}

func normalizeDays(days []int, month time.Month, year int) []int {
	totalDays := daysInMonth(month, year)
	for i := range len(days) {
		if days[i] == -1 || days[i] == -2 {
			days[i] += totalDays + 1
		}
	}
	sort.Ints(days)
	return days
}

func shiftMonthly(now time.Time, days []int, months []int) time.Time {
	currYear, currMonth, currDay := now.Date()

	allowedMonths := make(map[int]bool)
	for _, m := range months {
		allowedMonths[m] = true
	}

	for {
		if len(allowedMonths) > 0 && !allowedMonths[int(currMonth)] {
			currMonth++
			if currMonth > 12 {
				currMonth = 1
				currYear++
			}
			currDay = 0
			continue
		}

		monthDays := make([]int, len(days))
		copy(monthDays, days)
		monthDays = normalizeDays(monthDays, currMonth, currYear)
		index := sort.Search(len(days), func(i int) bool { return monthDays[i] > currDay })
		if index != len(days) {
			return time.Date(currYear, currMonth, monthDays[index], 0, 0, 0, 0, time.UTC)
		}

		currMonth++
		if currMonth == 12 {
			currMonth = 1
			currYear++
		}
		currDay = 0
	}
}

// NextDate calculates the next occurrence of a date based on a given repetition pattern.
func NextDate(now, date time.Time, repeat string) string {
	nextdate := ""
	switch repeat[0] {
	case 'd':
		days, _ := strconv.Atoi(repeat[2:])
		nextdate = shiftDaily(now, date, days).Format(lib.DateFormat)
	case 'y':
		nextdate = shiftYearly(now, date).Format(lib.DateFormat)
	case 'w':
		weekdays := strings.Split(repeat[2:], ",")
		days := make([]int, len(weekdays))
		for i, weekday := range weekdays {
			days[i], _ = strconv.Atoi(weekday)
		}
		sort.Ints(days)
		nextdate = shiftWeekly(now, days).Format(lib.DateFormat)
	case 'm':
		strs := strings.Split(repeat, " ")

		strDays := strings.Split(strs[1], ",")
		days := make([]int, len(strDays))
		for i, strDay := range strDays {
			days[i], _ = strconv.Atoi(strDay)
		}
		sort.Ints(days)

		months := make([]int, 0)
		if len(strs) > 2 {
			strMonths := strings.Split(strs[2], ",")
			for _, strMonth := range strMonths {
				month, _ := strconv.Atoi(strMonth)
				months = append(months, month)
			}
		}
		nextdate = shiftMonthly(now, days, months).Format(lib.DateFormat)
	}
	return nextdate
}
