package app

import (
	"context"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/10Narratives/task-tracker/internal/config"
	mw_auth "github.com/10Narratives/task-tracker/internal/delivery/http/middleware/auth"
	mw_logging "github.com/10Narratives/task-tracker/internal/delivery/http/middleware/logging"
	next "github.com/10Narratives/task-tracker/internal/delivery/http/nextdate"
	"github.com/10Narratives/task-tracker/internal/delivery/http/singin"
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
	app.logger.Info("starting to create database connection")
	db, close, err := storage.OpenDB(app.cfg.Storage.DriverName, app.cfg.Storage.DataSourceName)
	if err != nil {
		app.logger.Error(err.Error())
		os.Exit(1)
	}
	defer close()
	app.logger.Info("database connection created successfully")

	app.logger.Info("starting to initialize task service")
	store := sqlite.New(db, app.cfg.Storage.PaginationLimit)
	err = store.Prepare()
	if err != nil {
		app.logger.Error("can not prepare database:" + err.Error())
		os.Exit(1)
	}
	service := tasks.New(store)
	app.logger.Info("task service initialized successfully")

	app.logger.Info("starting to initialize router")
	router := chi.NewRouter()
	router.Use(mw_logging.New(app.logger))

	router.Handle("/*", http.StripPrefix("/", http.FileServer(http.Dir(app.cfg.HTTP.FileServerPath))))

	router.Post("/api/signin", singin.New(app.logger))

	router.Route("/api", func(r chi.Router) {
		r.Use(mw_auth.Auth)

		router.Post("/api/task", register.New(app.logger, service))
		router.Get("/api/tasks", read.New(app.logger, service))
		router.Get("/api/task", readone.New(app.logger, service))
		router.Put("/api/task", update.New(app.logger, service))
		router.Delete("/api/task", delete.New(app.logger, service))
		router.Post("/api/task/done", complete.New(app.logger, service))
		router.Delete("/api/task/done", delete.New(app.logger, service))
	})

	router.Get("/api/nextdate", next.New(app.logger))

	router.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("http://localhost:8000/swagger/doc.json"),
	))

	app.logger.Info("router initialized successfully")

	srv := &http.Server{
		Addr:         app.cfg.HTTP.Address + ":" + app.cfg.HTTP.Port,
		Handler:      router,
		ReadTimeout:  app.cfg.HTTP.Timeout,
		WriteTimeout: app.cfg.HTTP.Timeout,
		IdleTimeout:  app.cfg.HTTP.IdleTimeout,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			app.logger.Error("failed to start HTTP server", slog.String("error", err.Error()))
		}
	}()

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGTERM)

	app.logger.Info("server started")

	<-done
	app.logger.Info("stopping server")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		app.logger.Error("failed to stop server")

		return
	}

	app.logger.Info("server stopped")
}
