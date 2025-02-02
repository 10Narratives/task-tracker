package sqlite

import (
	"database/sql"
	"errors"

	"github.com/10Narratives/task-tracker/internal/models"
	_ "github.com/mattn/go-sqlite3"
)

const PrepareSQLRequest = `
	CREATE TABLE IF NOT EXISTS scheduler (
    	id INTEGER PRIMARY KEY AUTOINCREMENT,
    	date TEXT NOT NULL,
    	title TEXT NOT NULL,
    	comment TEXT,
    	repeat TEXT CHECK(LENGTH(repeat) <= 128)
	);

	CREATE INDEX IF NOR EXISTS idx_scheduler_date ON scheduler(date);
`

var (
	ErrCanNotOpenDatabase     = errors.New("can not open database")
	ErrCanNotPrepareStatement = errors.New("can not prepare statement")
	ErrCanNotPrepareDatabase  = errors.New("can not prepare database")
)

type Storage struct {
	DB *sql.DB
}

func New(driver, dsn string) (*Storage, error) {
	db, err := sql.Open(driver, dsn)
	if err != nil {
		return nil, errors.Join(ErrCanNotOpenDatabase, err)
	}

	stmt, err := db.Prepare(PrepareSQLRequest)

	if err != nil {
		return nil, errors.Join(ErrCanNotPrepareStatement, err)
	}

	_, err = stmt.Exec()
	if err != nil {
		return nil, errors.Join(ErrCanNotPrepareDatabase, err)
	}

	return &Storage{DB: db}, nil
}

// Add
func (storage Storage) Add(task *models.Task) error {
	if task == nil {
		return errors.New("can not add task which pointer is equal to nil")
	}
	insertRes, err := storage.DB.Exec("INSERT INTO scheduler (date, title, comment, repeat) VALUES (?, ?, ?, ?)",
		task.Date, task.Title, task.Comment, task.Repeat)
	if err != nil {
		return err
	}

	lastTaskID, err := insertRes.LastInsertId()
	if err != nil {
		return nil
	}

	task.ID = lastTaskID
	return nil
}

// Delete
// GetTasks
// GetTask
