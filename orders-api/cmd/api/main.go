/*
Função main

Responsabilidades:
- carregar configurações;
- conectar ao RabbitMQ;
- criar publisher e consumer;
- executar migrations;
- abrir conexão com PostgreSQL;
- montar a aplicação;
- iniciar o servidor HTTP.
*/

package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "orders-api/docs"

	"orders-api/config"
	"orders-api/infrastructure/database"
	"orders-api/infrastructure/logging"
	"orders-api/infrastructure/messaging/rabbitmq"
	app "orders-api/internal/app"
)

func main() {
	logger := logging.New()
	slog.SetDefault(logger)

	cfg := config.Load()

	rabbitConnection, err := rabbitmq.NewConnection(cfg.RabbitMQURL)
	if err != nil {
		slog.Error(
			"failed to connect to rabbitmq",
			"error", err,
		)
		return
	}
	defer rabbitConnection.Close()

	rabbitPublisher, err := rabbitmq.NewPublisher(
		rabbitConnection,
		logger,
	)
	if err != nil {
		slog.Error(
			"failed to create rabbitmq publisher",
			"error", err,
		)
		return
	}
	defer rabbitPublisher.Close()

	rabbitConsumer, err := rabbitmq.NewConsumer(
		rabbitConnection,
		logger,
	)
	if err != nil {
		slog.Error(
			"failed to create rabbitmq consumer",
			"error", err,
		)
		return
	}
	defer rabbitConsumer.Close()

	if err := database.RunMigrations(cfg.DatabaseURL); err != nil {
		slog.Error(
			"failed to run migrations",
			"error", err,
		)
		return
	}

	db, err := database.New(cfg)
	if err != nil {
		slog.Error(
			"failed to connect to database",
			"error", err,
		)
		return
	}
	defer db.Close()
	httpHandler, orderEventHandler, err := app.NewApplication(
		cfg,
		db,
		rabbitPublisher,
	)

	if err != nil {
		slog.Error(
			"failed to create application",
			"error", err,
		)
		return
	}

	server := &http.Server{
		Addr:    ":" + cfg.AppPort,
		Handler: httpHandler,
	}

	ctx, stop := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
	)
	defer stop()

	go func() {
		slog.Info(
			"server started",
			"port", cfg.AppPort,
			"environment", cfg.AppEnv,
		)

		if err := server.ListenAndServe(); err != nil &&
			err != http.ErrServerClosed {
			slog.Error(
				"failed to start server",
				"error", err,
			)
		}
	}()

	go func() {
		err := rabbitConsumer.Consume(
			ctx,
			orderEventHandler.Handle,
		)

		if err != nil {
			slog.Error(
				"rabbitmq consumer stopped",
				"error", err,
			)
		}
	}()

	<-ctx.Done()

	slog.Info("shutdown signal received")

	shutdownCtx, cancel := context.WithTimeout(
		context.Background(),
		10*time.Second,
	)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		slog.Error(
			"server shutdown failed",
			"error", err,
		)
	}

	slog.Info("server stopped")
}
