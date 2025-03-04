package complete

import (
	"context"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/go-chi/render"
)

const op = "http.Complete"

type Response struct {
	Err string `json:"error,omitempty"`
}

//go:generate go run github.com/vektra/mockery/v2@v2.52.1 --name=TaskCompleter
type TaskCompleter interface {
	Complete(ctx context.Context, id int64) error
}

// @Summary Complete task by its ID
// @Description Mark task as completed, either by deleting it or modifying its deadline
// @Produce json
// @Param id query int true "Task ID"
// @Success 200 {object} Response
// @Failure 400 {object} Response "Invalid task ID"
// @Failure 500 {object} Response "Failed to complete task"
// @Router /api/task/done [post]
func New(log *slog.Logger, tc TaskCompleter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger := log.With("op", op)

		param := r.URL.Query().Get("id")
		id, err := strconv.Atoi(param)
		if err != nil {
			logger.Error("gotten invalid id")
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, Response{Err: "gotten invalid id"})
			return
		}

		err = tc.Complete(context.Background(), int64(id))
		if err != nil {
			logger.Error("failed to complete task")
			w.WriteHeader(http.StatusInternalServerError)
			render.JSON(w, r, Response{Err: "failed to complete task"})
			return
		}

		logger.Info("task was completed")
		render.JSON(w, r, Response{})
	}
}
