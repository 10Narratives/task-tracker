package yearly_test

import (
	"testing"
	"time"

	"github.com/10Narratives/task-tracker/internal/services/nextdate/dateiters/yearly"
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
			name: "Successful initialization",
			args: args{
				timeStep: "y",
			},
			wantIter: func(tt require.TestingT, got interface{}, _ ...interface{}) {
				iter, ok := got.(yearly.Yearly)
				require.True(t, ok)
				assert.Equal(t, yearly.Yearly{}, iter)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			yearly := yearly.New(tt.args.timeStep)
			tt.wantIter(t, yearly)
		})
	}
}

func TestYearly_Next(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		timeStep  string
		now       time.Time
		startDate time.Time
		wantDate  time.Time
	}{
		{"Now is equal to startDate", "y", time.Date(2024, 2, 1, 0, 0, 0, 0, time.UTC), time.Date(2024, 2, 1, 0, 0, 0, 0, time.UTC), time.Date(2024, 2, 1, 0, 0, 0, 0, time.UTC)},
		{"Now before startDate", "y", time.Date(2024, 2, 1, 0, 0, 0, 0, time.UTC), time.Date(2024, 2, 2, 0, 0, 0, 0, time.UTC), time.Date(2024, 2, 2, 0, 0, 0, 0, time.UTC)},
		{"Small step", "y", time.Date(2024, 2, 2, 0, 0, 0, 0, time.UTC), time.Date(2024, 2, 1, 0, 0, 0, 0, time.UTC), time.Date(2025, 2, 1, 0, 0, 0, 0, time.UTC)},
		{"Large step", "y", time.Date(2024, 2, 2, 0, 0, 0, 0, time.UTC), time.Date(1024, 2, 1, 0, 0, 0, 0, time.UTC), time.Date(2025, 2, 1, 0, 0, 0, 0, time.UTC)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			yearly := yearly.New(tt.timeStep)

			newDate := yearly.Next(tt.now, tt.startDate)
			assert.Equal(t, tt.wantDate, newDate)
		})
	}
}
