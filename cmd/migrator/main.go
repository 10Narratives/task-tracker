package main

import (
	"errors"
	"fmt"
	"log/slog"
	"os"

	migratorcfg "github.com/10Narratives/task-tracker/internal/config/migrator"
	"github.com/10Narratives/task-tracker/internal/lib/logging/sl"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/sqlite3"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	cfg := migratorcfg.MustLoad()
	log := sl.MustLogger(
		sl.WithLevel(cfg.Logger.Level),
		sl.WithFormat(cfg.Logger.Format),
		sl.WithOutput(cfg.Logger.Output),
	)

	log.Info("running migrator with configuration: ", slog.Any("config", cfg.Migrate))

	m, err := migrate.New(
		"file://"+cfg.Migrate.MigrationsPath,
		fmt.Sprintf("sqlite3://%s", cfg.Migrate.StoragePath),
	)
	if err != nil {
		if os.IsNotExist(err) {
			log.Error("migrations directory does not exist", slog.String("path", cfg.Migrate.MigrationsPath))
		} else {
			log.Error("cannot create new migrate instance", slog.Any("error", err))
		}
		os.Exit(1)
	}

	if err := m.Up(); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			log.Warn("no migrations to apply")
			return
		}

		log.Error("cannot apply migrations: ", slog.Any("error", err.Error()))
		os.Exit(1)
	}

	log.Info("migrations applied")
}
