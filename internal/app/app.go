package app

import (
	"log"
	"log/slog"
	"net/http"

	"github.com/10Narratives/task-tracker/internal/config"
	"github.com/10Narratives/task-tracker/internal/storage/sqlite"
	"github.com/10Narratives/task-tracker/pkg/logging"
	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
)

type App struct {
	cfg     *config.Config
	logger  *slog.Logger
	storage *sqlite.TaskStorage
	router  *chi.Mux
}

// TODO: Make database connection

func New() App {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Can not load environment variables.")
	}

	cfg := config.MustConfig()
	logger := logging.MustLogger(cfg.Env)

	router := chi.NewRouter()
	router.Handle("/*", http.StripPrefix("/", http.FileServer(http.Dir(cfg.HTTP.FileServerPath))))

	return App{cfg, logger, nil, router}
}

func (app *App) Run() {
	fullAddr := app.cfg.HTTP.Address + ":" + app.cfg.HTTP.Port
	if err := http.ListenAndServe(fullAddr, app.router); err != nil {
		app.logger.Error("Can not start up http server", slog.String("occurred", err.Error()))
	}
}
