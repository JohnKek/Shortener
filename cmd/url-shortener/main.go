package main

import (
	"UrlShort/internal/config"
	mwLogger "UrlShort/internal/http-server/middleware/logger"
	"UrlShort/internal/storage/sqlite"
	"UrlShort/utils"
	"UrlShort/utils/logger/handlers/slogpretty"
	"fmt"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"log/slog"
	"os"
)

const (
	ENV_LOCAL = "local"
	ENV_DEV   = "dev"
	ENV_PROD  = "prod"
)

func main() {
	cfg := config.MustLoad()
	fmt.Println(cfg)
	logger := setupLogger(cfg.Env)
	logger.Info(
		"Starting url shortener",
		slog.String("env", cfg.Env))
	logger.Debug("Debug logging enabled")
	storage, err := sqlite.New(cfg.StoragePath)
	if err != nil {
		logger.Error("failed to init sqlite storage", utils.Err(err))
		os.Exit(1)
	}
	_ = storage
	router := chi.NewRouter()
	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(mwLogger.New(logger))
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)

	// TODO: init storage: sqllite
	// TODO: init router: chi, render

}
func setupLogger(env string) *slog.Logger {
	var log *slog.Logger
	switch env {
	case ENV_LOCAL:
		log = setupPrettySlog()
	case ENV_DEV:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case ENV_PROD:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	}
	return log
}

func setupPrettySlog() *slog.Logger {
	opts := slogpretty.PrettyHandlerOptions{
		SlogOpts: &slog.HandlerOptions{
			Level: slog.LevelDebug,
		},
	}

	handler := opts.NewPrettyHandler(os.Stdout)

	return slog.New(handler)
}
