package main

/*
Função main

Responsabilidades:
- carregar configurações;
- executar migrations;
- abrir conexão com PostgreSQL;
- montar dependências da aplicação;
- registrar rotas;
- configurar middlewares;
- iniciar servidor HTTP.

Métodos:
- main()
*/

import (
	"log"
	"net/http"
	"time"

	"desafio-go/config"

	"desafio-go/infrastructure/database"

	"desafio-go/infrastructure/http/handler"
	"desafio-go/infrastructure/http/middleware"
	"desafio-go/infrastructure/http/routes"

	"desafio-go/infrastructure/repository/postgres"

	authusecase "desafio-go/internal/usecase/auth"
	customerusecase "desafio-go/internal/usecase/customer"
	orderusecase "desafio-go/internal/usecase/order"
	productusecase "desafio-go/internal/usecase/product"
	userusecase "desafio-go/internal/usecase/user"
)

func main() {

	cfg := config.Load()

	if err := database.RunMigrations(cfg.DatabaseURL); err != nil {
		log.Fatal(err)
	}

	db, err := database.New(cfg)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	transactionManager := database.NewTransactionManager(db.Pool)

	repositoryFactory := postgres.NewRepositoryFactory()

	// Repositórios
	productRepository := repositoryFactory.Product(db.Pool)
	customerRepository := repositoryFactory.Customer(db.Pool)
	userRepository := repositoryFactory.User(db.Pool)
	orderRepository := repositoryFactory.Order(db.Pool)

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

	createCustomer := customerusecase.NewCreateCustomerUseCase(
		customerRepository,
	)

	getCustomer := customerusecase.NewGetCustomerUseCase(
		customerRepository,
	)

	listCustomers := customerusecase.NewListCustomersUseCase(
		customerRepository,
	)

	updateCustomer := customerusecase.NewUpdateCustomerUseCase(
		customerRepository,
	)

	deleteCustomer := customerusecase.NewDeleteCustomerUseCase(
		customerRepository,
	)

	createUser := userusecase.NewCreateUserUseCase(
		userRepository,
	)

	getUser := userusecase.NewGetUserUseCase(
		userRepository,
	)

	listUsers := userusecase.NewListUsersUseCase(
		userRepository,
	)

	updateUser := userusecase.NewUpdateUserUseCase(
		userRepository,
	)

	deleteUser := userusecase.NewDeleteUserUseCase(
		userRepository,
	)

	createOrder := orderusecase.NewCreateOrderUseCase(
		transactionManager,
		repositoryFactory,
	)

	getOrder := orderusecase.NewGetOrderUseCase(
		orderRepository,
	)

	listOrders := orderusecase.NewListOrdersUseCase(
		orderRepository,
	)

	payOrder := orderusecase.NewPayOrderUseCase(
		orderRepository,
	)

	cancelOrder := orderusecase.NewCancelOrderUseCase(
		transactionManager,
		repositoryFactory,
	)

	// JWT
	jwtTTL, err := time.ParseDuration(cfg.JWTExpiresIn)
	if err != nil {
		log.Fatal(err)
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
	)

	httpHandler := middleware.Recovery(
		middleware.Logger(
			middleware.CORS(router),
		),
	)

	server := &http.Server{
		Addr:    ":" + cfg.AppPort,
		Handler: httpHandler,
	}

	log.Printf("Servidor iniciado em %s", cfg.AppPort)

	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
