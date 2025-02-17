package delete

import (
	"context"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/go-chi/render"
)

const op = "http.Delete"

type Response struct {
	Err string `json:"error,omitempty"`
}

//go:generate go run github.com/vektra/mockery/v2@v2.52.1 --name=TaskRemover
type TaskRemover interface {
	Delete(ctx context.Context, id int64) error
}

// @Summary Delete task by its ID
// @Description Permanently remove a task from the system
// @Produce json
// @Param id query int true "Task ID"
// @Success 200 {object} Response
// @Failure 400 {object} Response "Invalid task ID"
// @Failure 500 {object} Response "Failed to delete task"
// @Router /api/task [delete]
func New(logger *slog.Logger, tr TaskRemover) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		param := r.URL.Query().Get("id")
		logger := logger.With(slog.String("op", op), slog.String("id", param))

		id, err := strconv.Atoi(param)
		if err != nil {
			logger.Error("gotten invalid id" + err.Error())
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
