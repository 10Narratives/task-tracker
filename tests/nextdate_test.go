package tests

import (
	"testing"
	"time"

	"github.com/10Narratives/task-tracker/internal/services/nextdate/dateiters"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDailyNext(t *testing.T) {
	tests := []struct {
		name     string
		now      time.Time
		step     string
		start    time.Time
		expected time.Time
		wantErr  error
	}{
		{
			name:     "Daily step 7 days",
			start:    time.Date(2025, 2, 1, 0, 0, 0, 0, time.UTC),
			step:     "d 7",
			now:      time.Date(2025, 2, 2, 0, 0, 0, 0, time.UTC),
			expected: time.Date(2025, 2, 8, 0, 0, 0, 0, time.UTC),
			wantErr:  nil,
		},
		{
			name:     "Daily step 1 day",
			start:    time.Date(2024, 2, 1, 0, 0, 0, 0, time.UTC),
			step:     "d 1",
			now:      time.Date(2024, 2, 4, 0, 0, 0, 0, time.UTC),
			expected: time.Date(2024, 2, 4, 0, 0, 0, 0, time.UTC),
			wantErr:  nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			daily, err := dateiters.NewDaily(tt.step)
			require.NoError(t, err)

			result := daily.Next(tt.now, tt.start)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestWeeklyNext(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name     string
		now      time.Time
		step     string
		start    time.Time
		expected time.Time
		wantErr  error
	}{
		{
			name:     "Weekly step on 1",
			start:    time.Date(2025, 2, 1, 0, 0, 0, 0, time.UTC),
			step:     "w 1",
			now:      time.Date(2025, 2, 2, 0, 0, 0, 0, time.UTC),
			expected: time.Date(2025, 2, 3, 0, 0, 0, 0, time.UTC),
			wantErr:  nil,
		},
		{
			name:     "Weekly step on 2",
			start:    time.Date(2025, 2, 1, 0, 0, 0, 0, time.UTC),
			step:     "w 2",
			now:      time.Date(2025, 2, 2, 0, 0, 0, 0, time.UTC),
			expected: time.Date(2025, 2, 4, 0, 0, 0, 0, time.UTC),
			wantErr:  nil,
		},
		{
			name:     "Weekly step on 3",
			start:    time.Date(2025, 2, 1, 0, 0, 0, 0, time.UTC),
			step:     "w 3",
			now:      time.Date(2025, 2, 2, 0, 0, 0, 0, time.UTC),
			expected: time.Date(2025, 2, 5, 0, 0, 0, 0, time.UTC),
			wantErr:  nil,
		},
		{
			name:     "Weekly step on 4",
			start:    time.Date(2025, 2, 1, 0, 0, 0, 0, time.UTC),
			step:     "w 4",
			now:      time.Date(2025, 2, 2, 0, 0, 0, 0, time.UTC),
			expected: time.Date(2025, 2, 6, 0, 0, 0, 0, time.UTC),
			wantErr:  nil,
		},
		{
			name:     "Weekly step on 5",
			start:    time.Date(2025, 2, 1, 0, 0, 0, 0, time.UTC),
			step:     "w 5",
			now:      time.Date(2025, 2, 2, 0, 0, 0, 0, time.UTC),
			expected: time.Date(2025, 2, 7, 0, 0, 0, 0, time.UTC),
			wantErr:  nil,
		},
		{
			name:     "Weekly step on 6",
			start:    time.Date(2025, 2, 1, 0, 0, 0, 0, time.UTC),
			step:     "w 6",
			now:      time.Date(2025, 2, 2, 0, 0, 0, 0, time.UTC),
			expected: time.Date(2025, 2, 8, 0, 0, 0, 0, time.UTC),
			wantErr:  nil,
		},
		{
			name:     "Weekly step on 7",
			start:    time.Date(2025, 2, 1, 0, 0, 0, 0, time.UTC),
			step:     "w 7",
			now:      time.Date(2025, 2, 2, 0, 0, 0, 0, time.UTC),
			expected: time.Date(2025, 2, 2, 0, 0, 0, 0, time.UTC),
			wantErr:  nil,
		},
		{
			name:     "Weekly step on 1,6",
			start:    time.Date(2025, 2, 1, 0, 0, 0, 0, time.UTC),
			step:     "w 1,6",
			now:      time.Date(2025, 2, 2, 0, 0, 0, 0, time.UTC),
			expected: time.Date(2025, 2, 3, 0, 0, 0, 0, time.UTC),
			wantErr:  nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			daily, err := dateiters.NewWeekly(tt.step)
			require.NoError(t, err)

			result := daily.Next(tt.now, tt.start)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestYearlyNext(t *testing.T) {
	tests := []struct {
		name     string
		now      time.Time
		step     string
		start    time.Time
		expected time.Time
		wantErr  error
	}{
		{
			name:     "Many years later",
			start:    time.Date(1000, 2, 8, 0, 0, 0, 0, time.UTC),
			step:     "y",
			now:      time.Date(2025, 2, 1, 0, 0, 0, 0, time.UTC),
			expected: time.Date(2025, 2, 8, 0, 0, 0, 0, time.UTC),
			wantErr:  nil,
		},
		{
			name:     "Yearly step February-March",
			start:    time.Date(2024, 2, 29, 0, 0, 0, 0, time.UTC),
			step:     "y",
			now:      time.Date(2025, 2, 4, 0, 0, 0, 0, time.UTC),
			expected: time.Date(2025, 3, 1, 0, 0, 0, 0, time.UTC),
			wantErr:  nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			daily, err := dateiters.NewYearly(tt.step)
			require.NoError(t, err)

			result := daily.Next(tt.now, tt.start)
			assert.Equal(t, tt.expected, result)
		})
	}
}
