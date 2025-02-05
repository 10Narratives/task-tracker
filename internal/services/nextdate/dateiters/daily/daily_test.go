package daily_test

import (
	"testing"
	"time"

	"github.com/10Narratives/task-tracker/internal/services/nextdate/dateiters/daily"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
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
				timeStep: `d 7`,
			},
			wantIter: func(tt require.TestingT, got interface{}, _ ...interface{}) {
				iter, ok := got.(daily.Daily)
				require.True(t, ok)

				assert.Equal(t, 7, iter.DayStep)
			},
		},
		{
			name: "successful initialization",
			args: args{
				timeStep: `d 50`,
			},
			wantIter: func(tt require.TestingT, got interface{}, _ ...interface{}) {
				iter, ok := got.(daily.Daily)
				require.True(t, ok)

				assert.Equal(t, 50, iter.DayStep)
			},
		},
		{
			name: "successful initialization",
			args: args{
				timeStep: `d 11`,
			},
			wantIter: func(tt require.TestingT, got interface{}, _ ...interface{}) {
				iter, ok := got.(daily.Daily)
				require.True(t, ok)

				assert.Equal(t, 11, iter.DayStep)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			daily := daily.New(tt.args.timeStep)
			tt.wantIter(t, daily)
		})
	}
}

func TestDaily_Next(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		timeStep  string
		now       time.Time
		startDate time.Time
		wantDate  time.Time
	}{
		{"Now is equal to startDate", "d 1", time.Date(2024, 2, 1, 0, 0, 0, 0, time.UTC), time.Date(2024, 2, 1, 0, 0, 0, 0, time.UTC), time.Date(2024, 2, 1, 0, 0, 0, 0, time.UTC)},
		{"Day step 7", "d 7", time.Date(2024, 2, 2, 0, 0, 0, 0, time.UTC), time.Date(2024, 2, 1, 0, 0, 0, 0, time.UTC), time.Date(2024, 2, 8, 0, 0, 0, 0, time.UTC)},
		{"Now before startDate", "d 10", time.Date(2024, 2, 1, 0, 0, 0, 0, time.UTC), time.Date(2024, 2, 5, 0, 0, 0, 0, time.UTC), time.Date(2024, 2, 5, 0, 0, 0, 0, time.UTC)},
		{"Large step", "d 28", time.Date(2025, 2, 2, 0, 0, 0, 0, time.UTC), time.Date(2025, 2, 1, 0, 0, 0, 0, time.UTC), time.Date(2025, 3, 1, 0, 0, 0, 0, time.UTC)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			daily := daily.New(tt.timeStep)

			newDate := daily.Next(tt.now, tt.startDate)
			assert.Equal(t, tt.wantDate, newDate)
		})
	}
}
