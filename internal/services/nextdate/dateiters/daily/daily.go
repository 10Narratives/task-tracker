package daily

import (
	"fmt"
	"time"

	"github.com/10Narratives/task-tracker/internal/services/nextdate"
)

const timeStepPattern = `^d (?:[1-9]|[1-9][0-9]|[1-3][0-9]{2}|400)$`

// Daily represents an iterator that advances by a fixed number of days.
type Daily struct {
	DayStep int // The number of days to step forward.
}

// New creates a new Daily iterator from a given time step string.
func New(timeStep string) (Daily, error) {
	daily := Daily{}

	err := nextdate.Validate(timeStep, timeStepPattern)
	if err != nil {
		return daily, fmt.Errorf("%w: daily format is d <number> where number in range [1, 400]", err)
	}
	fmt.Sscanf(timeStep, "d %d", &daily.DayStep)

	return daily, nil
}

// Next calculates the next occurrence based on the current time and the start date.
// It advances from startDate in increments of DayStep days until it reaches or surpasses now.
func (iter Daily) Next(now, startDate time.Time) time.Time {
	for now.After(startDate) {
		startDate = startDate.AddDate(0, 0, iter.DayStep)
	}
	return startDate
}
