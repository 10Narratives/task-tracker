package dateiters

import (
	"errors"
	"strconv"
	"strings"
	"time"

	"github.com/10Narratives/task-tracker/internal/services/nextdate"
)

var MonthlyTimeStepPattern = `^m (-?[1-9]|-1|-2|[12][0-9]|3[01])(,(-?[1-9]|-1|-2|[12][0-9]|3[01]))*( (1[0-2]|[1-9])(,(1[0-2]|[1-9]))*)?$`

var ErrInvalidMonthlyTimeStepFormat = errors.New("The format for monthly task is `m [1-31, -1, -2] [1-12]`")

type Monthly struct {
	Days   []int
	Months []int
}

func NewMonthly(timeStep string) (nextdate.DateIterator, error) {
	monthly := Monthly{}
	if err := nextdate.Validate(timeStep, MonthlyTimeStepPattern); err != nil {
		return monthly, errors.Join(err, ErrInvalidMonthlyTimeStepFormat)
	}

	sliced := strings.Split(timeStep, " ")

	days := strings.Split(sliced[1], ",")
	monthly.Days = make([]int, len(days))
	for i, day := range days {
		monthly.Days[i], _ = strconv.Atoi(day)
	}

	months := []string{}
	if len(sliced) > 2 {
		months = strings.Split(sliced[2], ",")
	}
	monthly.Months = make([]int, len(months))
	for i, month := range months {
		monthly.Months[i], _ = strconv.Atoi(month)
	}

	return monthly, nil
}

// m n
func (iter Monthly) Next(now, date time.Time) time.Time {

	return date
}
