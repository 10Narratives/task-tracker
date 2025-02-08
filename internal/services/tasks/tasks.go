package tasks

import (
	"context"
	"time"

	"github.com/10Narratives/task-tracker/internal/models"
	"github.com/10Narratives/task-tracker/internal/services/nextdate"
)

// TaskStorage is an interface for working with task storage.
// It defines methods for creating, reading, updating, and deleting tasks.
//
//go:generate go run github.com/vektra/mockery/v2@v2.28.2 --name=TaskStorage
type TaskStorage interface {

	// Create adds a new task to the storage and returns its ID and any error encountered.
	Create(ctx context.Context, date, title, comment, repeat string) (int64, error)

	// Read retrieves a task by its ID.
	// It returns the task and any error encountered.
	Read(ctx context.Context, id int64) (models.Task, error)

	// ReadGroup retrieves all tasks from the storage.
	// It returns a slice of tasks and any error encountered.
	ReadGroup(ctx context.Context) ([]models.Task, error)

	// ReadByDate retrieves tasks associated with a specific date.
	// It returns a slice of tasks and any error encountered.
	ReadByDate(ctx context.Context, date string) ([]models.Task, error)

	// ReadByPayload retrieves tasks that match a specific search payload.
	// It returns a slice of tasks and any error encountered.
	ReadByPayload(ctx context.Context, payload string) ([]models.Task, error)

	// Update modifies an existing task in the storage.
	// It returns any error encountered during the update.
	Update(ctx context.Context, t *models.Task) error

	// Delete removes a task from the storage by its ID.
	// It returns any error encountered during deletion.
	Delete(ctx context.Context, id int64) error
}

// TaskService manages tasks within the application.
type TaskService struct {
	// TaskStorage is the instance of task storage.
	storage TaskStorage
}

// New creates a new TaskService with the given TaskStorage.
func New(storage TaskStorage) TaskService {
	return TaskService{storage: storage}
}

// Register creates a new task with the specified details.
// It returns the ID of the created task and any error encountered.
func (service TaskService) Register(ctx context.Context, date, title, comment, repeat string) (int64, error) {
	return service.storage.Create(ctx, date, title, comment, repeat)
}

// Task retrieves a task by its ID.
// It returns the task and any error encountered.
func (service TaskService) Task(ctx context.Context, id int64) (models.Task, error) {
	return service.storage.Read(ctx, id)
}

// Tasks retrieves a list of tasks based on the search criteria.
// If search is empty, it returns all tasks. If search is a date, it returns tasks for that date.
// Otherwise, it searches for tasks matching the payload.
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

// Delete removes a task by its ID.
// It returns any error encountered during deletion.
func (service TaskService) Delete(ctx context.Context, id int64) error {
	return service.storage.Delete(ctx, id)
}

// Update modifies an existing task with the given details.
// It returns any error encountered during the update.
func (service TaskService) Update(ctx context.Context, id int64, date, title, comment, repeat string) error {
	task := models.Task{ID: id, Date: date, Title: title, Comment: comment, Repeat: repeat}
	return service.storage.Update(ctx, &task)
}

// Complete marks a task as complete.
// If the task is not recurring, it will be deleted.
// For recurring tasks, it updates the task date for the next occurrence.
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
		return nil
	}

	parsed, _ := time.Parse(nextdate.DateLayout, task.Date)
	task.Date = nextdate.NextDate(time.Now(), parsed, task.Repeat)

	return nil
}
