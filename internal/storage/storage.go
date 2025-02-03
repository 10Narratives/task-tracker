package storage

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/10Narratives/task-tracker/internal/models"
)

var (
	// ErrTaskNotFound    = errors.New("task not found")
	ErrNilTaskCreation = errors.New("cannot create task using nil pointer")
	ErrNilTaskUpdate   = errors.New("cannot update task using nil pointer")
	ErrEmptyDate       = errors.New("date cannot be empty for task retrieval")
	ErrEmptyPayload    = errors.New("payload cannot be empty for task retrieval")
	ErrLimitOutOfRange = errors.New("limit must be within the allowed range")
)

type TaskStorage interface {
	Create(ctx context.Context, t *models.Task) error
	Read(ctx context.Context, id int64) (models.Task, error)
	ReadGroup(ctx context.Context) ([]models.Task, error)
	ReadByDate(ctx context.Context, date string) ([]models.Task, error)
	ReadByPayload(ctx context.Context, payload string) ([]models.Task, error)
	Update(ctx context.Context, t *models.Task) error
	Delete(ctx context.Context, id int64) error
}

func OpenDB(driver string, dsn string) (*sql.DB, func(), error) {
	db, err := sql.Open(driver, dsn)
	if err != nil {
		return nil, nil, fmt.Errorf("cannot open database: %w", err)
	}

	closeFunc := func() {
		_ = db.Close()
	}

	return db, closeFunc, nil
}
