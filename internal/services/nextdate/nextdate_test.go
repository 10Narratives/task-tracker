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
		{
			name: "successful yearly move",
			args: args{now: time.Date(2024, 1, 26, 0, 0, 0, 0, time.UTC), date: time.Date(2023, 2, 20, 0, 0, 0, 0, time.UTC), repeat: "y"},
			want: "20240220",
		},
		{
			name: "successful yearly move",
			args: args{now: time.Date(2024, 2, 26, 0, 0, 0, 0, time.UTC), date: time.Date(2023, 1, 26, 0, 0, 0, 0, time.UTC), repeat: "y"},
			want: "20250126",
		},
		{
			name: "successful yearly move",
			args: args{now: time.Date(2024, 1, 9, 0, 0, 0, 0, time.UTC), date: time.Date(2023, 1, 10, 0, 0, 0, 0, time.UTC), repeat: "y"},
			want: "20240110",
		},
		{
			name: "successful yearly move",
			args: args{now: time.Date(2024, 1, 10, 0, 0, 0, 0, time.UTC), date: time.Date(2023, 1, 9, 0, 0, 0, 0, time.UTC), repeat: "y"},
			want: "20250109",
		},
		{
			name: "successful yearly move",
			args: args{now: time.Date(2024, 1, 9, 0, 0, 0, 0, time.UTC), date: time.Date(2023, 1, 9, 0, 0, 0, 0, time.UTC), repeat: "y"},
			want: "20250109",
		},
		{
			name: "successful weekly move",
			args: args{now: time.Date(2024, 1, 26, 0, 0, 0, 0, time.UTC), date: time.Date(2024, 1, 25, 0, 0, 0, 0, time.UTC), repeat: "w 1,2,3"},
			want: "20240129",
		},
		{
			name: "successful weekly move",
			args: args{now: time.Date(2024, 1, 26, 0, 0, 0, 0, time.UTC), date: time.Date(2024, 1, 26, 0, 0, 0, 0, time.UTC), repeat: "w 7"},
			want: "20240128",
		},
		{
			name: "successful weekly move",
			args: args{now: time.Date(2024, 1, 26, 0, 0, 0, 0, time.UTC), date: time.Date(2023, 1, 26, 0, 0, 0, 0, time.UTC), repeat: "w 4,5"},
			want: "20240201",
		},
		{
			name: "m 13",
			args: args{now: time.Date(2024, 1, 26, 0, 0, 0, 0, time.UTC), date: time.Date(2023, 11, 6, 0, 0, 0, 0, time.UTC), repeat: "m 13"},
			want: "20240213",
		},
		{
			name: "m 16,5",
			args: args{now: time.Date(2024, 1, 26, 0, 0, 0, 0, time.UTC), date: time.Date(2024, 1, 16, 0, 0, 0, 0, time.UTC), repeat: "m 16,5"},
			want: "20240205",
		},
		{
			name: "m 25,26,7",
			args: args{now: time.Date(2024, 1, 26, 0, 0, 0, 0, time.UTC), date: time.Date(2024, 1, 26, 0, 0, 0, 0, time.UTC), repeat: "m 25,26,7"},
			want: "20240207",
		}, //
		{
			name: "m 31",
			args: args{now: time.Date(2024, 1, 26, 0, 0, 0, 0, time.UTC), date: time.Date(2024, 4, 9, 0, 0, 0, 0, time.UTC), repeat: "m 31"},
			want: "20240531",
		}, //
		{
			name: "m 10,17 12,8,1",
			args: args{now: time.Date(2024, 1, 26, 0, 0, 0, 0, time.UTC), date: time.Date(2024, 3, 29, 0, 0, 0, 0, time.UTC), repeat: "m 10,17 12,8,1"},
			want: "20240810",
		},
		{
			name: "m 07,19 05,6",
			args: args{now: time.Date(2024, 1, 26, 0, 0, 0, 0, time.UTC), date: time.Date(2023, 3, 11, 0, 0, 0, 0, time.UTC), repeat: "m 07,19 05,6"},
			want: "20240507",
		},
		{
			name: "m 1 1,2",
			args: args{now: time.Date(2024, 1, 26, 0, 0, 0, 0, time.UTC), date: time.Date(2023, 3, 11, 0, 0, 0, 0, time.UTC), repeat: "m 1 1,2"},
			want: "20240201",
		},
		{
			name: "m -1",
			args: args{now: time.Date(2024, 1, 26, 0, 0, 0, 0, time.UTC), date: time.Date(2024, 1, 27, 0, 0, 0, 0, time.UTC), repeat: "m -1"},
			want: "20240131",
		},
		{
			name: "m -2",
			args: args{now: time.Date(2024, 1, 26, 0, 0, 0, 0, time.UTC), date: time.Date(2024, 2, 22, 0, 0, 0, 0, time.UTC), repeat: "m -2"},
			want: "20240228",
		}, //
		{
			name: "m -1,-2",
			args: args{now: time.Date(2024, 1, 26, 0, 0, 0, 0, time.UTC), date: time.Date(2024, 3, 26, 0, 0, 0, 0, time.UTC), repeat: "m -1,-2"},
			want: "20240330",
		}, //
		{
			name: "m -1,18",
			args: args{now: time.Date(2024, 1, 26, 0, 0, 0, 0, time.UTC), date: time.Date(2024, 2, 1, 0, 0, 0, 0, time.UTC), repeat: "m -1,18"},
			want: "20240218",
		}, //
	}

	for _, tc := range tests {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			//t.Parallel()

			res := nextdate.NextDate(tc.args.now, tc.args.date, tc.args.repeat)
			assert.Equal(t, tc.want, res)
		})
	}
}
