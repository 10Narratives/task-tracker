package nextdate_test

import (
	"testing"
	"time"

	"github.com/10Narratives/task-tracker/internal/services/nextdate"
	"github.com/stretchr/testify/assert"
)

func TestNextDate(t *testing.T) {
	t.Parallel()

	type args struct {
		now    time.Time
		date   time.Time
		repeat string
	}

	tests := []struct {
		name     string
		args     args
		wantDate string
	}{
		{
			name: "day step",
			args: args{
				now:    time.Date(2024, 2, 2, 0, 0, 0, 0, time.UTC),
				date:   time.Date(2024, 2, 1, 0, 0, 0, 0, time.UTC),
				repeat: "d 7",
			},
			wantDate: "20240208",
		},
		{
			name: "year step",
			args: args{
				now:    time.Date(2024, 2, 2, 0, 0, 0, 0, time.UTC),
				date:   time.Date(2024, 2, 1, 0, 0, 0, 0, time.UTC),
				repeat: "y",
			},
			wantDate: "20250201",
		},
		{
			name: "week step",
			args: args{
				now:    time.Date(2024, 1, 27, 0, 0, 0, 0, time.UTC),
				date:   time.Date(2024, 1, 26, 0, 0, 0, 0, time.UTC),
				repeat: "w 7",
			},
			wantDate: "20240128",
		},
		// TODO: Make test case for month
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			date := nextdate.NextDate(tt.args.now, tt.args.date, tt.args.repeat)
			assert.Equal(t, tt.wantDate, date)
		})
	}
}
