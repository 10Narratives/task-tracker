package sqlite_test

import (
	"context"
	"database/sql"
	"errors"
	"regexp"
	"testing"

	"github.com/10Narratives/task-tracker/internal/models"
	"github.com/10Narratives/task-tracker/internal/storage"
	"github.com/10Narratives/task-tracker/internal/storage/sqlite"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTaskStorage_Prepare(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		mock    func(dbMock sqlmock.Sqlmock)
		wantErr require.ErrorAssertionFunc
	}{
		{
			name: "success",
			mock: func(dbMock sqlmock.Sqlmock) {
				dbMock.ExpectPrepare(`CREATE TABLE IF NOT EXISTS scheduler`).
					WillReturnError(nil) // No error in preparing statement

				dbMock.ExpectExec(`CREATE TABLE IF NOT EXISTS scheduler`).
					WillReturnResult(sqlmock.NewResult(0, 0)) // Simulating successful execution
			},
			wantErr: require.NoError,
		},
		{
			name: "prepare statement error",
			mock: func(dbMock sqlmock.Sqlmock) {
				dbMock.ExpectPrepare(`CREATE TABLE IF NOT EXISTS scheduler`).
					WillReturnError(errors.New("prepare statement error"))
			},
			wantErr: func(tt require.TestingT, err error, i ...interface{}) {
				require.EqualError(tt, err, "can not prepare statement: prepare statement error", i...)
			},
		},
		{
			name: "exec statement error",
			mock: func(dbMock sqlmock.Sqlmock) {
				dbMock.ExpectPrepare(`CREATE TABLE IF NOT EXISTS scheduler`).
					WillReturnError(nil) // No error in preparing statement

				dbMock.ExpectExec(`CREATE TABLE IF NOT EXISTS scheduler`).
					WillReturnError(errors.New("exec statement error")) // Error when executing
			},
			wantErr: func(tt require.TestingT, err error, i ...interface{}) {
				require.EqualError(tt, err, "can not prepare database: exec statement error", i...)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			db, dbMock, err := sqlmock.New()
			require.NoError(t, err)
			defer db.Close()

			storage := sqlite.New(db, 10)
			tt.mock(dbMock)

			err = storage.Prepare()
			tt.wantErr(t, err)

			require.NoError(t, dbMock.ExpectationsWereMet())
		})
	}
}

func TestTaskStorage_Create(t *testing.T) {
	t.Parallel()

	type args struct {
		ctx  context.Context
		task *models.Task
	}

	var (
		id      int64  = 100
		date    string = "20240203"
		title   string = "test title"
		comment string = "test comment"
		repeat  string = "test repeat"
	)

	tests := []struct {
		name     string
		mocks    func(dbMock sqlmock.Sqlmock)
		args     args
		wantTask require.ValueAssertionFunc
		wantErr  require.ErrorAssertionFunc
	}{
		{
			name: "success",
			mocks: func(dbMock sqlmock.Sqlmock) {
				dbMock.
					ExpectExec(`INSERT INTO scheduler`).
					WithArgs(date, title, comment, repeat).
					WillReturnResult(sqlmock.NewResult(id, 1))
			},
			args: args{
				ctx: context.Background(),
				task: &models.Task{
					Date:    date,
					Title:   title,
					Comment: comment,
					Repeat:  repeat,
				},
			},
			wantTask: func(tt require.TestingT, got interface{}, i ...interface{}) {
				task, ok := got.(*models.Task)
				require.True(t, ok)
				require.NotNil(t, task, i...)
				assert.Equal(t, id, task.ID, i...)
				assert.Equal(t, date, task.Date, i...)
				assert.Equal(t, title, task.Title, i...)
				assert.Equal(t, comment, task.Comment, i...)
				assert.Equal(t, repeat, task.Repeat, i...)
			},
			wantErr: require.NoError,
		},
		{
			name: "database error",
			mocks: func(dbMock sqlmock.Sqlmock) {
				dbMock.
					ExpectExec(`INSERT INTO scheduler`).
					WithArgs(date, title, comment, repeat).
					WillReturnError(errors.New("database error"))
			},
			args: args{
				ctx: context.Background(),
				task: &models.Task{
					Date:    date,
					Title:   title,
					Comment: comment,
					Repeat:  repeat,
				},
			},
			wantTask: func(tt require.TestingT, got interface{}, i ...interface{}) {
				task, ok := got.(*models.Task)
				require.True(t, ok)
				require.NotNil(t, task, i...)
				assert.Equal(t, int64(0), task.ID, i...)
				assert.Equal(t, date, task.Date, i...)
				assert.Equal(t, title, task.Title, i...)
				assert.Equal(t, comment, task.Comment, i...)
				assert.Equal(t, repeat, task.Repeat, i...)
			},
			wantErr: func(tt require.TestingT, err error, i ...interface{}) {
				assert.EqualError(t, err, "cannot insert task into database: database error", i...)
			},
		},
		{
			name:  "nil task",
			mocks: func(dbMock sqlmock.Sqlmock) {},
			args: args{
				ctx:  context.Background(),
				task: nil,
			},
			wantTask: require.Nil,
			wantErr: func(tt require.TestingT, err error, i ...interface{}) {
				require.EqualError(t, err, storage.ErrNilTaskCreation.Error(), i...)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			db, dbMock, err := sqlmock.New()
			require.NoError(t, err)

			storage := sqlite.New(db, 10)
			tt.mocks(dbMock)

			err = storage.Create(tt.args.ctx, tt.args.task)
			tt.wantErr(t, err)
			tt.wantTask(t, tt.args.task)

			require.NoError(t, dbMock.ExpectationsWereMet())
		})
	}
}

