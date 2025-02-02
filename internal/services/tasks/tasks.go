package tasks

import (
	"github.com/10Narratives/task-tracker/internal/models"
	"github.com/10Narratives/task-tracker/internal/storage/sqlite"
)

type TasksService struct {
	storage *sqlite.Storage
}

func New(storage *sqlite.Storage) TasksService {
	return TasksService{storage: storage}
}

func (service TasksService) Register(date, title, comment, repeat string) (models.Task, error) {
	task := models.Task{
		Date:    date,
		Title:   title,
		Comment: comment,
		Repeat:  repeat,
	}

	err := service.storage.Add(&task)
	if err != nil {
		return models.Task{}, err
	}

	return task, nil
}
