package nextdate

import (
	"errors"
	"regexp"
	"time"
)

type DateIterator interface {
	Next(now, startDate time.Time) time.Time
}

var ErrInvalidTimeStepFormat = errors.New("invalid time step format")

func Validate(timeStep, pattern string) error {
	re := regexp.MustCompile(pattern)
	if !re.MatchString(timeStep) {
		return ErrInvalidTimeStepFormat
	}
	return nil
}

func NextDate(now, date time.Time, repeat string) (string, error) {
	// Validate repeat

	return "", nil
}
