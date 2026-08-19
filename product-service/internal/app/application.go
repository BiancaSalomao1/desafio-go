package app

import (
	"context"
	"log/slog"
	"net/http"

	"product-service/config"
	"product-service/infrastructure/database"
	"product-service/infrastructure/http/handler"
	"product-service/infrastructure/http/routes"
	"product-service/infrastructure/messaging/rabbitmq"
	"product-service/infrastructure/repository/postgres"
	"product-service/internal/messaging"
	productusecase "product-service/internal/usecase/product"
)

// NewApplication monta toda a aplicação do Product Service.
//
// Fluxo:
//
// HTTP
// ↓
// Router
// ↓
// Handler
// ↓
// Use Cases
// ↓
// ProductRepository
// ↓
// PostgreSQL
//
// Fluxo da Saga:
//
// RabbitMQ
// ↓
// Consumer
// ↓
// StockHandler
// ↓
// ReserveStockUseCase / ReleaseStockUseCase
// ↓
// ProductRepository
// ↓
// PostgreSQL
// ↓
// Publisher
// ↓
// RabbitMQ
func NewApplication(
	cfg *config.Config,
	db *database.Database,
) (http.Handler, error) {

	logger := slog.Default()

	// ------------------------------------------------------------
	// Repository
	// ------------------------------------------------------------

	productRepository := postgres.NewProductRepository(
		db.Pool,
	)

	// ------------------------------------------------------------
	// Product Use Cases
	// ------------------------------------------------------------

	createProduct := productusecase.NewCreateProductUseCase(
		productRepository,
	)

	getProduct := productusecase.NewGetProductUseCase(
		productRepository,
	)

	listProducts := productusecase.NewListProductsUseCase(
		productRepository,
	)

	updateProduct := productusecase.NewUpdateProductUseCase(
		productRepository,
	)

	deleteProduct := productusecase.NewDeleteProductUseCase(
		productRepository,
	)

	// ------------------------------------------------------------
	// Stock Use Cases
	// ------------------------------------------------------------

	reserveStock := productusecase.NewReserveStockUseCase(
		productRepository,
	)

	releaseStock := productusecase.NewReleaseStockUseCase(
		productRepository,
	)

	// ------------------------------------------------------------
	// RabbitMQ Connection
	// ------------------------------------------------------------

	rabbitConnection, err := rabbitmq.NewConnection(
		cfg.RabbitMQURL,
	)
	if err != nil {
		return nil, err
	}

	// ------------------------------------------------------------
	// RabbitMQ Publisher
	// ------------------------------------------------------------

	publisher, err := rabbitmq.NewPublisher(
		rabbitConnection,
		logger,
	)
	if err != nil {
		rabbitConnection.Close()
		return nil, err
	}

	// ------------------------------------------------------------
	// Stock Handler
	// ------------------------------------------------------------

	stockHandler := messaging.NewStockHandler(
		logger,
		reserveStock,
		releaseStock,
		publisher,
	)

	// ------------------------------------------------------------
	// RabbitMQ Consumer
	// ------------------------------------------------------------

	consumer, err := rabbitmq.NewConsumer(
		rabbitConnection,
		logger,
	)
	if err != nil {
		publisher.Close()
		rabbitConnection.Close()

		return nil, err
	}

	// ------------------------------------------------------------
	// Start RabbitMQ Consumer
	// ------------------------------------------------------------

	go func() {
		ctx := context.Background()

		if err := consumer.Consume(
			ctx,
			stockHandler.Handle,
		); err != nil {
			logger.Error(
				"rabbitmq consumer stopped",
				"service", "product-service",
				"error", err,
			)
		}
	}()

	// ------------------------------------------------------------
	// HTTP Handler
	// ------------------------------------------------------------

	productHandler := handler.NewProductHandler(
		createProduct,
		getProduct,
		listProducts,
		updateProduct,
		deleteProduct,
	)

	return routes.NewRouter(productHandler), nil
}
