package nextdate

import (
	"time"

	"github.com/10Narratives/task-tracker/internal/services/nextdate"
)

const YearlyRule string = `^y$`

type Yearly struct{}

func NewYearly(repeatRule string) (Yearly, error) {
	var yearly = Yearly{}
	return yearly, nextdate.Validate(repeatRule, YearlyRule)
}

func (iter Yearly) NextDate(startDate string) (string, error) {
	parsedStartDate, err := nextdate.StringToTime(startDate)
	if err != nil {
		return "", err
	}

	now := time.Now()
	for now.After(parsedStartDate) {
		parsedStartDate = parsedStartDate.AddDate(1, 0, 0)
	}

	return parsedStartDate.Format(nextdate.DateLayout), nil
}
