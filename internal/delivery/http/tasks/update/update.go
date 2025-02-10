package update

import (
	"context"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/10Narratives/task-tracker/internal/delivery/http/validation"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
)

// TODO: Написать документацию для структуры ответа и интерфейса
// TODO: Написать документацию для обработчика с помощью Swagger

const op = "http.Update"

type Request struct {
	ID      string `json:"id" validate:"required"`
	Date    string `json:"date" validate:"required,dateformat"`
	Title   string `json:"title" validate:"required,title"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat" validate:"repeat"`
}

type Response struct {
	Err string `json:"error,omitempty"`
}

//go:generate go run github.com/vektra/mockery/v2@v2.52.1 --name=TaskUpdater
type TaskUpdater interface {
	Update(ctx context.Context, id int64, date, title, comment, repeat string) error
}

// @Summary Update an existing task
// @Description Modify the details of an existing task by its ID
// @Accept json
// @Produce json
// @Param request body Request true "Task data to update"
// @Success 200 {object} Response
// @Failure 400 {object} Response "Invalid request format or missing fields"
// @Failure 500 {object} Response "Failed to update task"
// @Router /api/task [put]
func New(logger *slog.Logger, tu TaskUpdater) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger := logger.With(slog.String("op", op))

		var req Request
		err := render.DecodeJSON(r.Body, &req)
		if err != nil {
			logger.Error("failed to decode request body")
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, Response{Err: "failed to decode request body"})
			return
		}

		logger.Info("request body decoded", slog.Any("request", req))

		v := validator.New()
		v.RegisterValidation("dateformat", validation.IsDateValid)
		v.RegisterValidation("title", validation.IsTitleValid)
		v.RegisterValidation("repeat", validation.IsRepeatValid)
		if err := v.Struct(req); err != nil {
			validationErr := err.(validator.ValidationErrors)
			logger.Error("invalid request")
			logger.Error(validation.ValidationErrorMsg(validationErr))
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, Response{Err: validation.ValidationErrorMsg(validationErr)})
			return
		}

		id, _ := strconv.Atoi(req.ID)
		err = tu.Update(context.Background(), int64(id), req.Date, req.Title, req.Comment, req.Repeat)
		if err != nil {
			logger.Error(err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			render.JSON(w, r, Response{Err: "failed to update task"})
			return
		}

		logger.Info("task was updated")
		render.JSON(w, r, Response{})
	}
}
