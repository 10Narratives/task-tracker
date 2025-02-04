package nextdate

import (
	"errors"
	"time"

	"github.com/10Narratives/task-tracker/internal/services/nextdate/dateiters/daily"
	"github.com/10Narratives/task-tracker/internal/services/nextdate/dateiters/monthly"
	"github.com/10Narratives/task-tracker/internal/services/nextdate/dateiters/weekly"
	"github.com/10Narratives/task-tracker/internal/services/nextdate/dateiters/yearly"
)

const DateLayout = "20060201"

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
