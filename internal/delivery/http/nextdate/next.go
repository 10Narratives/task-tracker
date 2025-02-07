package next

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/10Narratives/task-tracker/internal/delivery/http/validation"
	"github.com/10Narratives/task-tracker/internal/services/nextdate"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
)

// TODO: Написать unit-тесты для этого обработчика
// TODO: Написать документацию для обработчика с помощью Swagger

type URLParams struct {
	Now    string `json:"now" validate:"required,dateformat"`
	Date   string `json:"date" validate:"required,dateformat"`
	Repeat string `json:"repeat" validate:"required,repeat"`
}

type Response struct {
	NextDate string `json:"nextdate,omitempty"`
	Err      string `json:"error,omitempty"`
}

func New(logger *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query()
		params := URLParams{
			Now:    query.Get("now"),
			Date:   query.Get("date"),
			Repeat: query.Get("repeat"),
		}

		v := validator.New()
		v.RegisterValidation("dateformat", validation.IsDateValid)
		v.RegisterValidation("repeat", validation.IsRepeatValid)
		if err := v.Struct(params); err != nil {
			validationErr := err.(validator.ValidationErrors)

			logger.Error("invalid request")
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, Response{Err: validation.ValidationErrorMsg(validationErr)})

			return
		}

		logger.Info("getting nextdate")

		now, _ := time.Parse("20060102", params.Now)
		date, _ := time.Parse("20060102", params.Now)
		nextdate := nextdate.NextDate(now, date, params.Repeat)
		render.JSON(w, r, Response{NextDate: nextdate})
	}
}
