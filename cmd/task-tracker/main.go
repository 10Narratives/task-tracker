package main

import (
	"fmt"
	"log"
	"log/slog"
	"net/http"

	"github.com/10Narratives/task-tracker/internal/config"
	"github.com/10Narratives/task-tracker/internal/http-server/mw_logger"
	"github.com/10Narratives/task-tracker/internal/logging"
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
	logger.Info("Start up task-tracker",
		slog.String("env", cfg.Env),
		slog.String("address", cfg.HTTP.Address),
		slog.String("port", cfg.HTTP.Port))
	logger.Debug("Debug messages are enabled")

	// TODO: Initialize storage

	router := chi.NewRouter()
	router.Use(middleware.RequestID)
	router.Use(mw_logger.New(logger))
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)
	router.Handle("/*", http.StripPrefix("/", http.FileServer(http.Dir(cfg.HTTP.FileServerPath))))

	fullAddr := cfg.HTTP.Address + ":" + cfg.HTTP.Port
	if err := http.ListenAndServe(fullAddr, router); err != nil {
		fmt.Println(err)
		logger.Error("Can not start up http server")
	}

}
