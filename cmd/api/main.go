package main

import (
	"fmt"

	"desafio-go/infrastructure/repository/memory"

	customerusecase "desafio-go/internal/usecase/customer"
	orderusecase "desafio-go/internal/usecase/order"
	productusecase "desafio-go/internal/usecase/product"
	userusecase "desafio-go/internal/usecase/user"
)

func main() {

	fmt.Println("===================================")
	fmt.Println("      SERVIÇO DE PEDIDOS")
	fmt.Println("===================================")

	//----------------------------------
	// Repositories
	//----------------------------------

	productRepository := memory.NewMemoryProductRepository()
	customerRepository := memory.NewMemoryCustomerRepository()
	userRepository := memory.NewMemoryUserRepository()
	orderRepository := memory.NewMemoryOrderRepository()

	//----------------------------------
	// Product UseCases
	//----------------------------------

	createProduct := productusecase.NewCreateProductUseCase(productRepository)
	getProduct := productusecase.NewGetProductUseCase(productRepository)
	listProducts := productusecase.NewListProductsUseCase(productRepository)

	_ = createProduct
	_ = getProduct
	_ = listProducts

	//----------------------------------
	// Customer UseCases
	//----------------------------------

	createCustomer := customerusecase.NewCreateCustomerUseCase(customerRepository)
	getCustomer := customerusecase.NewGetCustomerUseCase(customerRepository)

	_ = createCustomer
	_ = getCustomer

	//----------------------------------
	// User UseCases
	//----------------------------------

	createUser := userusecase.NewCreateUserUseCase(userRepository)

	_ = createUser

	//----------------------------------
	// Order UseCases
	//----------------------------------

	createOrder := orderusecase.NewCreateOrderUseCase(
		orderRepository,
		productRepository,
		customerRepository,
	)

	getOrder := orderusecase.NewGetOrderUseCase(orderRepository)

	payOrder := orderusecase.NewPayOrderUseCase(orderRepository)

	cancelOrder := orderusecase.NewCancelOrderUseCase(
		orderRepository,
		productRepository,
	)

	_ = createOrder
	_ = getOrder
	_ = payOrder
	_ = cancelOrder

	fmt.Println("Arquitetura inicializada com sucesso.")
}
