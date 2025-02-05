package dateiters

import (
	"time"

	"github.com/10Narratives/task-tracker/internal/services/nextdate/dateiters/daily"
	"github.com/10Narratives/task-tracker/internal/services/nextdate/dateiters/monthly"
	"github.com/10Narratives/task-tracker/internal/services/nextdate/dateiters/weekly"
	"github.com/10Narratives/task-tracker/internal/services/nextdate/dateiters/yearly"
)

type DateIterator interface {
	Next(now, date time.Time) time.Time
}

func NewDateIterator(repeat string) DateIterator {
	var dateIter DateIterator
	switch repeat[0] {
	case 'd':
		dateIter = daily.New(repeat)
		break
	case 'w':
		dateIter = weekly.New(repeat)
		break
	case 'm':
		dateIter = monthly.New(repeat)
		break
	case 'y':
		dateIter = yearly.New(repeat)
		break
	}

	return dateIter
}
