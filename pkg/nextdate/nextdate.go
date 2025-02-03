package nextdate

import (
	"errors"
	"regexp"
	"time"
)

const DateLayout string = "20060102"

var (
	ErrCanNotParseStartDate  error = errors.New("can not parse start date")
	ErrInvalidTimeStepFormat error = errors.New("invalid time step format")
)

// DateIterator is an interface that defines methods for iterating through dates.
// Implementing types should provide a mechanism to get the next date based on the
// current date and a starting date, as well as validate a given time step.
type DateIterator interface {
	// Next returns the next date based on the provided current date (now)
	// and a starting date (startDate). The method should define the logic
	// for calculating the next date.
	Next(now, startDate time.Time) time.Time
}

func Validate(timeStep string, pattern string) error {
	re := regexp.MustCompile(pattern)
	if !re.MatchString(timeStep) {
		return ErrInvalidTimeStepFormat
	}
	return nil
}
