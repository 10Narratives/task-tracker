package nextdate

import (
	"slices"
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

func shiftMonthly(base time.Time, days []int, months []int) time.Time {
	year, month, day := base.Date()
	loc := base.Location()

	// Сортируем список дней, чтобы найти ближайший
	sort.Ints(days)

	for i := 0; i < 24; i++ { // Проверяем ближайшие 24 месяца
		// Проверяем, подходит ли месяц
		if len(months) > 0 && !slices.Contains(months, int(month)) {
			month++
			if month > 12 {
				year++
				month = 1
			}
			continue
		}

		// Определяем последний день месяца
		lastDay := time.Date(year, month+1, 0, 0, 0, 0, 0, loc).Day()

		// Проверяем возможные дни в этом месяце
		for _, d := range days {
			targetDay := d
			if d == -1 {
				targetDay = lastDay
			} else if d == -2 {
				targetDay = lastDay - 1
			}

			// Если день корректен и не в прошлом, возвращаем его
			if (year > base.Year() || month > base.Month() || targetDay >= day) && targetDay <= lastDay {
				return time.Date(year, month, targetDay, 0, 0, 0, 0, loc)
			}
		}

		// Переход на следующий месяц
		month++
		if month > 12 {
			year++
			month = 1
		}
	}

	return base // Если не найдено, возвращаем исходную дату
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
