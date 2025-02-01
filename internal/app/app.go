package app

import (
	"log"
	"log/slog"
	"net/http"
	"os"

	"github.com/10Narratives/task-tracker/internal/config"
	"github.com/10Narratives/task-tracker/internal/logging"
	"github.com/10Narratives/task-tracker/internal/storage/sqlite"
	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
)

type App struct {
	cfg     *config.Config
	logger  *slog.Logger
	storage *sqlite.Storage
	router  *chi.Mux
}

func New() App {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Can not load environment variables.")
	}

	cfg := config.MustConfig()
	logger := logging.MustLogger(cfg.Env)

	storage, err := sqlite.New(cfg.Storage.DriverName, cfg.Storage.DataSourceName)
	if err != nil {
		logger.Warn("Can not initialize storage")
		os.Exit(1)
	}

	router := chi.NewRouter()
	router.Handle("/*", http.StripPrefix("/", http.FileServer(http.Dir(cfg.HTTP.FileServerPath))))

	return App{cfg, logger, storage, router}
}

func (app *App) Run() {
	fullAddr := app.cfg.HTTP.Address + ":" + app.cfg.HTTP.Port
	if err := http.ListenAndServe(fullAddr, app.router); err != nil {
		app.logger.Error("Can not start up http server", slog.String("occurred", err.Error()))
	}
}
