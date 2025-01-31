package service

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"time"
)

/*
d n - повторение задачи через определённое количество дней (1 <= n <= 400)

y - повторение задачи через год (не требует дополнительных параметров)

w <через запятую от 1 до 7> - задача назначается в указанные дни недели, где 1 - это понедельник, 7 - это воскресенье
m <через запятую от 1 до 31, -1, -2> [через запятую от 1 до 12] задача назначается в указанные дни месяца.
При этом вторая последовательность чисел опциональна и указывает на определённые месяцы.

*/

const layout string = "20060102"

var regexps map[rune]string = map[rune]string{
	'd': `^d (?:[1-9]|[1-9][0-9]|[1-3][0-9]{2}|400)$`,
	'y': `^y$`,
	'w': `^w [1-7]$`,
	'm': `^m ((?:-1|-2|[1-9]|[1-2][0-9]|3[0-1])(?:,(?:-1|-2|[1-9]|[1-2][0-9]|3[0-1]))*)(?: ([1-9]|1[0-2])(?:,([1-9]|1[0-2]))*)?$`,
}

func isRepeatValid(repeat string) error {
	if len(repeat) < 1 {
		return errors.New("Can not compute next date. Repeat rule is empty.")
	}
	option := rune(repeat[0])

	rule, wasFound := regexps[option]
	if !wasFound {
		return errors.New("Can not compute next date. Option is unknown")
	}

	re, err := regexp.Compile(rule)
	if err != nil {
		return fmt.Errorf("Can not compile regexp, %w", err)
	}

	if !re.MatchString(repeat) {
		return errors.New("Gotten repeat rule is illegal")
	}
	return nil
}

func computeDayRepeat(now time.Time, date time.Time, dayRepeat int) string {
	for now.After(date) {
		date = date.AddDate(0, 0, dayRepeat)
	}
	return date.Format(layout)
}

func computeYearRepeat(now time.Time, date time.Time) string {
	for now.After(date) {
		date = date.AddDate(1, 0, 0)
	}
	return date.Format(layout)
}

func NextDate(now time.Time, date string, repeat string) (string, error) {
	err := isRepeatValid(repeat)
	if err != nil {
		return "", err
	}

	parsedDate, err := time.Parse(layout, date)
	if err != nil {
		return "", fmt.Errorf("Can not parse date, %w", err)
	}

	nextDate, option := "", rune(repeat[0])
	switch option {
	case 'd':
		dayRepeat, _ := strconv.Atoi(repeat[2:])
		nextDate = computeDayRepeat(now, parsedDate, dayRepeat)
		break
	case 'y':
		nextDate = computeYearRepeat(now, parsedDate)
	}

	return nextDate, nil
}
