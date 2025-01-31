package tests

import (
	"testing"
	"time"

	"github.com/10Narratives/task-tracker/internal/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDayRepeat(t *testing.T) {
	t.Parallel()
	now := time.Date(2024, 1, 26, 12, 0, 0, 0, time.UTC)

	tests := []struct {
		name      string
		now       time.Time
		date      string
		repeat    string
		want      string
		expectErr bool
	}{
		{"Next Date 1", now, "20240113", "d 7", "20240127", false},
		{"Next Date 2", now, "20240120", "d 20", "20240209", false},
		{"Next Date 4", now, "20231225", "d 12", "20240130", false},
		{"Invalid repeat", now, "20240113", "d", "", true},
		{"Invalid repeat 2", now, "20240113", "d 555", "", true},
		{"Empty repeat", now, "20240113", "", "", true},
		{"Empty date", now, "", "d 7", "", true},
		{"Parse date error", now, "bad so bad", "d 7", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, err := service.NextDate(tt.now, tt.date, tt.repeat)

			if tt.expectErr {
				require.Error(t, err, "An error was expected, but there was none")
			} else {
				require.NoError(t, err, "The error was not expected, but it appeared")
				assert.Equal(t, tt.want, got, "The result is not as expected")
			}
		})
	}
}

/*
	{"16890220", "y", `20240220`},
	{"20250701", "y", `20260701`},
	{"20240101", "y", `20250101`},
	{"20231231", "y", `20241231`},
	{"20240229", "y", `20250301`},
	{"20240301", "y", `20250301`},
*/

func TestYearRepeat(t *testing.T) {
	t.Parallel()
	now := time.Date(2024, 1, 26, 12, 0, 0, 0, time.UTC)

	tests := []struct {
		name      string
		now       time.Time
		date      string
		repeat    string
		want      string
		expectErr bool
	}{
		{"Next Date 1", now, "16890220", "y", "20240220", false},
		{"Next Date 2", now, "20240101", "y", "20250101", false},
		{"Empty repeat", now, "20240113", "", "", true},
		{"Empty date", now, "", "y", "", true},
		{"Parse date error", now, "bad so bad", "y", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, err := service.NextDate(tt.now, tt.date, tt.repeat)

			if tt.expectErr {
				require.Error(t, err, "An error was expected, but there was none")
			} else {
				require.NoError(t, err, "The error was not expected, but it appeared")
				assert.Equal(t, tt.want, got, "The result is not as expected")
			}
		})
	}
}