func TestTaskStorage_Read(t *testing.T) {
	t.Parallel()

	type args struct {
		ctx context.Context
		id  int64
	}

	var (
		id      int64  = 100
		date    string = "20240203"
		title   string = "test title"
		comment string = "test comment"
		repeat  string = "test repeat"
	)

	tests := []struct {
		name     string
		mocks    func(dbMock sqlmock.Sqlmock)
		args     args
		wantTask require.ValueAssertionFunc
		wantErr  require.ErrorAssertionFunc
	}{
		{
			name: "successful reading",
			mocks: func(dbMock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "date", "title", "comment", "repeat"}).
					AddRow(id, date, title, comment, repeat)
				dbMock.ExpectQuery(`SELECT id, date, title, comment, repeat FROM scheduler WHERE id = ?`).
					WithArgs(id).WillReturnRows(rows)
			},
			args: args{
				ctx: context.Background(),
				id:  id,
			},
			wantTask: func(tt require.TestingT, got interface{}, i ...interface{}) {
				task, ok := got.(models.Task)
				require.True(t, ok)
				assert.Equal(t, id, task.ID, i...)
				assert.Equal(t, date, task.Date, i...)
				assert.Equal(t, title, task.Title, i...)
				assert.Equal(t, comment, task.Comment, i...)
				assert.Equal(t, repeat, task.Repeat, i...)
			},
			wantErr: require.NoError,
		},
		{
			name: "no rows",
			mocks: func(dbMock sqlmock.Sqlmock) {
				dbMock.ExpectQuery(`SELECT id, date, title, comment, repeat FROM scheduler WHERE id = ?`).
					WithArgs(id).WillReturnError(sql.ErrNoRows)
			},
			args: args{
				ctx: context.Background(),
				id:  id,
			},
			wantTask: func(tt require.TestingT, got interface{}, i ...interface{}) {
				task, ok := got.(models.Task)
				require.True(t, ok)
				require.Equal(t, models.Task{}, task)
			},
			wantErr: require.NoError,
		},
		{
			name: "database error",
			mocks: func(dbMock sqlmock.Sqlmock) {
				dbMock.
					ExpectQuery(`SELECT id, date, title, comment, repeat FROM scheduler WHERE id = ?`).
					WithArgs(id).
					WillReturnError(errors.New("database error"))
			},
			args: args{
				ctx: context.Background(),
				id:  id,
			},
			wantTask: func(tt require.TestingT, got interface{}, i ...interface{}) {
				task, ok := got.(models.Task)
				require.True(t, ok)
				require.Equal(t, models.Task{}, task)
			},
			wantErr: func(tt require.TestingT, err error, i ...interface{}) {
				assert.EqualError(t, err, "cannot read task from database: database error", i...)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			db, dbMock, err := sqlmock.New()
			require.NoError(t, err)

			storage := sqlite.New(db, 10)
			tt.mocks(dbMock)

			task, err := storage.Read(tt.args.ctx, tt.args.id)
			tt.wantErr(t, err)
			tt.wantTask(t, task)

			require.NoError(t, dbMock.ExpectationsWereMet())
		})
	}
}

