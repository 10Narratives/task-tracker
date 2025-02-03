package dateiters

import (
	"errors"
	"fmt"
	"time"

	"github.com/10Narratives/task-tracker/pkg/nextdate"
)

// DailyTimeStepPattern is a regular expression pattern that matches valid daily time step formats.
// The pattern allows numbers from 1 to 400.
var DailyTimeStepPattern = `^d (?:[1-9]|[1-9][0-9]|[1-3][0-9]{2}|400)$`

// ErrInvalidDailyTimeStepFormat is returned when the provided daily time step does not
// match the expected format defined by DailyTimeStep.
// The valid format is "d [1-400]".
var ErrInvalidDailyTimeStepFormat = errors.New("The format for daily task is 'd [1-400]'")

// Daily represents a daily iterator with a specified step in days.
// The DayStep field determines the number of days to step forward
// in each iteration.
type Daily struct {
	// DayStep specifies the number of days to advance on each iteration.
	DayStep int
}

// NewDaily creates a new Daily iterator based on the provided timeStep string.
// It parses the number of days the iterator should step forward and validates the format.
//
// Parameters:
//   - timeStep: A string specifying the number of days to step forward, formatted as "d <number>".
//     Must matches the DailyTimeStepPattern.
//
// Returns:
//   - A DateIterator and an error. The DateIterator is set up based on the parsed DayStep,
//     or an error if the timeStep format is invalid.
func NewDaily(timeStep string) (nextdate.DateIterator, error) {
	daily := Daily{}
	if err := nextdate.Validate(timeStep, DailyTimeStepPattern); err != nil {
		return daily, errors.Join(err, ErrInvalidDailyTimeStepFormat)
	}
	fmt.Sscanf(timeStep, "d %d", &daily.DayStep)
	return daily, nil
}

// Next computes the next date based on the current date (now) and the starting date (date).
// It iterates by adding the specified number of days (DayStep) until the new date is not
// before the current date.
//
// Parameters:
// - now: The current date and time.
// - date: The starting date from which to calculate the next date.
//
// Returns:
// - The next valid date that is not before the current date.
func (iter Daily) Next(now, date time.Time) time.Time {
	for now.After(date) {
		date = date.AddDate(0, 0, iter.DayStep)
	}
	return date
}
