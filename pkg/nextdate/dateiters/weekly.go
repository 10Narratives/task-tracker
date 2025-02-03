package dateiters

import (
	"errors"
	"strconv"
	"strings"
	"time"

	"github.com/10Narratives/task-tracker/pkg/nextdate"
)

// WeeklyTimeStepPattern defines the regular expression used to validate
// the format of the weekly time step. Valid formats include "w [1-7]"
// where the numbers represent the days of the week (1 for Monday to 7 for Sunday).
var WeeklyTimeStepPattern = `^w ([1-7](,[1-7])*)$`

// ErrInvalidWeeklyTimeStepFormat is the error returned when the
// format of the provided weekly time step does not match the expected pattern.
var ErrInvalidWeeklyTimeStepFormat = errors.New("The format for weekly task is 'w [1-7]'")

// Weekly represents an iterator that steps through specific weekdays.
// It contains a slice of integers that represent the weekdays to step through.
type Weekly struct {
	Weekdays []int
}

// NewWeekly creates a new Weekly iterator based on the provided timeStep string.
// It validates the format and extracts the weekdays to step forward.
//
// Parameters:
//   - timeStep: A string specifying the weekdays to step forward,
//     formatted as "w [1-7]". Multiple weekdays can be specified
//     as a comma-separated list (e.g., "w 1,3,5"). Must matches WeeklyTimeStepPattern
//
// Returns:
//   - A DateIterator and an error. The DateIterator is populated
//     with the validated weekdays or an error if the timeStep format is invalid.
func NewWeekly(timeStep string) (nextdate.DateIterator, error) {
	weekly := Weekly{}
	if err := nextdate.Validate(timeStep, WeeklyTimeStepPattern); err != nil {
		return weekly, errors.Join(err, ErrInvalidWeeklyTimeStepFormat)
	}
	weekdays := strings.Split(timeStep[2:], ",")
	weekly.Weekdays = make([]int, len(weekdays))
	for i, weekday := range weekdays {
		weekly.Weekdays[i], _ = (strconv.Atoi(weekday))
	}
	return weekly, nil
}

// Next computes the next eligible date based on the current date (now) and the starting date (date).
// It calculates how many days to add based on the specified weekdays and returns the next valid date.
//
// Parameters:
// - now: The current date and time.
// - date: The starting date from which to calculate the next date.
//
// Returns:
// - The next valid date that is not before the current date.
func (iter Weekly) Next(now, date time.Time) time.Time {
	startWeekday := int(date.Weekday())
	if startWeekday == 0 {
		startWeekday = 7 // Treat Sunday as 7
	}

	minDaysToAdd := 7
	for _, targetDay := range iter.Weekdays {
		diff := (targetDay - startWeekday + 7) % 7

		// If the date is today and we are looking for future dates, continue
		if diff == 0 && now.Before(date) {
			return date
		}

		// Update days to add to the minimum found
		if diff != 0 {
			minDaysToAdd = min(minDaysToAdd, diff)
		}
	}

	return date.AddDate(0, 0, minDaysToAdd)
}