func TestTaskStorage_ReadGroup(t *testing.T) {
	t.Parallel()

	type args struct {
		ctx context.Context
	}

	tests := []struct {
		name      string
		mocks     func(dbMock sqlmock.Sqlmock)
		args      args
		wantTasks require.ValueAssertionFunc
		wantErr   require.ErrorAssertionFunc
	}{
		{
			name: "successful reading",
			mocks: func(dbMock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "date", "title", "comment", "repeat"}).
					AddRow(1, "20240203", "Test title task 1", "Comment for task 1", "d 7").
					AddRow(2, "20240203", "Test title task 2", "Comment for task 2", "d 7").
					AddRow(3, "20240203", "Test title task 3", "Comment for task 3", "d 7")
				dbMock.ExpectQuery(`SELECT id, date, title, comment, repeat FROM scheduler ORDER BY date LIMIT ?`).
					WithArgs(3).
					WillReturnRows(rows)
			},
			args: args{
				ctx: context.Background(),
			},
			wantTasks: func(tt require.TestingT, got interface{}, i ...interface{}) {
				tasks, ok := got.([]models.Task)
				require.True(t, ok)
				assert.Len(t, tasks, 3)

				assert.Equal(t, int64(1), tasks[0].ID, i...)
				assert.Equal(t, "20240203", tasks[0].Date, i...)
				assert.Equal(t, "Test title task 1", tasks[0].Title, i...)
				assert.Equal(t, "Comment for task 1", tasks[0].Comment, i...)
				assert.Equal(t, "d 7", tasks[0].Repeat, i...)

				assert.Equal(t, int64(2), tasks[1].ID, i...)
				assert.Equal(t, "20240203", tasks[1].Date, i...)
				assert.Equal(t, "Test title task 2", tasks[1].Title, i...)
				assert.Equal(t, "Comment for task 2", tasks[1].Comment, i...)
				assert.Equal(t, "d 7", tasks[1].Repeat, i...)

				assert.Equal(t, int64(3), tasks[2].ID, i...)
				assert.Equal(t, "20240203", tasks[2].Date, i...)
				assert.Equal(t, "Test title task 3", tasks[2].Title, i...)
				assert.Equal(t, "Comment for task 3", tasks[2].Comment, i...)
				assert.Equal(t, "d 7", tasks[2].Repeat, i...)
			},
			wantErr: require.NoError,
		},
		{
			name: "no rows",
			mocks: func(dbMock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "date", "title", "comment", "repeat"})
				dbMock.ExpectQuery(`SELECT id, date, title, comment, repeat FROM scheduler ORDER BY date LIMIT ?`).
					WithArgs(3).
					WillReturnRows(rows)
			},
			args: args{
				ctx: context.Background(),
			},
			wantTasks: func(tt require.TestingT, got interface{}, i ...interface{}) {
				tasks, ok := got.([]models.Task)
				require.True(t, ok)
				require.Empty(t, tasks)
			},
			wantErr: require.NoError,
		},
		{
			name: "database error",
			mocks: func(dbMock sqlmock.Sqlmock) {
				dbMock.ExpectQuery(`SELECT id, date, title, comment, repeat FROM scheduler ORDER BY date LIMIT ?`).
					WithArgs(3).
					WillReturnError(errors.New("database error"))
			},
			args: args{
				ctx: context.Background(),
			},
			wantTasks: func(tt require.TestingT, got interface{}, i ...interface{}) {
				tasks, ok := got.([]models.Task)
				require.True(t, ok)
				require.Empty(t, tasks)
			},
			wantErr: func(tt require.TestingT, err error, i ...interface{}) {
				require.EqualError(tt, err, "cannot execute query: database error")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			db, dbMock, err := sqlmock.New()
			require.NoError(t, err)

			storage := sqlite.New(db, 3)
			tt.mocks(dbMock)

			tasks, err := storage.ReadGroup(tt.args.ctx)
			tt.wantErr(t, err)
			tt.wantTasks(t, tasks)

			require.NoError(t, dbMock.ExpectationsWereMet())
		})
	}
}

