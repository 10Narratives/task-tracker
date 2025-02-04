package nextdate_test

import (
	"testing"
	"time"

	"github.com/10Narratives/task-tracker/internal/services/nextdate"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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
		wantErr  require.ErrorAssertionFunc
	}{
		{
			name: "day step",
			args: args{
				now:    time.Date(2024, 2, 2, 0, 0, 0, 0, time.UTC),
				date:   time.Date(2024, 2, 1, 0, 0, 0, 0, time.UTC),
				repeat: "d 7",
			},
			wantDate: "20240208",
			wantErr:  require.NoError,
		},
		{
			name: "year step",
			args: args{
				now:    time.Date(2024, 2, 2, 0, 0, 0, 0, time.UTC),
				date:   time.Date(2024, 2, 1, 0, 0, 0, 0, time.UTC),
				repeat: "y",
			},
			wantDate: "20250201",
			wantErr:  require.NoError,
		},
		{
			name: "week step",
			args: args{
				now:    time.Date(2024, 1, 27, 0, 0, 0, 0, time.UTC),
				date:   time.Date(2024, 1, 26, 0, 0, 0, 0, time.UTC),
				repeat: "w 7",
			},
			wantDate: "20240128",
			wantErr:  require.NoError,
		},
		{
			name: "empty repeat",
			args: args{
				now:    time.Date(2024, 1, 27, 0, 0, 0, 0, time.UTC),
				date:   time.Date(2024, 1, 26, 0, 0, 0, 0, time.UTC),
				repeat: "",
			},
			wantDate: "",
			wantErr: func(tt require.TestingT, err error, i ...interface{}) {
				assert.EqualError(t, err, nextdate.ErrEmptyRepeat.Error())
			},
		},
		{
			name: "unsupported option",
			args: args{
				now:    time.Date(2024, 1, 27, 0, 0, 0, 0, time.UTC),
				date:   time.Date(2024, 1, 26, 0, 0, 0, 0, time.UTC),
				repeat: "l 500",
			},
			wantDate: "",
			wantErr: func(tt require.TestingT, err error, i ...interface{}) {
				assert.EqualError(t, err, nextdate.ErrUnsupportedOption.Error())
			},
		},
		// TODO: Make test case for month
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			date, err := nextdate.NextDate(tt.args.now, tt.args.date, tt.args.repeat)
			assert.Equal(t, tt.wantDate, date)
			tt.wantErr(t, err)
		})
	}
}
