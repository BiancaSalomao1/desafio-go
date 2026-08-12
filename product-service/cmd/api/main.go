package main

import (
	"log/slog"
	"net/http"

	"product-service/config"
	"product-service/infrastructure/database"
	app "product-service/internal/app"
)

func main() {
	cfg := config.Load()

	if err := database.RunMigrations(cfg.DatabaseURL); err != nil {
		slog.Error(
			"failed to run migrations",
			"service", "product-service",
			"error", err,
		)
		return
	}

	db, err := database.New(cfg)
	if err != nil {
		slog.Error(
			"failed to connect to database",
			"service", "product-service",
			"error", err,
		)
		return
	}
	defer db.Close()

	httpHandler, err := app.NewApplication(cfg, db)
	if err != nil {
		slog.Error(
			"failed to create application",
			"service", "product-service",
			"error", err,
		)
		return
	}

	server := &http.Server{
		Addr:    ":" + cfg.AppPort,
		Handler: httpHandler,
	}

	slog.Info(
		"product service started",
		"service", "product-service",
		"port", cfg.AppPort,
	)

	if err := server.ListenAndServe(); err != nil {
		slog.Error(
			"server stopped",
			"service", "product-service",
			"error", err,
		)
	}
}
