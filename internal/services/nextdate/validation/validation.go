package validation

import (
	"errors"
	"regexp"
)

var ErrInvalidTimeStepFormat = errors.New("invalid time step format")

func Validate(timeStep, pattern string) error {
	re := regexp.MustCompile(pattern)
	if !re.MatchString(timeStep) {
		return ErrInvalidTimeStepFormat
	}
	return nil
}
