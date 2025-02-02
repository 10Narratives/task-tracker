package nextdate

import (
	"errors"
	"regexp"
	"time"
)

var (
	DateLayout              string = "20060102"
	ErrCanNotParseStartDate error  = errors.New("can not parse start date")
	ErrTimeStepIsNotValid   error  = errors.New("time step is not valid")
)

type DateIterator interface {
	Next(now, startDate time.Time) time.Time
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
	_ = startDate

	err = Validate(timeStep, DailyRepeatRule)
	if err != nil {
		return "", err
	}
	daily, _ := NewDaily(timeStep)
	date = daily.Next(now, startDate).Format(DateLayout)

	return date, nil
}
