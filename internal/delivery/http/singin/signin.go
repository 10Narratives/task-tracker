package singin

import (
	"log/slog"
	"net/http"
	"os"

	"github.com/10Narratives/task-tracker/internal/delivery/http/validation"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v5"
)

// TODO: Make sign in handlers

const op = "http.Authentication"

type Request struct {
	Password string `json:"password" validate:"required"`
}

type Response struct {
	Token string `json:"token,omitempty"`
	Err   string `json:"error,omitempty"`
}

func New(log *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log := log.With(slog.String("op", op))

		//authHeader := r.Header.Get("Authorization")
		//if authHeader == "" {
		//	log.Error("unauthorized")
		//	w.WriteHeader(http.StatusUnauthorized)
		//	render.JSON(w, r, Response{Err: "unauthorized"})
		//	return
		//}

		var req Request
		err := render.DecodeJSON(r.Body, &req)
		if err != nil {
			log.Error("failed to decode request body")
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, Response{Err: "failed to decode request body"})
			return
		}
		log.Info("request body decoded")

		v := validator.New()
		if err := v.Struct(req); err != nil {
			validationErr := err.(validator.ValidationErrors)

			log.Error("invalid request")
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, Response{Err: validation.ValidationErrorMsg(validationErr)})

			return
		}

		// TODO: make JWT token
		pass := os.Getenv("PASSWORD")
		token := jwt.New(jwt.SigningMethodHS256)
		signedToken, err := token.SignedString(pass)
		if err != nil {
			log.Error("can not sign jwt token")
			w.WriteHeader(http.StatusInternalServerError)
			render.JSON(w, r, Response{Err: "can not sign jwt token"})
			return
		}

		render.JSON(w, r, Response{Token: signedToken})
	}
}
