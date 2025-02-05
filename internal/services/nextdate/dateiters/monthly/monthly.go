package monthly

import (
	"strconv"
	"strings"
	"time"
)

const timeStepPattern = `^m (-?[1-9]|-1|-2|[12][0-9]|3[01])(,(-?[1-9]|-1|-2|[12][0-9]|3[01]))*( (1[0-2]|[1-9])(,(1[0-2]|[1-9]))*)?$`

type Monthly struct {
	Days   []int
	Months map[int]bool
}

func New(timeStep string) Monthly {
	monthly := Monthly{}

	slicedStep := strings.Split(timeStep, " ")
	days := strings.Split(slicedStep[1], ",")
	monthly.Days = make([]int, len(days))
	for i, day := range days {
		monthly.Days[i], _ = strconv.Atoi(day)
	}

	if len(slicedStep) < 3 {
		return monthly
	}

	months := strings.Split(slicedStep[2], ",")
	monthly.Months = map[int]bool{}
	for _, month := range months {
		convertedMonth, _ := strconv.Atoi(month)
		if _, wasFound := monthly.Months[convertedMonth]; !wasFound {
			monthly.Months[convertedMonth] = true
		}
	}

	return monthly
}

func (iter Monthly) Next(now, startDate time.Time) time.Time {
	// TODO: Make Next()
	return startDate
}
