package tests

import (
	"testing"
	"time"

	"github.com/10Narratives/task-tracker/internal/services/nextdate"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

/*
	{"20240113", "d 7", `20240127`},
	{"20240120", "d 20", `20240209`},
	{"20240202", "d 30", `20240303`},
	{"20240320", "d 401", ""},
	{"20231225", "d 12", `20240130`},
	{"20240228", "d 1", "20240229"},*/

func TestNextDate(t *testing.T) {
	tests := []struct {
		name      string
		now       time.Time
		date      string
		timeStep  string
		want      string
		expectErr error
	}{
		{
			name:     "Valid daily step 7 days",
			now:      time.Date(2025, 2, 1, 0, 0, 0, 0, time.UTC),
			date:     "20240201",
			timeStep: "d 7",
			want:     "20250206",
		},
		{
			name:     "Valid daily step 20 days",
			now:      time.Date(2024, 2, 1, 0, 0, 0, 0, time.UTC),
			date:     "20240120",
			timeStep: "d 20",
			want:     "20240209",
		},
		{
			name:     "Valid daily step 30 days",
			now:      time.Date(2024, 2, 3, 0, 0, 0, 0, time.UTC),
			date:     "20240202",
			timeStep: "d 30",
			want:     "20240303",
		},
		{
			name:      "Invalid date format",
			now:       time.Now(),
			date:      "invalid",
			timeStep:  "d 7",
			expectErr: nextdate.ErrCanNotParseStartDate,
		},
		{
			name:      "Invalid time step",
			now:       time.Now(),
			date:      "20240201",
			timeStep:  "d 0",
			expectErr: nextdate.ErrTimeStepIsNotValid,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := nextdate.NextDate(tt.now, tt.date, tt.timeStep)

			if tt.expectErr != nil {
				require.Error(t, err)
				assert.ErrorIs(t, err, tt.expectErr)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

func TestDailyNext(t *testing.T) {
	tests := []struct {
		name     string
		start    time.Time
		step     string
		now      time.Time
		expected time.Time
	}{
		{
			name:     "Daily step 7 days",
			start:    time.Date(2024, 2, 1, 0, 0, 0, 0, time.UTC),
			step:     "d 7",
			now:      time.Date(2025, 2, 1, 0, 0, 0, 0, time.UTC),
			expected: time.Date(2025, 2, 6, 0, 0, 0, 0, time.UTC),
		},
		{
			name:     "Daily step 1 day",
			start:    time.Date(2024, 2, 1, 0, 0, 0, 0, time.UTC),
			step:     "d 1",
			now:      time.Date(2024, 2, 5, 0, 0, 0, 0, time.UTC),
			expected: time.Date(2024, 2, 5, 0, 0, 0, 0, time.UTC),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			daily, err := nextdate.NewDaily(tt.step)
			require.NoError(t, err)

			result := daily.Next(tt.now, tt.start)
			assert.Equal(t, tt.expected, result)
		})
	}
}
