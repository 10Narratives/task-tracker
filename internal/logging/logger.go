package logging

import (
	"log"
	"log/slog"
	"os"

	"github.com/10Narratives/task-tracker/internal/config"
)

func MustLogger(env string) *slog.Logger {
	var logger *slog.Logger

	switch env {
	case config.EnvLocal:
		logger = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case config.EnvDev:
		logger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case config.EnvProd:
		logger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	default:
		log.Fatal("Can not initialize logger: env parameter is unknown.")
	}

	return logger
}
