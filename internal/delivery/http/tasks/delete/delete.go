package delete

import (
	"context"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

// TODO: Сделать валидацию параметров URL с помощью validator
// TODO: Написать unit-тестирование для этого обработчика
// TODO: Написать документацию для структуры ответа и интерфейса
// TODO: Написать документацию для обработчика с помощью Swagger

type Response struct {
	Err string `json:"error,omitempty"`
}

type TaskRemover interface {
	Delete(ctx context.Context, id int64) error
}

func New(logger *slog.Logger, tr TaskRemover) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		param := chi.URLParam(r, "id")
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

		render.JSON(w, r, Response{})
	}
}
