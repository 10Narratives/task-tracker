package validation

import (
	"errors"
	"regexp"
)

var TimeStepPatterns = map[rune]string{
	'd': "^d (?:[1-9]|[1-9][0-9]|[1-3][0-9]{2}|400)$",
	'y': "^y$",
	'w': "^w ([1-7](,[1-7])*)$",
	'm': "^m (-?[1-9]|-1|-2|[12][0-9]|3[01])(,(-?[1-9]|-1|-2|[12][0-9]|3[01]))*( (1[0-2]|[1-9])(,(1[0-2]|[1-9]))*)?$",
}

var (
	ErrUnsupportedOption   = errors.New("unsupported option")
	ErrInvalidRepeatFormat = errors.New("invalid repeat format")
)

// Validate checks whether the given time step string matches the specified pattern.
//
// Parameters:
//   - timeStep: The input string representing the time step (e.g., "d 7", "w 1,3", "m 15").
//   - pattern: A regular expression pattern defining the valid format for the time step.
//
// Returns:
//   - nil if the timeStep matches the pattern.
//   - ErrInvalidTimeStepFormat if the timeStep does not conform to the expected format.
func Validate(repeat, pattern string) error {
	re := regexp.MustCompile(pattern)
	if !re.MatchString(repeat) {
		return ErrInvalidRepeatFormat
	}
	return nil
}
