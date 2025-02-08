package logging

// TODO: Replace logger into `/pkg`

import (
	"log"
	"log/slog"
	"os"

	"github.com/10Narratives/task-tracker/internal/config"
	"github.com/10Narratives/task-tracker/pkg/logging/slogpretty"
)

// MustLogger initializes and returns a structured logger based on the provided environment.
//
// The logging format and log level are determined by the environment:
//   - Local (config.EnvLocal): Text format with Debug level.
//   - Development (config.EnvDev): JSON format with Debug level.
//   - Production (config.EnvProd): JSON format with Info level.
//
// If an unknown environment is provided, the function logs a fatal error and terminates the program.
//
// Parameters:
//   - env: The runtime environment (expected values: config.EnvLocal, config.EnvDev, config.EnvProd).
//
// Returns:
//   - *slog.Logger: Configured logger instance.
func MustLogger(env string) *slog.Logger {
	var logger *slog.Logger

	switch env {
	case config.EnvLocal:
		logger = setupPrettySlog()
	case config.EnvDev:
		logger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case config.EnvProd:
		logger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	default:
		log.Fatal("Can not initialize logger: env parameter is unknown.")
	}

	return logger
}

func setupPrettySlog() *slog.Logger {
	opts := slogpretty.PrettyHandlerOptions{
		SlogOpts: &slog.HandlerOptions{
			Level: slog.LevelDebug,
		},
	}

	handler := opts.NewPrettyHandler(os.Stdout)

	return slog.New(handler)
}
