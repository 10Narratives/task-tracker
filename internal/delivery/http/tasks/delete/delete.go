package delete

import (
	"context"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/go-chi/render"
)

// TODO: Написать документацию для структуры ответа и интерфейса
// TODO: Написать документацию для обработчика с помощью Swagger

const op = "http.Delete"

type Response struct {
	Err string `json:"error,omitempty"`
}

//go:generate go run github.com/vektra/mockery/v2@v2.52.1 --name=TaskRemover
type TaskRemover interface {
	Delete(ctx context.Context, id int64) error
}

func New(logger *slog.Logger, tr TaskRemover) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		param := r.URL.Query().Get("id")
		logger := logger.With(slog.String("op", op), slog.String("id", param))

		id, err := strconv.Atoi(param)
		if err != nil {
			logger.Error("gotten invalid id")
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, Response{Err: "gotten invalid id"})
			return
		}
		err = tr.Delete(context.Background(), int64(id))
		if err != nil {
			logger.Error("failed to delete task")
			w.WriteHeader(http.StatusInternalServerError)
			render.JSON(w, r, Response{Err: "failed to delete task"})
			return
		}
		logger.Info("task was deleted")
		render.JSON(w, r, Response{})
	}
}