func TestTaskStorage_ReadByDate(t *testing.T) {
	t.Parallel()

	type args struct {
		ctx  context.Context
		date string
	}

	var date string = "20240203"

	tests := []struct {
		name      string
		mocks     func(dbMock sqlmock.Sqlmock)
		args      args
		wantTasks require.ValueAssertionFunc
		wantErr   require.ErrorAssertionFunc
	}{
		{
			name: "successful reading",
			mocks: func(dbMock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "date", "title", "comment", "repeat"}).
					AddRow(1, "20240203", "Test title task 1", "Comment for task 1", "d 7").
					AddRow(2, "20240203", "Test title task 2", "Comment for task 2", "d 7").
					AddRow(3, "20240203", "Test title task 3", "Comment for task 3", "d 7")
				query := regexp.QuoteMeta("SELECT id, date, title, comment, repeat FROM scheduler WHERE date = ? LIMIT ?")
				dbMock.ExpectQuery(query).
					WithArgs(date, 3).
					WillReturnRows(rows)
			}, // SELECT id, date, title, comment, repeat FROM scheduler WHERE date = ? LIMIT ?
			args: args{
				ctx:  context.Background(),
				date: date,
			},
			wantTasks: func(tt require.TestingT, got interface{}, i ...interface{}) {
				tasks, ok := got.([]models.Task)
				require.True(t, ok)
				assert.Len(t, tasks, 3)

				assert.Equal(t, int64(1), tasks[0].ID, i...)
				assert.Equal(t, "20240203", tasks[0].Date, i...)
				assert.Equal(t, "Test title task 1", tasks[0].Title, i...)
				assert.Equal(t, "Comment for task 1", tasks[0].Comment, i...)
				assert.Equal(t, "d 7", tasks[0].Repeat, i...)

				assert.Equal(t, int64(2), tasks[1].ID, i...)
				assert.Equal(t, "20240203", tasks[1].Date, i...)
				assert.Equal(t, "Test title task 2", tasks[1].Title, i...)
				assert.Equal(t, "Comment for task 2", tasks[1].Comment, i...)
				assert.Equal(t, "d 7", tasks[1].Repeat, i...)

				assert.Equal(t, int64(3), tasks[2].ID, i...)
				assert.Equal(t, "20240203", tasks[2].Date, i...)
				assert.Equal(t, "Test title task 3", tasks[2].Title, i...)
				assert.Equal(t, "Comment for task 3", tasks[2].Comment, i...)
				assert.Equal(t, "d 7", tasks[2].Repeat, i...)
			},
			wantErr: require.NoError,
		},
		{
			name: "no rows",
			mocks: func(dbMock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "date", "title", "comment", "repeat"})
				query := regexp.QuoteMeta("SELECT id, date, title, comment, repeat FROM scheduler WHERE date = ? LIMIT ?")
				dbMock.ExpectQuery(query).
					WithArgs(date, 3).
					WillReturnRows(rows)
			},
			args: args{
				ctx:  context.Background(),
				date: date,
			},
			wantTasks: func(tt require.TestingT, got interface{}, i ...interface{}) {
				tasks, ok := got.([]models.Task)
				require.True(t, ok)
				require.Empty(t, tasks)
			},
			wantErr: require.NoError,
		},
		{
			name: "database error",
			mocks: func(dbMock sqlmock.Sqlmock) {
				query := regexp.QuoteMeta("SELECT id, date, title, comment, repeat FROM scheduler WHERE date = ? LIMIT ?")
				dbMock.ExpectQuery(query).
					WithArgs(date, 3).
					WillReturnError(errors.New("database error"))
			},
			args: args{
				ctx:  context.Background(),
				date: date,
			},
			wantTasks: func(tt require.TestingT, got interface{}, i ...interface{}) {
				tasks, ok := got.([]models.Task)
				require.True(t, ok)
				require.Empty(t, tasks)
			},
			wantErr: func(tt require.TestingT, err error, i ...interface{}) {
				require.EqualError(tt, err, "cannot execute query: database error")
			},
		},
		{
			name:  "empty date",
			mocks: func(dbMock sqlmock.Sqlmock) {},
			args: args{
				ctx:  context.Background(),
				date: "",
			},
			wantTasks: func(tt require.TestingT, got interface{}, i ...interface{}) {
				tasks, ok := got.([]models.Task)
				require.True(t, ok)
				require.Empty(t, tasks)
			},
			wantErr: func(tt require.TestingT, err error, i ...interface{}) {
				require.EqualError(tt, err, storage.ErrEmptyDate.Error())
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			db, dbMock, err := sqlmock.New()
			require.NoError(t, err)

			storage := sqlite.New(db, 3)
			tt.mocks(dbMock)

			tasks, err := storage.ReadByDate(tt.args.ctx, tt.args.date)
			tt.wantErr(t, err)
			tt.wantTasks(t, tasks)

			require.NoError(t, dbMock.ExpectationsWereMet())
		})
	}
}

