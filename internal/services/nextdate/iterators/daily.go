package nextdate

import (
	"strconv"
	"strings"
	"time"

	"github.com/10Narratives/task-tracker/internal/services/nextdate"
)

const DailyRule string = `^d (?:[1-9]|[1-9][0-9]|[1-3][0-9]{2}|400)$`

type Daily struct {
	DayStep int
}

func NewDaily(repeatRule string) (Daily, error) {
	var (
		daily Daily = Daily{}
		err   error = nil
	)

	err = nextdate.Validate(repeatRule, DailyRule)
	if err == nil {
		sliced := strings.Split(repeatRule, " ")
		daily.DayStep, _ = strconv.Atoi(sliced[1])
	}

	return daily, err
}

func (iter Daily) NextDate(startDate string) (string, error) {
	parsedStartDate, err := nextdate.StringToTime(startDate)
	if err != nil {
		return "", err
	}

	now := time.Now()
	for now.After(parsedStartDate) {
		parsedStartDate = parsedStartDate.AddDate(0, 0, iter.DayStep)
	}

	return parsedStartDate.Format(nextdate.DateLayout), nil
}
