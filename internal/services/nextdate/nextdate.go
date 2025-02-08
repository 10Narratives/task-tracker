package nextdate

import (
	"fmt"
	"time"

	"github.com/10Narratives/task-tracker/internal/services/nextdate/dateiters"
)

const DateLayout = "20060102"

// NextDate calculates the next occurrence of a date based on a given repetition pattern.
func NextDate(now, date time.Time, repeat string) string {
	iter := dateiters.NewDateIterator(repeat)
	fmt.Println("iter", iter)
	nextTime := iter.Next(now, date)
	fmt.Println("nextTime", nextTime)
	result := nextTime.Format(DateLayout)
	fmt.Println("result", result)
	return result
}