func TestTaskStorage_ReadByPayload(t *testing.T) {
	t.Parallel()

	type args struct {
		ctx     context.Context
		payload string
	}

	var payload string = "task"

	tests := []struct {
		name      string
		mocks     func(dbMock sqlmock.Sqlmock)
		args      args
		wantTasks require.ValueAssertionFunc
		wantErr   require.ErrorAssertionFunc
	}{
		{
			name: "successful reading",
			mocks: func(dbMock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "date", "title", "comment", "repeat"}).
					AddRow(1, "20240203", "Test title task 1", "Comment for task 1", "d 7").
					AddRow(2, "20240203", "Test title task 2", "Comment for task 2", "d 7").
					AddRow(3, "20240203", "Test title task 3", "Comment for task 3", "d 7")
				query := regexp.QuoteMeta("SELECT id, date, title, comment, repeat FROM scheduler WHERE title LIKE ? OR comment LIKE ? ORDER BY date LIMIT ?")
				dbMock.ExpectQuery(query).
					WithArgs(payload, payload, 3).
					WillReturnRows(rows)
			}, // SELECT id, date, title, comment, repeat FROM scheduler WHERE date = ? LIMIT ?
			args: args{
				ctx:     context.Background(),
				payload: payload,
			},
			wantTasks: func(tt require.TestingT, got interface{}, i ...interface{}) {
				tasks, ok := got.([]models.Task)
				require.True(t, ok)
				assert.Len(t, tasks, 3)

				assert.Equal(t, int64(1), tasks[0].ID, i...)
				assert.Equal(t, "20240203", tasks[0].Date, i...)
				assert.Equal(t, "Test title task 1", tasks[0].Title, i...)
				assert.Equal(t, "Comment for task 1", tasks[0].Comment, i...)
				assert.Equal(t, "d 7", tasks[0].Repeat, i...)

				assert.Equal(t, int64(2), tasks[1].ID, i...)
				assert.Equal(t, "20240203", tasks[1].Date, i...)
				assert.Equal(t, "Test title task 2", tasks[1].Title, i...)
				assert.Equal(t, "Comment for task 2", tasks[1].Comment, i...)
				assert.Equal(t, "d 7", tasks[1].Repeat, i...)

				assert.Equal(t, int64(3), tasks[2].ID, i...)
				assert.Equal(t, "20240203", tasks[2].Date, i...)
				assert.Equal(t, "Test title task 3", tasks[2].Title, i...)
				assert.Equal(t, "Comment for task 3", tasks[2].Comment, i...)
				assert.Equal(t, "d 7", tasks[2].Repeat, i...)
			},
			wantErr: require.NoError,
		},
		{
			name: "no rows",
			mocks: func(dbMock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "date", "title", "comment", "repeat"})
				query := regexp.QuoteMeta("SELECT id, date, title, comment, repeat FROM scheduler WHERE title LIKE ? OR comment LIKE ? ORDER BY date LIMIT ?")
				dbMock.ExpectQuery(query).
					WithArgs(payload, payload, 3).
					WillReturnRows(rows)
			},
			args: args{
				ctx:     context.Background(),
				payload: payload,
			},
			wantTasks: func(tt require.TestingT, got interface{}, i ...interface{}) {
				tasks, ok := got.([]models.Task)
				require.True(t, ok)
				require.Empty(t, tasks)
			},
			wantErr: require.NoError,
		},
		{
			name: "database error",
			mocks: func(dbMock sqlmock.Sqlmock) {
				query := regexp.QuoteMeta("SELECT id, date, title, comment, repeat FROM scheduler WHERE title LIKE ? OR comment LIKE ? ORDER BY date LIMIT ?")
				dbMock.ExpectQuery(query).
					WithArgs(payload, payload, 3).
					WillReturnError(errors.New("database error"))
			},
			args: args{
				ctx:     context.Background(),
				payload: payload,
			},
			wantTasks: func(tt require.TestingT, got interface{}, i ...interface{}) {
				tasks, ok := got.([]models.Task)
				require.True(t, ok)
				require.Empty(t, tasks)
			},
			wantErr: func(tt require.TestingT, err error, i ...interface{}) {
				require.EqualError(tt, err, "cannot execute query: database error")
			},
		},
		{
			name:  "empty payload",
			mocks: func(dbMock sqlmock.Sqlmock) {},
			args: args{
				ctx:     context.Background(),
				payload: "",
			},
			wantTasks: func(tt require.TestingT, got interface{}, i ...interface{}) {
				tasks, ok := got.([]models.Task)
				require.True(t, ok)
				require.Empty(t, tasks)
			},
			wantErr: func(tt require.TestingT, err error, i ...interface{}) {
				require.EqualError(tt, err, storage.ErrEmptyPayload.Error())
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			db, dbMock, err := sqlmock.New()
			require.NoError(t, err)

			storage := sqlite.New(db, 3)
			tt.mocks(dbMock)

			tasks, err := storage.ReadByPayload(tt.args.ctx, tt.args.payload)
			tt.wantErr(t, err)
			tt.wantTasks(t, tasks)

			require.NoError(t, dbMock.ExpectationsWereMet())
		})
	}
}

