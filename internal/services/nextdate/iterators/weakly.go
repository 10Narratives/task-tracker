package nextdate

// TODO: Make docs for this date iterator

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

// TODO: Refactor names in this function
func Convert(base time.Weekday) int {
	converted := base
	if base == 0 {
		converted = 7
	}
	return int(converted)
}

// TODO: Change startDate type on time.Time and delete StringToTime conversation from method
func (iter Weakly) NextDate(startDate string) (string, error) {
	parsedStartDate, err := nextdate.StringToTime(startDate)
	if err != nil {
		return "", err
	}

	// TODO: Find solution for calibration by O(1)
	now := time.Now()
	for now.After(parsedStartDate) {
		parsedStartDate = parsedStartDate.AddDate(0, 0, 1)
	}

	diff := (iter.TargetWeekday - Convert(parsedStartDate.Weekday())) % 7
	parsedStartDate = parsedStartDate.AddDate(0, 0, diff)

	return parsedStartDate.Format(nextdate.DateLayout), nil
}
