package main

import (
	"context"
	"errors"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/gofiber/fiber/v3"

	"jianco-api/config"
	"jianco-api/internal/app"
	"jianco-api/internal/storage"
	"jianco-api/pkg/database"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		slog.Error("failed to load config", "error", err)
		os.Exit(1)
	}

	if cfg.IsProduction() {
		slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, nil)))
	} else {
		slog.SetDefault(slog.New(slog.NewTextHandler(os.Stdout, nil)))
	}

	db, err := database.NewPostgres(cfg.DatabaseURL, cfg.IsProduction(), cfg.DBMaxOpenConn, cfg.DBMaxIdleConn)
	if err != nil {
		slog.Error("failed to connect database", "error", err)
		os.Exit(1)
	}

	storageClient, err := storage.New(cfg.R2)
	if err != nil {
		slog.Error("failed to initialize storage", "error", err)
		os.Exit(1)
	}

	application := app.New(cfg, db, storageClient)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	slog.Info("starting api", "addr", cfg.Addr(), "env", cfg.AppEnv)
	if err := application.Fiber.Listen(cfg.Addr(), fiber.ListenConfig{GracefulContext: ctx}); err != nil && !errors.Is(err, context.Canceled) {
		slog.Error("api stopped with error", "error", err)
		os.Exit(1)
	}
	slog.Info("api stopped")
}
