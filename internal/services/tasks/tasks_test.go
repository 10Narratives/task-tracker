package tasks_test

import (
	"context"
	"errors"
	"testing"

	"github.com/10Narratives/task-tracker/internal/models"
	"github.com/10Narratives/task-tracker/internal/services/tasks"
	"github.com/10Narratives/task-tracker/internal/services/tasks/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestTaskService_Register(t *testing.T) {
	var (
		id      int64           = 1
		ctx     context.Context = context.Background()
		date    string          = "20250207"
		title   string          = "Test title"
		comment string          = "Test comment"
		repeat  string          = "d 7"
	)

	type args struct {
		ctx     context.Context
		date    string
		title   string
		comment string
		repeat  string
	}

	tests := []struct {
		name       string
		mockSetup  func(m *mocks.TaskStorage)
		args       args
		wantResult require.ValueAssertionFunc
		wantErr    require.ErrorAssertionFunc
	}{
		{
			name: "successful registration",
			mockSetup: func(m *mocks.TaskStorage) {
				m.
					On("Create", ctx, date, title, comment, repeat).
					Return(id, nil)
			},
			args: args{
				ctx:     context.Background(),
				date:    date,
				title:   title,
				comment: comment,
				repeat:  repeat,
			},
			wantResult: func(tt require.TestingT, got interface{}, _ ...interface{}) {
				gotID, ok := got.(int64)
				require.True(t, ok)
				assert.Equal(t, id, gotID)
			},
			wantErr: require.NoError,
		},
		{
			name: "unsuccessful registration - database error is occurred",
			mockSetup: func(m *mocks.TaskStorage) {
				m.
					On("Create", ctx, date, title, comment, repeat).
					Return(int64(0), errors.New("database error"))
			},
			args: args{
				ctx:     context.Background(),
				date:    date,
				title:   title,
				comment: comment,
				repeat:  repeat,
			},
			wantResult: func(tt require.TestingT, got interface{}, _ ...interface{}) {
				gotID, ok := got.(int64)
				require.True(t, ok)
				assert.Equal(t, int64(0), gotID)
			},
			wantErr: func(tt require.TestingT, err error, i ...interface{}) {
				assert.EqualError(t, err, "database error")
			},
		},
	}

	for _, tc := range tests {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			storage := mocks.NewTaskStorage(t)
			tc.mockSetup(storage)

			service := tasks.New(storage)
			id, err := service.Register(tc.args.ctx, tc.args.date, tc.args.title, tc.args.comment, tc.args.repeat)
			tc.wantErr(t, err)
			tc.wantResult(t, id)

			storage.AssertExpectations(t)
		})
	}
}

func TestTaskService_Task(t *testing.T) {
	var (
		ctx     context.Context = context.Background()
		id      int64           = 1
		date    string          = "20250207"
		title   string          = "Test title"
		comment string          = "Test comment"
		repeat  string          = "d 7"
	)

	type args struct {
		ctx  context.Context
		id   int64
		task models.Task
	}

	tests := []struct {
		name       string
		mockSetup  func(m *mocks.TaskStorage)
		args       args
		wantResult require.ValueAssertionFunc
		wantErr    require.ErrorAssertionFunc
	}{
		{
			name: "successful task reading",
			mockSetup: func(m *mocks.TaskStorage) {
				m.On("Read", ctx, id).Return(models.Task{
					ID:      id,
					Date:    date,
					Title:   title,
					Comment: comment,
					Repeat:  repeat,
				}, nil)
			},
			args: args{
				ctx: ctx,
				id:  id,
				task: models.Task{
					ID:      id,
					Date:    date,
					Title:   title,
					Comment: comment,
					Repeat:  repeat,
				},
			},
			wantResult: func(tt require.TestingT, got interface{}, _ ...interface{}) {
				task, ok := got.(models.Task)
				require.True(t, ok)

				assert.Equal(t, id, task.ID)
				assert.Equal(t, date, task.Date)
				assert.Equal(t, title, task.Title)
				assert.Equal(t, comment, task.Comment)
				assert.Equal(t, repeat, task.Repeat)
			},
			wantErr: require.NoError,
		},
		{
			name: "unsuccessful reading - database error is occurred",
			mockSetup: func(m *mocks.TaskStorage) {
				m.
					On("Read", ctx, id).
					Return(models.Task{}, errors.New("database error"))
			},
			args: args{
				ctx: ctx,
				id:  id,
			},
			wantResult: func(tt require.TestingT, got interface{}, _ ...interface{}) {
				task, ok := got.(models.Task)
				require.True(t, ok)

				assert.Equal(t, int64(0), task.ID)
				assert.Equal(t, "", task.Date)
				assert.Equal(t, "", task.Title)
				assert.Equal(t, "", task.Comment)
				assert.Equal(t, "", task.Repeat)
			},
			wantErr: func(tt require.TestingT, err error, i ...interface{}) {
				assert.EqualError(t, err, "database error")
			},
		},
	}

	for _, tc := range tests {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			storage := mocks.NewTaskStorage(t)
			tc.mockSetup(storage)

			service := tasks.New(storage)
			task, err := service.Task(tc.args.ctx, tc.args.id)
			tc.wantResult(t, task)
			tc.wantErr(t, err)

			storage.AssertExpectations(t)
		})
	}
}

