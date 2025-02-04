package weekly

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/10Narratives/task-tracker/internal/services/nextdate"
)

const TimeStepPattern = `^w ([1-7](,[1-7])*)$`

// Weekly represents an iterator that advances to specified weekdays.
type Weekly struct {
	Weekdays []int // Days of the week to iterate over (1 = Monday, 7 = Sunday).
}

// New creates a new Weekly iterator from a given time step string.
// The time step must follow the format "w [1-7]" where 1 = Monday and 7 = Sunday.
// Multiple days can be specified as comma-separated values (e.g., "w 1,3,5").
func New(timeStep string) (Weekly, error) {
	weekly := Weekly{}
	err := nextdate.Validate(timeStep, TimeStepPattern)
	if err != nil {
		return weekly, fmt.Errorf("%w: weekly format is `w [1-7]`", err)
	}

	weekdays := strings.Split(timeStep[2:], ",")
	weekly.Weekdays = make([]int, len(weekdays))
	for i, weekday := range weekdays {
		weekly.Weekdays[i], _ = strconv.Atoi(weekday)
	}

	return weekly, nil
}

// Next calculates the next occurrence based on the current time and the start date.
// It advances to the next available weekday from the Weekdays slice, ensuring that
// the returned date is always in the future relative to now.
func (iter Weekly) Next(now, startDate time.Time) time.Time {
	// TODO: Find better solution for calibration startDate
	for now.After(startDate) {
		startDate = startDate.AddDate(0, 0, 1)
	}

	startWeekday := int(startDate.Weekday())
	if startWeekday == 0 {
		startWeekday = 7 // Treat Sunday as 7
	}

	minDaysToAdd := 7
	for _, targetDay := range iter.Weekdays {
		diff := (targetDay - startWeekday + 7) % 7

		// If the date is today and we are looking for future dates, continue
		if diff == 0 && now.Before(startDate) {
			return startDate
		}

		// Update days to add to the minimum found
		if diff != 0 {
			minDaysToAdd = min(minDaysToAdd, diff)
		}
	}

	return startDate.AddDate(0, 0, minDaysToAdd)
}
