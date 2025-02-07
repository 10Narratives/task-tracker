package tasks

import (
	"context"
	"time"

	"github.com/10Narratives/task-tracker/internal/models"
	"github.com/10Narratives/task-tracker/internal/services/nextdate"
	"github.com/10Narratives/task-tracker/internal/storage"
)

// TODO: Перенести сюда интерфейс хранилища
// TODO: Написать тесты для всех методов
// TODO: Написать документацию

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

func (service TaskService) Tasks(ctx context.Context, search string) ([]models.Task, error) {
	var (
		tasks []models.Task
		err   error
	)
	if search == "" {
		tasks, err = service.storage.ReadGroup(ctx)
	} else if _, isDate := time.Parse("20060102", search); isDate == nil {
		tasks, err = service.storage.ReadByDate(ctx, search)
	} else {
		tasks, err = service.storage.ReadByPayload(ctx, search)
	}
	return tasks, err
}

func (service TaskService) Delete(ctx context.Context, id int64) error {
	return service.storage.Delete(ctx, id)
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
