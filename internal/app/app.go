package app

import (
	"log"
	"log/slog"
	"net/http"
	"os"

	"github.com/10Narratives/task-tracker/internal/config"
	"github.com/10Narratives/task-tracker/internal/delivery/http/tasks/complete"
	"github.com/10Narratives/task-tracker/internal/delivery/http/tasks/delete"
	"github.com/10Narratives/task-tracker/internal/delivery/http/tasks/read"
	"github.com/10Narratives/task-tracker/internal/delivery/http/tasks/readone"
	"github.com/10Narratives/task-tracker/internal/delivery/http/tasks/register"
	"github.com/10Narratives/task-tracker/internal/delivery/http/tasks/update"
	"github.com/10Narratives/task-tracker/internal/services/tasks"
	"github.com/10Narratives/task-tracker/internal/storage"
	"github.com/10Narratives/task-tracker/internal/storage/sqlite"
	"github.com/10Narratives/task-tracker/pkg/logging"
	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"

	_ "github.com/10Narratives/task-tracker/docs"
	httpSwagger "github.com/swaggo/http-swagger/v2"
)

type App struct {
	cfg    *config.Config
	logger *slog.Logger
}

func New() App {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Can not load environment variables.")
	}

	cfg := config.MustConfig()
	logger := logging.MustLogger(cfg.Env)

	return App{cfg, logger}
}

func (app *App) Run() {
	app.logger.Info("Starting to create database connection")
	db, close, err := storage.OpenDB(app.cfg.Storage.DriverName, app.cfg.Storage.DriverName)
	if err != nil {
		app.logger.Error(err.Error())
		os.Exit(1)
	}
	defer close()
	app.logger.Info("database connection created successfully")

	app.logger.Info("starting to initialize task service")
	store := sqlite.New(db, 10)
	store.Prepare()
	service := tasks.New(store)
	app.logger.Info("task service initialized successfully")

	app.logger.Info("starting to initialize router")
	router := chi.NewRouter()
	router.Handle("/*", http.StripPrefix("/", http.FileServer(http.Dir(app.cfg.HTTP.FileServerPath))))
	router.Post("/api/task", register.New(app.logger, service))
	router.Get("/api/tasks", read.New(app.logger, service))
	router.Get("/api/task", readone.New(app.logger, service))
	router.Put("/api/task", update.New(app.logger, service))
	router.Delete("/api/task", delete.New(app.logger, service))
	router.Post("/api/task/done", complete.New(app.logger, service))
	router.Delete("/api/task/done", delete.New(app.logger, service))

	router.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("http://localhost:8000/swagger/doc.json"),
	))

	app.logger.Info("router initialized successfully")

	fullAddr := app.cfg.HTTP.Address + ":" + app.cfg.HTTP.Port
	if err := http.ListenAndServe(fullAddr, router); err != nil {
		app.logger.Error("Can not start up http server", slog.String("occurred", err.Error()))
	}
}
