package integration

/*
helper_test

Responsabilidades:

- carregar as configurações da aplicação;
- abrir conexão com o banco de testes;
- conectar ao RabbitMQ;
- criar publisher e consumer;
- montar toda a aplicação utilizando app.NewApplication;
- iniciar um servidor HTTP de testes;
- disponibilizar funções auxiliares para os testes de integração.

Este arquivo NÃO testa regras de negócio.

Seu objetivo é preparar um ambiente completo para execução dos
testes de integração utilizando a mesma configuração da aplicação.

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

Saga:

Orders API
      ↓
RabbitMQ Publisher
      ↓
Product Service
      ↓
RabbitMQ
      ↓
Orders API Consumer
*/

import (
	"context"
	"log/slog"
	"net/http/httptest"
	"testing"

	"github.com/joho/godotenv"

	"orders-api/config"
	"orders-api/infrastructure/database"
	"orders-api/infrastructure/messaging/rabbitmq"
	app "orders-api/internal/app"
	"orders-api/internal/security"
)

type TestServer struct {
	Server     *httptest.Server
	DB         *database.Database
	Config     *config.Config
	RabbitConn *rabbitmq.Connection
	Publisher  *rabbitmq.Publisher
	Consumer   *rabbitmq.Consumer
	Cancel     context.CancelFunc
}

func setup(t *testing.T) *TestServer {
	t.Helper()

	_ = godotenv.Load("../../.env")

	cfg := config.Load()

	// Executa as migrations.
	if err := database.RunMigrations(cfg.DatabaseURL); err != nil {
		t.Fatalf("error running migrations: %v", err)
	}

	// Conecta ao banco.
	db, err := database.New(cfg)
	if err != nil {
		t.Fatalf("error connecting database: %v", err)
	}

	// Limpa o banco antes de cada teste.
	cleanDatabase(t, db)

	// Conecta ao RabbitMQ.
	rabbitConn, err := rabbitmq.NewConnection(cfg.RabbitMQURL)
	if err != nil {
		db.Close()
		t.Fatalf("error connecting RabbitMQ: %v", err)
	}

	// Cria o publisher.
	publisher, err := rabbitmq.NewPublisher(
		rabbitConn,
		slog.Default(),
	)
	if err != nil {
		rabbitConn.Close()
		db.Close()

		t.Fatalf("error creating RabbitMQ publisher: %v", err)
	}
	tokenStore := security.NewMemoryTokenStore()

	// Monta a aplicação e obtém o handler de eventos da Saga.
	httpHandler, orderEventHandler, err := app.NewApplication(
		cfg,
		db,
		tokenStore,
		publisher,
	)

	if err != nil {
		publisher.Close()
		rabbitConn.Close()
		db.Close()

		t.Fatalf("error creating application: %v", err)
	}

	// Cria o consumer do Orders API.
	consumer, err := rabbitmq.NewConsumer(
		rabbitConn,
		slog.Default(),
	)
	if err != nil {
		publisher.Close()
		rabbitConn.Close()
		db.Close()

		t.Fatalf("error creating RabbitMQ consumer: %v", err)
	}

	// Inicia o consumer da Saga.
	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		if err := consumer.Consume(
			ctx,
			orderEventHandler.Handle,
		); err != nil {
			// O cancelamento normal do contexto não é erro de teste.
			if ctx.Err() == nil {
				slog.Error(
					"test RabbitMQ consumer stopped",
					"error", err,
				)
			}
		}
	}()

	server := httptest.NewServer(httpHandler)

	return &TestServer{
		Server:     server,
		DB:         db,
		Config:     cfg,
		RabbitConn: rabbitConn,
		Publisher:  publisher,
		Consumer:   consumer,
		Cancel:     cancel,
	}
}

func teardown(ts *TestServer) {
	if ts == nil {
		return
	}

	// Para primeiro o consumer.
	if ts.Cancel != nil {
		ts.Cancel()
	}

	if ts.Server != nil {
		ts.Server.Close()
	}

	if ts.Consumer != nil {
		ts.Consumer.Close()
	}

	if ts.Publisher != nil {
		ts.Publisher.Close()
	}

	if ts.RabbitConn != nil {
		ts.RabbitConn.Close()
	}

	if ts.DB != nil {
		ts.DB.Close()
	}
}
