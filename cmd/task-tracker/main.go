package main

import (
	"log"
	"log/slog"

	"github.com/10Narratives/task-tracker/internal/config"
	"github.com/10Narratives/task-tracker/internal/logger"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Can not load environment variables.")
	}

	cfg := config.MustLoad()

	logger := logger.MustLogger(cfg.Env)
	logger.Info("Start up task-tracker", slog.String("env", cfg.Env))
	logger.Debug("Debug messages are enabled")

	// TODO: Initialize storage

	// TODO: Initialize http server
}
