package nextdate

// TODO: Make docs for this date iterator

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

// TODO: Change startDate type on time.Time and delete StringToTime conversation from method
func (iter Daily) NextDate(startDate string) (string, error) {
	parsedStartDate, err := nextdate.StringToTime(startDate)
	if err != nil {
		return "", err
	}

	// TODO: Find better solution
	now := time.Now()
	for now.After(parsedStartDate) {
		parsedStartDate = parsedStartDate.AddDate(0, 0, iter.DayStep)
	}

	return parsedStartDate.Format(nextdate.DateLayout), nil
}
