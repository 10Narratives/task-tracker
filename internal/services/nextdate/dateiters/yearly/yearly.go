package yearly

import (
	"fmt"
	"time"

	"github.com/10Narratives/task-tracker/internal/services/nextdate"
)

const TimeStepPattern = `y`

// Yearly represents an iterator that advances by one year at a time.
type Yearly struct{}

// New creates a new Yearly iterator from a given time step string.
// The time step must be exactly "y".
// Returns an error if the format is invalid.
func New(timeStep string) (Yearly, error) {
	yearly := Yearly{}
	err := nextdate.Validate(timeStep, TimeStepPattern)
	if err != nil {
		err = fmt.Errorf("%w: yearly format is y", err)
	}
	return yearly, err
}

// Next calculates the next occurrence based on the current time and the start date.
// It advances from startDate in increments of one year until it reaches or surpasses now.
func (iter Yearly) Next(now, startDate time.Time) time.Time {
	for now.After(startDate) {
		startDate = startDate.AddDate(1, 0, 0)
	}
	return startDate
}
