package nextdate_test

import (
	"testing"
	"time"

	"github.com/10Narratives/task-tracker/internal/services/nextdate"
	"github.com/stretchr/testify/assert"
)

func TestNextDate(t *testing.T) {
	type args struct {
		now    time.Time
		date   time.Time
		repeat string
	}

	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "successful daily move",
			args: args{
				now:    time.Date(2025, 2, 5, 0, 0, 0, 0, time.UTC),
				date:   time.Date(2025, 2, 4, 0, 0, 0, 0, time.UTC),
				repeat: "d 3",
			},
			want: "20250207",
		},
	}

	for _, tc := range tests {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			res := nextdate.NextDate(tc.args.now, tc.args.now, tc.args.repeat)
			assert.Equal(t, tc.want, res)
		})
	}
}
