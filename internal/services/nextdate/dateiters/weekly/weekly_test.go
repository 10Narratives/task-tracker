package weekly_test

import (
	"testing"
	"time"

	"github.com/10Narratives/task-tracker/internal/services/nextdate/dateiters/weekly"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWeekly(t *testing.T) {
	t.Parallel()

	type args struct {
		timeStep string
	}

	tests := []struct {
		name     string
		args     args
		wantIter require.ValueAssertionFunc
	}{
		{
			name: "successful initialization",
			args: args{
				timeStep: "w 1,4",
			},
			wantIter: func(tt require.TestingT, got interface{}, _ ...interface{}) {
				iter, ok := got.(weekly.Weekly)
				require.True(t, ok)

				assert.Len(t, iter.Weekdays, 2)
				assert.Equal(t, iter.Weekdays[0], 1)
				assert.Equal(t, iter.Weekdays[1], 4)
			},
		},
		{
			name: "successful initialization with single weekday",
			args: args{
				timeStep: "w 1",
			},
			wantIter: func(tt require.TestingT, got interface{}, _ ...interface{}) {
				iter, ok := got.(weekly.Weekly)
				require.True(t, ok)

				assert.Len(t, iter.Weekdays, 1)
				assert.Equal(t, iter.Weekdays[0], 1)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			weekly := weekly.New(tt.args.timeStep)
			tt.wantIter(t, weekly)
		})
	}
}

func TestWeekly_Next(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		timeStep  string
		now       time.Time
		startDate time.Time
		wantDate  time.Time
	}{
		{"Single day", "w 7", time.Date(2024, 1, 27, 0, 0, 0, 0, time.UTC), time.Date(2024, 1, 26, 0, 0, 0, 0, time.UTC), time.Date(2024, 1, 28, 0, 0, 0, 0, time.UTC)},
		{"Many days", "w 1,2,3", time.Date(2024, 1, 26, 0, 0, 0, 0, time.UTC), time.Date(2024, 1, 25, 0, 0, 0, 0, time.UTC), time.Date(2024, 1, 29, 0, 0, 0, 0, time.UTC)},
		{"Month ago", "w 1,2", time.Date(2025, 2, 5, 0, 0, 0, 0, time.UTC), time.Date(2025, 1, 25, 0, 0, 0, 0, time.UTC), time.Date(2025, 2, 10, 0, 0, 0, 0, time.UTC)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			weekly := weekly.New(tt.timeStep)

			newDate := weekly.Next(tt.now, tt.startDate)
			assert.Equal(t, tt.wantDate, newDate)
		})
	}
}
