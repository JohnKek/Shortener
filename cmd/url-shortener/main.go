package main

import (
	"UrlShort/internal/config"
	"UrlShort/internal/storage/sqlite"
	"UrlShort/utils"
	"fmt"
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
	logger.Info("Starting url shortener", slog.String("env", cfg.Env))
	logger.Debug("Debug logging enabled")
	storage, err := sqlite.New(cfg.StoragePath)
	if err != nil {
		logger.Error("failed to init sqlite storage", utils.Err(err))
		os.Exit(1)
	}
	err = storage.SaveURL("https://www.google.com", "google")
	if err != nil {
		logger.Error("failed to save url", utils.Err(err))
	}
	// TODO: init storage: sqllite
	// TODO: init router: chi, render

}
func setupLogger(env string) *slog.Logger {
	var log *slog.Logger
	switch env {
	case ENV_LOCAL:
		log = slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case ENV_DEV:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case ENV_PROD:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	}
	return log
}
