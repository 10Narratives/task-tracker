package nextdate

import (
	"strconv"
	"strings"
	"time"

	"github.com/10Narratives/task-tracker/internal/services/nextdate"
)

const WeaklyRule string = `^w [1-7]&`

type Weakly struct {
	TargetWeekday int
}

func NewWeakly(repeatRule string) (Weakly, error) {
	var (
		weakly Weakly = Weakly{}
		err    error  = nil
	)

	err = nextdate.Validate(repeatRule, WeaklyRule)
	if err == nil {
		sliced := strings.Split(repeatRule, " ")
		weakly.TargetWeekday, _ = strconv.Atoi(sliced[1])
	}

	return weakly, err
}

func Convert(base time.Weekday) int {
	converted := base
	if base == 0 {
		converted = 7
	}
	return int(converted)
}

func (iter Weakly) NextDate(startDate string) (string, error) {
	parsedStartDate, err := nextdate.StringToTime(startDate)
	if err != nil {
		return "", err
	}

	now := time.Now()
	for now.After(parsedStartDate) {
		parsedStartDate = parsedStartDate.AddDate(0, 0, 1)
	}

	diff := (iter.TargetWeekday - Convert(parsedStartDate.Weekday())) % 7
	parsedStartDate = parsedStartDate.AddDate(0, 0, diff)

	return parsedStartDate.Format(nextdate.DateLayout), nil
}