func TestTaskStorage_Update(t *testing.T) {
	t.Parallel()

	type args struct {
		ctx  context.Context
		task *models.Task
	}

	var (
		id      int64  = 1
		date    string = "20250204"
		title   string = "test title"
		comment string = "test comment"
		repeat  string = "test repeat"
	)

	tests := []struct {
		name    string
		mocks   func(dbMock sqlmock.Sqlmock)
		args    args
		wantErr require.ErrorAssertionFunc
	}{
		{
			name: "successful update",
			mocks: func(dbMock sqlmock.Sqlmock) {
				query := regexp.QuoteMeta("UPDATE scheduler SET date = ?, title = ?, comment = ?, repeat = ? WHERE id = ?")
				dbMock.ExpectExec(query).
					WithArgs(date, title, comment, repeat, id).
					WillReturnResult(sqlmock.NewResult(0, 1))
			},
			args: args{
				ctx: context.Background(),
				task: &models.Task{
					ID:      id,
					Date:    date,
					Title:   title,
					Comment: comment,
					Repeat:  repeat,
				},
			},
			wantErr: require.NoError,
		},
		{
			name: "no rows affected",
			mocks: func(dbMock sqlmock.Sqlmock) {
				query := regexp.QuoteMeta("UPDATE scheduler SET date = ?, title = ?, comment = ?, repeat = ? WHERE id = ?")
				dbMock.ExpectExec(query).
					WithArgs(date, title, comment, repeat, id).
					WillReturnResult(sqlmock.NewResult(0, 0))
			},
			args: args{
				ctx: context.Background(),
				task: &models.Task{
					ID:      id,
					Date:    date,
					Title:   title,
					Comment: comment,
					Repeat:  repeat,
				},
			},
			wantErr: require.NoError,
		},
		{
			name: "database error",
			mocks: func(dbMock sqlmock.Sqlmock) {
				query := regexp.QuoteMeta("UPDATE scheduler SET date = ?, title = ?, comment = ?, repeat = ? WHERE id = ?")
				dbMock.ExpectExec(query).
					WithArgs(date, title, comment, repeat, id).
					WillReturnError(errors.New("database error"))
			},
			args: args{
				ctx: context.Background(),
				task: &models.Task{
					ID:      id,
					Date:    date,
					Title:   title,
					Comment: comment,
					Repeat:  repeat,
				},
			},
			wantErr: func(tt require.TestingT, err error, i ...interface{}) {
				require.EqualError(tt, err, "failed to update task: database error", i...)
			},
		},
		{
			name: "nil task update",
			mocks: func(dbMock sqlmock.Sqlmock) {
			},
			args: args{
				ctx:  context.Background(),
				task: nil,
			},
			wantErr: func(tt require.TestingT, err error, i ...interface{}) {
				require.EqualError(tt, err, storage.ErrNilTaskUpdate.Error(), i...)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			db, dbMock, err := sqlmock.New()
			require.NoError(t, err)

			storage := sqlite.New(db, 3)
			tt.mocks(dbMock)

			err = storage.Update(tt.args.ctx, tt.args.task)
			tt.wantErr(t, err)

			require.NoError(t, dbMock.ExpectationsWereMet())
		})
	}
}

