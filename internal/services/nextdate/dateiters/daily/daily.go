package daily

import (
	"fmt"
	"time"
)

// Daily represents an iterator that advances by a fixed number of days.
type Daily struct {
	DayStep int // The number of days to step forward.
}

// New creates a new Daily iterator from a given time step string.
func New(timeStep string) Daily {
	daily := Daily{}
	fmt.Sscanf(timeStep, "d %d", &daily.DayStep)
	return daily
}

// Next calculates the next occurrence based on the current time and the start date.
// It advances from startDate in increments of DayStep days until it reaches or surpasses now.
func (iter Daily) Next(now, startDate time.Time) time.Time {
	for now.After(startDate) {
		startDate = startDate.AddDate(0, 0, iter.DayStep)
	}
	return startDate
}
