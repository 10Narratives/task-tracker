package read

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/10Narratives/task-tracker/internal/models"
	"github.com/go-chi/render"
)

const op = "http.Read"

type Response struct {
	Tasks []models.Task `json:"tasks,omitempty"`
	Err   string        `json:"error,omitempty"`
}

//go:generate go run github.com/vektra/mockery/v2@v2.52.1 --name=TaskReader
type TaskReader interface {
	Tasks(ctx context.Context, search string) ([]models.Task, error)
}

// @Summary Get tasks
// @Description Retrieve tasks optionally filtered by a search query
// @Produce json
// @Param search query string false "Search filter"
// @Success 200 {object} Response
// @Failure 500 {object} Response "Failed to read tasks"
// @Router /api/tasks [get]
func New(log *slog.Logger, tr TaskReader) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		search := r.URL.Query().Get("search")

		logger := log.With(slog.String("op", op), slog.String("search", search))
		tasks, err := tr.Tasks(context.Background(), search)
		if err != nil {
			logger.Error(err.Error())
			logger.Error("failed to read tasks")
			w.WriteHeader(http.StatusInternalServerError)
			render.JSON(w, r, Response{Err: "failed to read tasks"})
		} else {
			logger.Info("success task reading")
			w.WriteHeader(http.StatusOK)
			render.JSON(w, r, Response{Tasks: tasks})
		}
	}
}