func TestTaskStorage_Delete(t *testing.T) {
	t.Parallel()

	var id int64 = 1

	type args struct {
		ctx context.Context
		id  int64
	}

	tests := []struct {
		name    string
		mocks   func(dbMock sqlmock.Sqlmock)
		args    args
		wantErr require.ErrorAssertionFunc
	}{
		{
			name: "successful deletion",
			mocks: func(dbMock sqlmock.Sqlmock) {
				query := regexp.QuoteMeta("DELETE FROM scheduler WHERE id = ?")
				dbMock.ExpectExec(query).
					WithArgs(id).
					WillReturnResult(sqlmock.NewResult(0, 1))
			},
			args: args{
				ctx: context.Background(),
				id:  id,
			},
			wantErr: require.NoError,
		},
		{
			name: "no rows affected",
			mocks: func(dbMock sqlmock.Sqlmock) {
				query := regexp.QuoteMeta("DELETE FROM scheduler WHERE id = ?")
				dbMock.ExpectExec(query).
					WithArgs(id).
					WillReturnResult(sqlmock.NewResult(0, 0))
			},
			args: args{
				ctx: context.Background(),
				id:  id,
			},
			wantErr: require.NoError,
		},
		{
			name: "database error",
			mocks: func(dbMock sqlmock.Sqlmock) {
				query := regexp.QuoteMeta("DELETE FROM scheduler WHERE id = ?")
				dbMock.ExpectExec(query).
					WithArgs(id).
					WillReturnError(errors.New("database error"))
			},
			args: args{
				ctx: context.Background(),
				id:  id,
			},
			wantErr: func(tt require.TestingT, err error, i ...interface{}) {
				require.EqualError(tt, err, "failed to delete task: database error", i...)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			db, dbMock, err := sqlmock.New()
			require.NoError(t, err)

			storage := sqlite.New(db, 3)
			tt.mocks(dbMock)

			err = storage.Delete(tt.args.ctx, tt.args.id)
			tt.wantErr(t, err)

			require.NoError(t, dbMock.ExpectationsWereMet())
		})
	}
}
