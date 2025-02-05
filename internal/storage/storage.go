package storage

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/10Narratives/task-tracker/internal/models"
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
