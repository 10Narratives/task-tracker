package nextdate

import (
	"fmt"
	"time"
)

var DailyRepeatRule = `^d (?:[1-9]|[1-9][0-9]|[1-3][0-9]{2}|400)$`

type Daily struct {
	DayStep int
}

func NewDaily(timeStep string) (Daily, error) {
	daily := Daily{}
	if err := Validate(timeStep, DailyRepeatRule); err != nil {
		return daily, err
	}
	fmt.Sscanf(timeStep, "d %d", &daily.DayStep)
	return daily, nil
}

func (iter Daily) Next(now, startDate time.Time) time.Time {
	for now.After(startDate) {
		startDate = startDate.AddDate(0, 0, iter.DayStep)
	}
	return startDate
}
