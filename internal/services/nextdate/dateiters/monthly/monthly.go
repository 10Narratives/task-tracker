package monthly

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/10Narratives/task-tracker/internal/services/nextdate"
)

const timeStepPattern = `^m (-?[1-9]|-1|-2|[12][0-9]|3[01])(,(-?[1-9]|-1|-2|[12][0-9]|3[01]))*( (1[0-2]|[1-9])(,(1[0-2]|[1-9]))*)?$`

type Monthly struct {
	Days   []int
	Months map[int]bool
}

func New(timeStep string) (Monthly, error) {
	monthly := Monthly{}
	err := nextdate.Validate(timeStep, timeStepPattern)
	if err != nil {
		return monthly, fmt.Errorf("%w: monthly format is m <1-31,-1,-2> [1-12]", err)
	}

	slicedStep := strings.Split(timeStep, " ")
	days := strings.Split(slicedStep[1], ",")
	monthly.Days = make([]int, len(days))
	for i, day := range days {
		monthly.Days[i], _ = strconv.Atoi(day)
	}

	if len(slicedStep) < 3 {
		return monthly, nil
	}

	months := strings.Split(slicedStep[2], ",")
	monthly.Months = map[int]bool{}
	for _, month := range months {
		convertedMonth, _ := strconv.Atoi(month)
		if _, wasFound := monthly.Months[convertedMonth]; !wasFound {
			monthly.Months[convertedMonth] = true
		}
	}

	return monthly, nil
}

func (iter Monthly) Next(now, startDate time.Time) time.Time {
	// TODO: Make Next()
	return startDate
}
