package nextdate

import (
	"time"

	"github.com/10Narratives/task-tracker/internal/services/nextdate/dateiters"
)

const DateLayout = "20060102"

// NextDate calculates the next occurrence of a date based on a given repetition pattern.
func NextDate(now, date time.Time, repeat string) string {
	return dateiters.NewDateIterator(repeat).Next(now, date).Format(DateLayout)
}
