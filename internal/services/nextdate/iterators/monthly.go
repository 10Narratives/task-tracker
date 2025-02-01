package nextdate

// TODO: Make docs for this date iterator

import (
	"strconv"
	"strings"

	"github.com/10Narratives/task-tracker/internal/services/nextdate"
)

const MonthlyRule string = `^m ((?:-1|-2|[1-9]|[1-2][0-9]|3[0-1])(?:,(?:-1|-2|[1-9]|[1-2][0-9]|3[0-1]))*)(?: ([1-9]|1[0-2])(?:,([1-9]|1[0-2]))*)?$`

type Monthly struct {
	TargetDays   []int
	TargetMonths []int
}

// TODO: Refactor names or find standard alternative
func convert(source string) []int {
	sliced := strings.Split(source, ",")
	result := make([]int, len(sliced))
	for i, data := range sliced {
		result[i], _ = strconv.Atoi(data)
	}
	return result
}

func NewMonthly(repeatRule string) (Monthly, error) {
	var (
		monthly Monthly = Monthly{}
		err     error   = nil
	)

	err = nextdate.Validate(repeatRule, MonthlyRule)
	if err == nil {
		sliced := strings.Split(repeatRule, " ")
		monthly.TargetDays = convert(sliced[1])
		if len(sliced) > 2 {
			monthly.TargetMonths = convert(sliced[2])
		}

	}

	return monthly, err
}

// TODO: Make NextDate(base time.Time) time.Time
