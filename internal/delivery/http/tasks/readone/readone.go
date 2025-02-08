package readone

import (
	"context"
	"log/slog"
	"net/http"
	"reflect"
	"strconv"

	"github.com/10Narratives/task-tracker/internal/models"
	"github.com/go-chi/render"
)

// TODO: Написать документацию для структуры ответа и интерфейса
// TODO: Написать документацию для обработчика с помощью Swagger

const op = "http.Readone"

type Response struct {
	ID      string `json:"id,omitempty"`
	Date    string `json:"date,omitempty"`
	Title   string `json:"title,omitempty"`
	Comment string `json:"comment,omitempty"`
	Repeat  string `json:"repeat,omitempty"`
	Err     string `json:"error,omitempty"`
}

//go:generate go run github.com/vektra/mockery/v2@v2.52.1 --name=TaskReader
type TaskReader interface {
	Task(ctx context.Context, id int64) (models.Task, error)
}

func New(logger *slog.Logger, tr TaskReader) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger := logger.With(slog.String("op", op))
		param := r.URL.Query().Get("id")
		id, err := strconv.Atoi(param)
		if err != nil {
			logger.Error("failed to covert id")
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, Response{Err: "gotten invalid id"})
			return
		}

		task, err := tr.Task(context.Background(), int64(id))
		if err != nil {
			logger.Error("failed to find task by id")
			w.WriteHeader(http.StatusInternalServerError)
			render.JSON(w, r, Response{Err: "failed to find task by id"})
			return
		}

		if reflect.DeepEqual(models.Task{}, task) {
			logger.Error("task not found")
			w.WriteHeader(http.StatusNotFound)
			render.JSON(w, r, Response{Err: "task not found"})
			return
		}

		logger.Info("task was found")
		w.WriteHeader(http.StatusOK)
		render.JSON(w, r, Response{
			ID:      param,
			Date:    task.Date,
			Title:   task.Title,
			Comment: task.Comment,
			Repeat:  task.Repeat,
		})
		return
	}
}
