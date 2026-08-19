package integration

/*
helper_test

Responsabilidades:

- carregar as configurações da aplicação;
- abrir conexão com o banco de testes;
- montar toda a aplicação utilizando app.NewApplication;
- iniciar um servidor HTTP de testes;
- disponibilizar funções auxiliares para os testes de integração.

Este arquivo NÃO testa regras de negócio.

Seu objetivo é preparar um ambiente completo para execução dos
testes de integração utilizando exatamente a mesma configuração da
aplicação em produção.

Fluxo:

HTTP Request
      ↓
Middlewares
      ↓
Router
      ↓
Handler
      ↓
Use Case
      ↓
Repository
      ↓
PostgreSQL
*/

import (
	"net/http/httptest"
	"testing"

	"github.com/joho/godotenv"

	"orders-api/config"
	"orders-api/infrastructure/database"
	app "orders-api/internal/app"
)

type TestServer struct {
	Server *httptest.Server
	DB     *database.Database
	Config *config.Config
}

func setup(t *testing.T) *TestServer {
	t.Helper()

	if err := godotenv.Load("../../.env"); err != nil {
		t.Fatalf("error loading .env: %v", err)
	}

	cfg := config.Load()

	// Executa as migrations
	if err := database.RunMigrations(cfg.DatabaseURL); err != nil {
		t.Fatalf("error running migrations: %v", err)
	}

	// Conecta ao banco
	db, err := database.New(cfg)
	if err != nil {
		t.Fatalf("error connecting database: %v", err)
	}

	// Limpa o banco antes do teste
	cleanDatabase(t, db)

	// Monta a aplicação
	httpHandler, _, err := app.NewApplication(cfg, db)
	if err != nil {
		db.Close()
		t.Fatalf("error creating application: %v", err)
	}

	server := httptest.NewServer(httpHandler)

	return &TestServer{
		Server: server,
		DB:     db,
		Config: cfg,
	}
}

func teardown(ts *TestServer) {
	ts.Server.Close()
	ts.DB.Close()
}