func TestTaskService_Tasks(t *testing.T) {
	type args struct {
		ctx    context.Context
		search string
	}

	tests := []struct {
		name       string
		mockSetup  func(m *mocks.TaskStorage)
		args       args
		wantResult require.ValueAssertionFunc
		wantErr    require.ErrorAssertionFunc
	}{
		{
			name: "successful reading of task group",
			mockSetup: func(m *mocks.TaskStorage) {
				m.
					On("ReadGroup", mock.Anything).
					Return([]models.Task{
						{ID: 1, Date: "20250206", Title: "Task 1", Comment: "Task 1 comment", Repeat: "d 7"},
						{ID: 2, Date: "20250206", Title: "Task 2", Comment: "Task 2 comment", Repeat: "d 7"},
						{ID: 3, Date: "20250206", Title: "Task 3", Comment: "Task 3 comment", Repeat: "d 7"},
						{ID: 4, Date: "20250206", Title: "Task 4", Comment: "Task 4 comment", Repeat: "d 7"},
					}, nil)
			},
			wantResult: func(tt require.TestingT, got interface{}, _ ...interface{}) {
				tasks, ok := got.([]models.Task)
				require.True(t, ok)

				assert.Len(t, tasks, 4)

				assert.Equal(t, int64(1), tasks[0].ID)
				assert.Equal(t, "20250206", tasks[0].Date)
				assert.Equal(t, "Task 1", tasks[0].Title)
				assert.Equal(t, "Task 1 comment", tasks[0].Comment)
				assert.Equal(t, "d 7", tasks[0].Repeat)

				assert.Equal(t, int64(2), tasks[1].ID)
				assert.Equal(t, "20250206", tasks[1].Date)
				assert.Equal(t, "Task 2", tasks[1].Title)
				assert.Equal(t, "Task 2 comment", tasks[1].Comment)
				assert.Equal(t, "d 7", tasks[1].Repeat)

				assert.Equal(t, int64(3), tasks[2].ID)
				assert.Equal(t, "20250206", tasks[2].Date)
				assert.Equal(t, "Task 3", tasks[2].Title)
				assert.Equal(t, "Task 3 comment", tasks[2].Comment)
				assert.Equal(t, "d 7", tasks[2].Repeat)

				assert.Equal(t, int64(4), tasks[3].ID)
				assert.Equal(t, "20250206", tasks[3].Date)
				assert.Equal(t, "Task 4", tasks[3].Title)
				assert.Equal(t, "Task 4 comment", tasks[3].Comment)
				assert.Equal(t, "d 7", tasks[3].Repeat)
			},
			wantErr: require.NoError,
		},
		{
			name: "successful reading of task group by date",
			mockSetup: func(m *mocks.TaskStorage) {
				m.
					On("ReadByDate", mock.Anything, "20250206").
					Return([]models.Task{
						{ID: 1, Date: "20250206", Title: "Task 1", Comment: "Task 1 comment", Repeat: "d 7"},
						{ID: 2, Date: "20250206", Title: "Task 2", Comment: "Task 2 comment", Repeat: "d 7"},
						{ID: 3, Date: "20250206", Title: "Task 3", Comment: "Task 3 comment", Repeat: "d 7"},
						{ID: 4, Date: "20250206", Title: "Task 4", Comment: "Task 4 comment", Repeat: "d 7"},
					}, nil)
			},
			args: args{search: "06.02.2025"},
			wantResult: func(tt require.TestingT, got interface{}, _ ...interface{}) {
				tasks, ok := got.([]models.Task)
				require.True(t, ok)

				assert.Len(t, tasks, 4)

				assert.Equal(t, int64(1), tasks[0].ID)
				assert.Equal(t, "20250206", tasks[0].Date)
				assert.Equal(t, "Task 1", tasks[0].Title)
				assert.Equal(t, "Task 1 comment", tasks[0].Comment)
				assert.Equal(t, "d 7", tasks[0].Repeat)

				assert.Equal(t, int64(2), tasks[1].ID)
				assert.Equal(t, "20250206", tasks[1].Date)
				assert.Equal(t, "Task 2", tasks[1].Title)
				assert.Equal(t, "Task 2 comment", tasks[1].Comment)
				assert.Equal(t, "d 7", tasks[1].Repeat)

				assert.Equal(t, int64(3), tasks[2].ID)
				assert.Equal(t, "20250206", tasks[2].Date)
				assert.Equal(t, "Task 3", tasks[2].Title)
				assert.Equal(t, "Task 3 comment", tasks[2].Comment)
				assert.Equal(t, "d 7", tasks[2].Repeat)

				assert.Equal(t, int64(4), tasks[3].ID)
				assert.Equal(t, "20250206", tasks[3].Date)
				assert.Equal(t, "Task 4", tasks[3].Title)
				assert.Equal(t, "Task 4 comment", tasks[3].Comment)
				assert.Equal(t, "d 7", tasks[3].Repeat)
			},
			wantErr: require.NoError,
		},
		{
			name: "successful reading of task group by payload",
			mockSetup: func(m *mocks.TaskStorage) {
				m.
					On("ReadByPayload", mock.Anything, "Task").
					Return([]models.Task{
						{ID: 1, Date: "20250206", Title: "Task 1", Comment: "Task 1 comment", Repeat: "d 7"},
						{ID: 2, Date: "20250206", Title: "Task 2", Comment: "Task 2 comment", Repeat: "d 7"},
						{ID: 3, Date: "20250206", Title: "Task 3", Comment: "Task 3 comment", Repeat: "d 7"},
						{ID: 4, Date: "20250206", Title: "Task 4", Comment: "Task 4 comment", Repeat: "d 7"},
					}, nil)
			},
			args: args{search: "Task"},
			wantResult: func(tt require.TestingT, got interface{}, _ ...interface{}) {
				tasks, ok := got.([]models.Task)
				require.True(t, ok)

				assert.Len(t, tasks, 4)

				assert.Equal(t, int64(1), tasks[0].ID)
				assert.Equal(t, "20250206", tasks[0].Date)
				assert.Equal(t, "Task 1", tasks[0].Title)
				assert.Equal(t, "Task 1 comment", tasks[0].Comment)
				assert.Equal(t, "d 7", tasks[0].Repeat)

				assert.Equal(t, int64(2), tasks[1].ID)
				assert.Equal(t, "20250206", tasks[1].Date)
				assert.Equal(t, "Task 2", tasks[1].Title)
				assert.Equal(t, "Task 2 comment", tasks[1].Comment)
				assert.Equal(t, "d 7", tasks[1].Repeat)

				assert.Equal(t, int64(3), tasks[2].ID)
				assert.Equal(t, "20250206", tasks[2].Date)
				assert.Equal(t, "Task 3", tasks[2].Title)
				assert.Equal(t, "Task 3 comment", tasks[2].Comment)
				assert.Equal(t, "d 7", tasks[2].Repeat)

				assert.Equal(t, int64(4), tasks[3].ID)
				assert.Equal(t, "20250206", tasks[3].Date)
				assert.Equal(t, "Task 4", tasks[3].Title)
				assert.Equal(t, "Task 4 comment", tasks[3].Comment)
				assert.Equal(t, "d 7", tasks[3].Repeat)
			},
			wantErr: require.NoError,
		},
	}

	for _, tc := range tests {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			storage := mocks.NewTaskStorage(t)
			tc.mockSetup(storage)

			service := tasks.New(storage)
			tasks, err := service.Tasks(tc.args.ctx, tc.args.search)
			tc.wantResult(t, tasks)
			tc.wantErr(t, err)

			storage.AssertExpectations(t)
		})
	}
}

