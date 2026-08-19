/*
application

Responsabilidades:

- montar toda a aplicação;
- criar repositórios;
- criar casos de uso;
- criar handlers;
- registrar rotas;
- configurar middlewares.

Este arquivo NÃO inicia o servidor HTTP.

Seu objetivo é disponibilizar uma aplicação totalmente configurada,
permitindo que tanto o main.go quanto os testes de integração utilizem
a mesma configuração.

Fluxo:

HTTP Request
      ↓
Middlewares
      ↓
Router
      ↓
Handlers
      ↓
Use Cases
      ↓
Repositories
      ↓
PostgreSQL
*/

package app

import (
	"log/slog"
	"net/http"
	"orders-api/config"
	"orders-api/infrastructure/database"
	"orders-api/infrastructure/http/handler"
	"orders-api/infrastructure/http/middleware"
	"orders-api/infrastructure/http/routes"
	"orders-api/infrastructure/repository/postgres"
	"time"

	authusecase "orders-api/internal/usecase/auth"
	customerusecase "orders-api/internal/usecase/customer"
	orderusecase "orders-api/internal/usecase/order"
	productusecase "orders-api/internal/usecase/product"
	userusecase "orders-api/internal/usecase/user"

	"orders-api/infrastructure/messaging/rabbitmq"
	"orders-api/internal/messaging"
)

// NewApplication monta toda a aplicação e retorna o handler HTTP
// configurado com todas as dependências.
func NewApplication(
	cfg *config.Config,
	db *database.Database,
	publishers ...messaging.EventPublisher,
) (http.Handler, *rabbitmq.OrderEventHandler, error) {

	transactionManager := database.NewTransactionManager(db.Pool)

	repositoryFactory := postgres.NewRepositoryFactory()

	productRepository := repositoryFactory.Product(db.Pool)
	customerRepository := repositoryFactory.Customer(db.Pool)
	userRepository := repositoryFactory.User(db.Pool)
	orderRepository := repositoryFactory.Order(db.Pool)

	createProduct := productusecase.NewCreateProductUseCase(productRepository)
	getProduct := productusecase.NewGetProductUseCase(productRepository)
	listProducts := productusecase.NewListProductsUseCase(productRepository)
	updateProduct := productusecase.NewUpdateProductUseCase(productRepository)
	deleteProduct := productusecase.NewDeleteProductUseCase(productRepository)

	createCustomer := customerusecase.NewCreateCustomerUseCase(customerRepository)
	getCustomer := customerusecase.NewGetCustomerUseCase(customerRepository)
	listCustomers := customerusecase.NewListCustomersUseCase(customerRepository)
	updateCustomer := customerusecase.NewUpdateCustomerUseCase(customerRepository)
	deleteCustomer := customerusecase.NewDeleteCustomerUseCase(customerRepository)

	createUser := userusecase.NewCreateUserUseCase(userRepository)
	getUser := userusecase.NewGetUserUseCase(userRepository)
	listUsers := userusecase.NewListUsersUseCase(userRepository)
	updateUser := userusecase.NewUpdateUserUseCase(userRepository)
	deleteUser := userusecase.NewDeleteUserUseCase(userRepository)

	createOrder := orderusecase.NewCreateOrderUseCase(
		transactionManager,
		repositoryFactory,
		publishers...,
	)
	handleStockReserved := orderusecase.NewHandleStockReservedUseCase(
		orderRepository,
	)

	handleStockReservationFailed := orderusecase.NewHandleStockReservationFailedUseCase(
		orderRepository,
	)

	orderEventHandler := rabbitmq.NewOrderEventHandler(
		handleStockReserved,
		handleStockReservationFailed,
		slog.Default(),
	)

	getOrder := orderusecase.NewGetOrderUseCase(orderRepository)

	listOrders := orderusecase.NewListOrdersUseCase(orderRepository)

	payOrder := orderusecase.NewPayOrderUseCase(orderRepository)

	cancelOrder := orderusecase.NewCancelOrderUseCase(
		transactionManager,
		repositoryFactory,
		publishers...,
	)

	jwtTTL, err := time.ParseDuration(cfg.JWTExpiresIn)
	if err != nil {
		return nil, nil, err
	}

	loginUseCase := authusecase.NewLoginUseCase(
		userRepository,
		cfg.JWTSecret,
		jwtTTL,
	)

	authHandler := handler.NewAuthHandler(
		loginUseCase,
	)

	productHandler := handler.NewProductHandler(
		createProduct,
		getProduct,
		listProducts,
		updateProduct,
		deleteProduct,
	)

	customerHandler := handler.NewCustomerHandler(
		createCustomer,
		getCustomer,
		listCustomers,
		updateCustomer,
		deleteCustomer,
	)

	userHandler := handler.NewUserHandler(
		createUser,
		getUser,
		listUsers,
		updateUser,
		deleteUser,
	)

	orderHandler := handler.NewOrderHandler(
		createOrder,
		getOrder,
		listOrders,
		payOrder,
		cancelOrder,
	)

	router := routes.NewRouter(
		productHandler,
		customerHandler,
		userHandler,
		orderHandler,
		authHandler,
		cfg.JWTSecret,
	)

	httpHandler := middleware.Recovery(
		middleware.Logger(
			middleware.CORS(router),
		),
	)

	return httpHandler, orderEventHandler, nil
}
