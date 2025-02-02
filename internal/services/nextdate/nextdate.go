package nextdate

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
	"time"
)

const DateLayout string = "20060102"

var (
	YearlyRepeatRule = `y`
	DailyRepeatRule  = `^d (?:[1-9]|[1-9][0-9]|[1-3][0-9]{2}|400)$`
)

var (
	ErrCanNotParseStartDate error = errors.New("can not parse start date")
	ErrTimeStepIsNotValid   error = errors.New("time step is not valid")
)

type DateIterator interface {
	Next(now, startDate time.Time) time.Time
}

type Yearly struct{}

func NewYearly(timeStep string) (DateIterator, error) {
	return Yearly{}, Validate(timeStep, YearlyRepeatRule)
}

func (iter Yearly) Next(now, startDate time.Time) time.Time {
	for now.After(startDate) {
		startDate = startDate.AddDate(1, 0, 0)
	}
	return startDate
}

type Daily struct {
	DayStep int
}

func NewDaily(timeStep string) (DateIterator, error) {
	daily := Daily{}
	if err := Validate(timeStep, DailyRepeatRule); err != nil {
		return daily, err
	}
	fmt.Sscanf(timeStep, "d %d", &daily.DayStep)
	return daily, nil
}

func (iter Daily) Next(now, startDate time.Time) time.Time {
	for now.After(startDate) {
		startDate = startDate.AddDate(0, 0, iter.DayStep)
	}
	return startDate
}

func Validate(timeStep string, repeatRule string) error {
	re := regexp.MustCompile(repeatRule)
	if !re.MatchString(timeStep) {
		return ErrTimeStepIsNotValid
	}
	return nil
}

func NextDate(now time.Time, date string, timeStep string) (string, error) {
	startDate, err := time.Parse(DateLayout, date)
	if err != nil {
		return "", errors.Join(ErrCanNotParseStartDate, err)
	}
	if len(strings.Trim(timeStep, " ")) < 1 {
		return "", nil // TODO: Make error instance
	}

	var iter DateIterator
	switch timeStep[0] {
	case 'd':
		iter, err = NewDaily(timeStep)
	case 'y':
		iter, err = NewYearly(timeStep)
	}
	if err != nil {
		return "", err
	}

	date = iter.Next(now, startDate).Format(DateLayout)

	return date, nil
}
