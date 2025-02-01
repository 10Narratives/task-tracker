package main

import (
	"log"
	"log/slog"
	"net/http"
	"os"

	"github.com/10Narratives/task-tracker/internal/config"
	"github.com/10Narratives/task-tracker/internal/logging"
	"github.com/10Narratives/task-tracker/internal/storage/sqlite"
	"github.com/10Narratives/task-tracker/internal/transport/http/mw_logger"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Can not load environment variables.")
	}

	cfg := config.MustConfig()

	logger := logging.MustLogger(cfg.Env)
	logger.Info("Start up task-tracker", slog.String("env", cfg.Env))
	logger.Debug("Debug messages are enabled")

	logger.Info("Open database",
		slog.String("driver", cfg.Storage.DriverName),
		slog.String("data_source_name", cfg.Storage.DataSourceName),
	)
	storage, err := sqlite.New(cfg.Storage.DriverName, cfg.Storage.DataSourceName)
	if err != nil {
		logger.Error("Failed to init storage", slog.String("occurred", err.Error()))
		os.Exit(1)
	}

	_ = storage

	logger.Info("Start up HTTP server",
		slog.String("address", cfg.HTTP.Address),
		slog.String("port", cfg.HTTP.Port),
		slog.String("timeout", cfg.HTTP.Timeout.String()),
		slog.String("idle_timeout", cfg.HTTP.IdleTimeout.String()),
		slog.String("file_server_path", cfg.HTTP.FileServerPath),
	)
	router := chi.NewRouter()
	router.Use(middleware.RequestID)
	router.Use(mw_logger.New(logger))
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)
	router.Handle("/*", http.StripPrefix("/", http.FileServer(http.Dir(cfg.HTTP.FileServerPath))))

	fullAddr := cfg.HTTP.Address + ":" + cfg.HTTP.Port
	if err := http.ListenAndServe(fullAddr, router); err != nil {
		logger.Error("Can not start up http server", slog.String("occurred", err.Error()))
	}

}
