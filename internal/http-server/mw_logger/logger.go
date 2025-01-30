package mw_logger

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5/middleware"
)

// New creates and returns an HTTP middleware that logs incoming requests.
//
// This middleware logs request details, including:
//   - HTTP method
//   - Request path
//   - Remote address
//   - User agent
//   - Request ID (retrieved from context)
//
// Additionally, it logs the response details after the request completes:
//   - HTTP status code
//   - Number of bytes written
//   - Duration of request processing
//
// Parameters:
//   - logger: A *slog.Logger instance used for structured logging.
//
// Returns:
//   - A middleware function that wraps an http.Handler and logs request details.
func New(logger *slog.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		logger = logger.With(
			slog.String("component", "middleware/logger"),
		)

		logger.Info("Logger middleware is enabled")

		fn := func(w http.ResponseWriter, r *http.Request) {
			entry := slog.With(
				slog.String("method", r.Method),
				slog.String("path", r.URL.Path),
				slog.String("remote_addr", r.RemoteAddr),
				slog.String("user_agent", r.UserAgent()),
				slog.String("request_id", middleware.GetReqID(r.Context())),
			)
			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)

			startTime := time.Now()
			defer func() {
				entry.Info("request completed",
					slog.Int("status", ww.Status()),
					slog.Int("bytes", ww.BytesWritten()),
					slog.String("duration", time.Since(startTime).String()),
				)
			}()

			next.ServeHTTP(ww, r)
		}

		return http.HandlerFunc(fn)
	}
}
