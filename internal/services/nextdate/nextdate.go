package nextdate

import (
	"errors"
	"regexp"
	"time"
)

var (
	ErrTimeStepIsNotValid error = errors.New("time step is not valid")
)

type DateIterator interface {
	Next(startDate time.Time) time.Time
}

func Validate(timeStep string, repeatRule string) error {
	re := regexp.MustCompile(repeatRule)
	if !re.MatchString(timeStep) {
		return ErrTimeStepIsNotValid
	}
	return nil
}
