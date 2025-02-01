package nextdate

// TODO: Make docs for this date iterator

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

// TODO: Change startDate type on time.Time and delete StringToTime conversation from method
func (iter Yearly) NextDate(startDate string) (string, error) {
	parsedStartDate, err := nextdate.StringToTime(startDate)
	if err != nil {
		return "", err
	}

	// TODO: Find better solution
	now := time.Now()
	for now.After(parsedStartDate) {
		parsedStartDate = parsedStartDate.AddDate(1, 0, 0)
	}

	return parsedStartDate.Format(nextdate.DateLayout), nil
}
