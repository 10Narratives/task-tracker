package nextdate_test

import (
	"testing"
	"time"

	"github.com/10Narratives/task-tracker/internal/services/nextdate"
	"github.com/stretchr/testify/assert"
)

/*
	{"20240113", "d 7", `20240127`},
	{"20240120", "d 20", `20240209`},
	{"20240202", "d 30", `20240303`},
	{"20240320", "d 401", ""},
	{"20231225", "d 12", `20240130`},
	{"20240228", "d 1", "20240229"},
*/

// 20240126

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
			name: "successful daily move - d 3",
			args: args{now: time.Date(2025, 2, 5, 0, 0, 0, 0, time.UTC), date: time.Date(2025, 2, 4, 0, 0, 0, 0, time.UTC), repeat: "d 3"},
			want: "20250207",
		},
		{
			name: "successful daily move - d 7",
			args: args{now: time.Date(2025, 2, 5, 0, 0, 0, 0, time.UTC), date: time.Date(2025, 2, 4, 0, 0, 0, 0, time.UTC), repeat: "d 7"},
			want: "20250211",
		},
		{
			name: "successful daily move - d 20",
			args: args{now: time.Date(2024, 1, 26, 0, 0, 0, 0, time.UTC), date: time.Date(2024, 1, 20, 0, 0, 0, 0, time.UTC), repeat: "d 20"},
			want: "20240209",
		},
		{
			name: "successful daily move - d 30",
			args: args{now: time.Date(2024, 1, 26, 0, 0, 0, 0, time.UTC), date: time.Date(2024, 2, 2, 0, 0, 0, 0, time.UTC), repeat: "d 30"},
			want: "20240303",
		},
		{
			name: "successful daily move - same day",
			args: args{now: time.Date(2024, 1, 26, 0, 0, 0, 0, time.UTC), date: time.Date(2024, 1, 26, 0, 0, 0, 0, time.UTC), repeat: "d 1"},
			want: "20240127",
		},
		{
			name: "successful daily move - future",
			args: args{now: time.Date(2024, 1, 26, 0, 0, 0, 0, time.UTC), date: time.Date(2024, 1, 27, 0, 0, 0, 0, time.UTC), repeat: "d 1"},
			want: "20240128",
		},
	}

	for _, tc := range tests {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			res := nextdate.NextDate(tc.args.now, tc.args.date, tc.args.repeat)
			assert.Equal(t, tc.want, res)
		})
	}
}
