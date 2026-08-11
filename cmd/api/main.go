/*
Função main

Responsabilidades:
- carregar configurações;
- executar migrations;
- abrir conexão com PostgreSQL;
- montar a aplicação;
- iniciar o servidor HTTP.

A montagem da aplicação (repositórios, casos de uso,
handlers, rotas e middlewares) é realizada em
internal/app/application.go.
*/

// Package main Desafio Go API
//
// @title           Desafio Go API
// @version         1.0
// @description     API REST desenvolvida em Go utilizando Clean Architecture.
// @termsOfService  http://swagger.io/terms/
//
// @contact.name   Bianca Salomão
// @contact.url    https://github.com/BiancaSalomao1
//
// @license.name  MIT
//
// @host      localhost:8080
// @BasePath  /
//
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
package main

import (
	"log/slog"
	"net/http"

	"desafio-go/infrastructure/logging"

	_ "desafio-go/docs"

	"desafio-go/config"
	"desafio-go/infrastructure/database"
	app "desafio-go/internal/app"
)

func main() {

	logger := logging.New()

	slog.SetDefault(logger)

	// Carrega as configurações da aplicação
	cfg := config.Load()

	// Executa as migrations
	if err := database.RunMigrations(cfg.DatabaseURL); err != nil {
		slog.Error(
			"failed to run migrations",
			"error", err,
		)
		return
	}

	// Abre conexão com o banco
	db, err := database.New(cfg)
	if err != nil {
		slog.Error("failed to connect to database", "error", err)
	}
	defer db.Close()

	// Monta toda a aplicação
	httpHandler, err := app.NewApplication(cfg, db)
	if err != nil {
		slog.Error("failed to create application", "error", err)
	}

	// Configura servidor HTTP
	server := &http.Server{
		Addr:    ":" + cfg.AppPort,
		Handler: httpHandler,
	}

	slog.Info(
		"server started",
		"port", cfg.AppPort,
		"environment", cfg.AppEnv,
	)
	// Inicia servidor
	if err := server.ListenAndServe(); err != nil {
		slog.Error("failed to start server", "error", err)
	}
}
