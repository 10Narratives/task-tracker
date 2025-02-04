package nextdate

import (
	"errors"
	"time"

	"github.com/10Narratives/task-tracker/internal/services/nextdate/dateiters/daily"
	"github.com/10Narratives/task-tracker/internal/services/nextdate/dateiters/monthly"
	"github.com/10Narratives/task-tracker/internal/services/nextdate/dateiters/weekly"
	"github.com/10Narratives/task-tracker/internal/services/nextdate/dateiters/yearly"
)

const DateLayout = "20060102"

var (
	ErrEmptyRepeat       = errors.New("empty repeat")
	ErrUnsupportedOption = errors.New("unsupported option")
)

type DateIterator interface {
	Next(now, startDate time.Time) time.Time
}

func newDateIterator(repeat string) (DateIterator, error) {
	var dateIter DateIterator
	var err error = nil
	switch repeat[0] {
	case 'd':
		dateIter, err = daily.New(repeat)
		break
	case 'w':
		dateIter, err = weekly.New(repeat)
		break
	case 'm':
		dateIter, err = monthly.New(repeat)
		break
	case 'y':
		dateIter, err = yearly.New(repeat)
		break
	default:
		err = ErrUnsupportedOption
	}

	return dateIter, err
}

// NextDate calculates the next occurrence of a date based on a given repetition pattern.
//
// Parameters:
//   - now: The current time, used as a reference.
//   - date: The starting date from which the iteration begins.
//   - repeat: A string defining the repetition pattern. Supported formats:
//   - "d <number>" → Advances by a fixed number of days (1-400).
//   - "w <number>" → Advances to specific weekdays (1 = Monday, 7 = Sunday).
//   - "m <days> [months]" → Advances to specific days of the month (1-31, -1 for last day).
//   - "y" → Advances by one year.
//
// Returns:
//   - A string representing the next calculated date in the format defined by DateLayout.
//   - An error if the repeat pattern is empty or unsupported.
//
// Errors:
//   - ErrEmptyRepeat: If the repeat string is empty.
//   - ErrUnsupportedOption: If the repeat pattern is invalid.
func NextDate(now, date time.Time, repeat string) (string, error) {
	if len(repeat) == 0 {
		return "", ErrEmptyRepeat
	}

	dateIter, err := newDateIterator(repeat)
	if err != nil {
		return "", err
	}

	return dateIter.Next(now, date).Format(DateLayout), nil
}
