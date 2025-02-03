package dateiters

import (
	"errors"
	"time"

	"github.com/10Narratives/task-tracker/pkg/nextdate"
)

var YearlyTimeStepPattern = `y`

var ErrInvalidYearlyTimeStepFormat = errors.New("The format fir yearly tasks is `y`")

type Yearly struct{}

func NewYearly(timeStep string) (nextdate.DateIterator, error) {
	if err := nextdate.Validate(timeStep, YearlyTimeStepPattern); err != nil {
		return Yearly{}, errors.Join(err, ErrInvalidYearlyTimeStepFormat)
	}
	return Yearly{}, nil
}

func (iter Yearly) Next(now, date time.Time) time.Time {
	for now.After(date) {
		date = date.AddDate(1, 0, 0)
	}
	return date
}
