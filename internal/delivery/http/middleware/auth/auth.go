package mw_auth

import (
	"errors"
	"net/http"
	"os"

	"github.com/golang-jwt/jwt/v5"
)

func Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		pass := os.Getenv("PASSWORD")
		if pass == "" {
			next.ServeHTTP(w, r)
			return
		}

		cookie, err := r.Cookie("token")
		if err != nil {
			http.Error(w, "Authentication required", http.StatusUnauthorized)
			return
		}

		tokenString := cookie.Value
		_, err = validateToken(tokenString, pass)
		if err != nil {
			http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func validateToken(tokenString, secret string) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}
		return []byte(secret), nil
	})

	if err != nil || !token.Valid {
		return nil, errors.New("invalid token")
	}

	return token, nil
}
