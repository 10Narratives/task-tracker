package yearly

import (
	"time"
)

const timeStepPattern = `^y$`

// Yearly represents an iterator that advances by one year at a time.
type Yearly struct{}

// New creates a new Yearly iterator from a given time step string.
func New(timeStep string) Yearly {
	return Yearly{}
}

// Next calculates the next occurrence based on the current time and the start date.
// It advances from startDate in increments of one year until it reaches or surpasses now.
func (iter Yearly) Next(now, startDate time.Time) time.Time {
	for now.After(startDate) {
		startDate = startDate.AddDate(1, 0, 0)
	}
	return startDate
}
