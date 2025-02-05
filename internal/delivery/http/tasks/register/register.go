package register

import (
	"context"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/10Narratives/task-tracker/internal/delivery/http/validation"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
)

type Request struct {
	Date    string `json:"date" validate:"required,dateformat"`
	Title   string `json:"title" validate:"required,title"`
	Comment string `json:"comment,omitempty"`
	Repeat  string `json:"repeat" validate:"repeat"`
}

type Response struct {
	ID  string `json:"id,omitempty"`
	Err string `json:"error,omitempty"`
}

//go:generate go run github.com/vektra/mockery/v2@v2.28.2 --name=TaskRegistrar
type TaskRegistrar interface {
	Register(ctx context.Context, date, title, comment, repeat string) (int64, error)
}

func New(log *slog.Logger, ts TaskRegistrar) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req Request
		err := render.DecodeJSON(r.Body, &req)
		if errors.Is(err, io.EOF) {
			log.Error("request body is empty")
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, Response{Err: "empty request"})

			return
		}

		if err != nil {
			log.Error("failed to decode request body")
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, Response{Err: "failed to decode request body"})

			return
		}

		log.Info("request body decoded", slog.Any("request", req))

		v := validator.New()
		v.RegisterValidation("dateformat", validation.IsDateValid)
		v.RegisterValidation("title", validation.IsTitleValid)
		v.RegisterValidation("repeat", validation.IsRepeatValid)
		if err := v.Struct(req); err != nil {
			validationErr := err.(validator.ValidationErrors)

			log.Error("invalid request")
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, Response{Err: validation.ValidationErrorMsg(validationErr)})

			return
		}

		id, err := ts.Register(context.Background(), req.Date, req.Title, req.Comment, req.Repeat)
		if err != nil {
			log.Error(err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			render.JSON(w, r, Response{Err: "failed to add task"})

			return
		}

		render.JSON(w, r, Response{ID: strconv.Itoa(int(id))})
	}
}