func TestTaskService_Delete(t *testing.T) {
	var (
		ctx       = context.Background()
		id  int64 = 100
	)

	type args struct {
		ctx context.Context
		id  int64
	}

	tests := []struct {
		name      string
		mockSetup func(m *mocks.TaskStorage)
		args      args
		// wantResult require.ValueAssertionFunc
		wantErr require.ErrorAssertionFunc
	}{
		{
			name: "successful deletion",
			mockSetup: func(m *mocks.TaskStorage) {
				m.On("Delete", ctx, id).Return(nil)
			},
			args:    args{ctx: ctx, id: id},
			wantErr: require.NoError,
		},
		{
			name: "unsuccessful deletion - database error is occurred",
			mockSetup: func(m *mocks.TaskStorage) {
				m.On("Delete", ctx, id).Return(errors.New("database error"))
			},
			args: args{ctx: ctx, id: id},
			wantErr: func(tt require.TestingT, err error, i ...interface{}) {
				assert.EqualError(t, err, "database error")
			},
		},
	}

	for _, tc := range tests {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			storage := mocks.NewTaskStorage(t)
			tc.mockSetup(storage)

			service := tasks.New(storage)
			err := service.Delete(tc.args.ctx, tc.args.id)
			//	tc.wantResult(t, tasks)
			tc.wantErr(t, err)

			storage.AssertExpectations(t)
		})
	}
}

