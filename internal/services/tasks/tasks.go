package tasks

import (
	"context"
	"errors"
	"time"

	"github.com/10Narratives/task-tracker/internal/models"
	"github.com/10Narratives/task-tracker/internal/services/nextdate"
	"github.com/10Narratives/task-tracker/internal/storage"
)

type TaskService struct {
	storage storage.TaskStorage
}

// New creates an instance of TaskService
func New(storage storage.TaskStorage) TaskService {
	return TaskService{storage: storage}
}

// Register creates task model
func (service TaskService) Register(ctx context.Context, date, title, comment, repeat string) (int64, error) {
	return service.storage.Create(ctx, date, title, comment, repeat)
}

func (service TaskService) Task(ctx context.Context, id int64) (models.Task, error) {
	return service.storage.Read(ctx, id)
}

func (service TaskService) Tasks(ctx context.Context) ([]models.Task, error) {
	return service.storage.ReadGroup(ctx)
}

func (service TaskService) Delete(ctx context.Context, id int64) error {
	return service.storage.Delete(ctx, id)
}

const (
	ByDate int = iota
	BySubstring
)

func (service TaskService) Search(ctx context.Context, condition int, payload string) ([]models.Task, error) {
	var (
		tasks []models.Task
		err   error
	)
	switch condition {
	case ByDate:
		tasks, err = service.storage.ReadByDate(ctx, payload)
		break
	case BySubstring:
		tasks, err = service.storage.ReadByPayload(ctx, payload)
		break
	default:
		tasks = make([]models.Task, 0)
		err = errors.New("unknown search condition")
		break
	}

	return tasks, err
}

func (service TaskService) Update(ctx context.Context, id int64, date, title, comment, repeat string) error {
	task := models.Task{ID: id, Date: date, Title: title, Comment: comment, Repeat: repeat}
	return service.storage.Update(ctx, &task)
}

func (service TaskService) Complete(ctx context.Context, id int64) error {
	task, err := service.storage.Read(ctx, id)
	if err != nil {
		return err
	}

	if len(task.Repeat) == 0 {
		err = service.Delete(ctx, id)
		if err != nil {
			return err
		}
	}

	parsed, _ := time.Parse(nextdate.DateLayout, task.Date)
	task.Date = nextdate.NextDate(time.Now(), parsed, task.Repeat)

	return nil
}
