package validation

import (
	"errors"
	"regexp"
)

var ErrInvalidTimeStepFormat = errors.New("invalid time step format")

// Validate checks whether the given time step string matches the specified pattern.
//
// Parameters:
//   - timeStep: The input string representing the time step (e.g., "d 7", "w 1,3", "m 15").
//   - pattern: A regular expression pattern defining the valid format for the time step.
//
// Returns:
//   - nil if the timeStep matches the pattern.
//   - ErrInvalidTimeStepFormat if the timeStep does not conform to the expected format.
func Validate(timeStep, pattern string) error {
	re := regexp.MustCompile(pattern)
	if !re.MatchString(timeStep) {
		return ErrInvalidTimeStepFormat
	}
	return nil
}