func TestTaskService_Update(t *testing.T) {
	var (
		ctx           = context.Background()
		id      int64 = 100
		date          = "20250205"
		title         = "Title"
		comment       = "Comment"
		repeat        = "d 100"
	)

	type args struct {
		ctx  context.Context
		task *models.Task
	}

	tests := []struct {
		name      string
		mockSetup func(m *mocks.TaskStorage)
		args      args
		wantErr   require.ErrorAssertionFunc
	}{
		{
			name: "successful update",
			mockSetup: func(m *mocks.TaskStorage) {
				m.On("Update", ctx, &models.Task{ID: id, Date: date, Title: title, Comment: comment, Repeat: repeat}).Return(nil)
			},
			args:    args{ctx: ctx, task: &models.Task{ID: id, Date: date, Title: title, Comment: comment, Repeat: repeat}},
			wantErr: require.NoError,
		},
		{
			name: "unsuccessful update - database error is occurred",
			mockSetup: func(m *mocks.TaskStorage) {
				m.
					On("Update", ctx, &models.Task{ID: id, Date: date, Title: title, Comment: comment, Repeat: repeat}).
					Return(errors.New("database error"))
			},
			args: args{ctx: ctx, task: &models.Task{ID: id, Date: date, Title: title, Comment: comment, Repeat: repeat}},
			wantErr: func(tt require.TestingT, err error, i ...interface{}) {
				assert.EqualError(t, err, "database error")
			},
		},
	}

	for _, tc := range tests {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			storage := mocks.NewTaskStorage(t)
			tc.mockSetup(storage)

			service := tasks.New(storage)
			err := service.Update(tc.args.ctx, id, date, title, comment, repeat)
			tc.wantErr(t, err)

			storage.AssertExpectations(t)
		})
	}
}

func TestTaskService_Complete(t *testing.T) {
	type args struct {
		ctx context.Context
		id  int64
	}

	tests := []struct {
		name      string
		mockSetup func(m *mocks.TaskStorage)
		args      args
		wantErr   require.ErrorAssertionFunc
	}{
		{
			name: "successful complete - with nextdate",
			mockSetup: func(m *mocks.TaskStorage) {
				m.
					On("Read", mock.Anything, int64(100)).
					Return(models.Task{ID: 100, Date: "20250402", Title: "Title", Comment: "Comment", Repeat: "d 7"}, nil)
				m.
					On("Update", mock.Anything, mock.Anything).
					Return(nil)
			},
			args:    args{ctx: context.Background(), id: 100},
			wantErr: require.NoError,
		},
		{
			name: "unsuccessful complete - with nextdate",
			mockSetup: func(m *mocks.TaskStorage) {
				m.
					On("Read", mock.Anything, int64(100)).
					Return(models.Task{}, errors.New("database error"))
			},
			args: args{ctx: context.Background(), id: 100},
			wantErr: func(tt require.TestingT, err error, i ...interface{}) {
				assert.EqualError(t, err, "database error")
			},
		},
		{
			name: "successful complete - without nextdate",
			mockSetup: func(m *mocks.TaskStorage) {
				m.
					On("Read", mock.Anything, int64(100)).
					Return(models.Task{ID: 100, Date: "20250402", Title: "Title", Comment: "Comment", Repeat: ""}, nil)
				m.
					On("Delete", mock.Anything, int64(100)).
					Return(nil)
			},
			args:    args{ctx: context.Background(), id: 100},
			wantErr: require.NoError,
		},
		{
			name: "unsuccessful complete - without nextdate",
			mockSetup: func(m *mocks.TaskStorage) {
				m.
					On("Read", mock.Anything, int64(100)).
					Return(models.Task{ID: 100, Date: "20250402", Title: "Title", Comment: "Comment", Repeat: ""}, nil)
				m.
					On("Delete", mock.Anything, int64(100)).
					Return(errors.New("database error"))
			},
			args: args{ctx: context.Background(), id: 100},
			wantErr: func(tt require.TestingT, err error, i ...interface{}) {
				assert.EqualError(t, err, "database error")
			},
		},
	}

	for _, tc := range tests {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			storage := mocks.NewTaskStorage(t)
			tc.mockSetup(storage)

			service := tasks.New(storage)
			err := service.Complete(tc.args.ctx, tc.args.id)
			tc.wantErr(t, err)

			storage.AssertExpectations(t)
		})
	}
}
