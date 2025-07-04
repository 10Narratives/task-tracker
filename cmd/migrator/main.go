package main

import (
	"log/slog"

	migratorcfg "github.com/10Narratives/task-tracker/internal/config/migrator"
	"github.com/10Narratives/task-tracker/internal/lib/logging/sl"
)

func main() {
	cfg := migratorcfg.MustLoad()
	log := sl.MustLogger(
		sl.WithLevel(cfg.Logger.Level),
		sl.WithFormat(cfg.Logger.Format),
		sl.WithOutput(cfg.Logger.Output),
	)

	log.Info("running migrator with configuration: ", slog.Any("config", cfg.Migrate))
}
