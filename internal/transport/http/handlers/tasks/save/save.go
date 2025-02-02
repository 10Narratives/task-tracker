package save

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
)

var (
	ErrJSONDeserialization = errors.New("ошибка десериализации JSON")
	ErrNoTaskTitle         = errors.New("не указан заголовок задачи")
	ErrWrongDateFormat     = errors.New("дата представлена в формате, отличном от 20060102")
	ErrInvalidRepeatRule   = errors.New("правило повторения указано в неправильном формате")
)

type Request struct {
	Date    string `json:"date"`
	Title   string `json:"title" validate:"required,title"`
	Comment string `json:"comment,omitempty"`
	Repeat  string `json:"repeat"`
}

type Response struct {
	ID  string `json:"id,omitempty"`
	Err string `json:"error,omitempty"`
}

type TaskSaver interface {
	Save(date, title, comment, repeat string) (int64, error)
}

func New(log *slog.Logger, taskSaver TaskSaver) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req Request

		err := render.DecodeJSON(r.Body, &req)
		if err != nil {
			render.JSON(w, r, Response{
				Err: err.Error(),
			})

			return
		}

		if err := validator.New().Struct(req); err != nil {
			render.JSON(w, r, Response{
				Err: "invalid request",
			})

			return
		}
	}
}
