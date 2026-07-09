package main

import (
	"fmt"

	"desafio-go/infrastructure/repository/memory"

	"desafio-go/internal/domain"
	customerusecase "desafio-go/internal/usecase/customer"
	orderusecase "desafio-go/internal/usecase/order"
	productusecase "desafio-go/internal/usecase/product"
	userusecase "desafio-go/internal/usecase/user"
)

func main() {

	fmt.Println("      SERVIÇO DE PEDIDOS")

	// Repositories

	productRepository := memory.NewMemoryProductRepository()
	customerRepository := memory.NewMemoryCustomerRepository()
	userRepository := memory.NewMemoryUserRepository()
	orderRepository := memory.NewMemoryOrderRepository()

	// Product UseCases

	createProduct := productusecase.NewCreateProductUseCase(productRepository)
	getProduct := productusecase.NewGetProductUseCase(productRepository)
	listProducts := productusecase.NewListProductsUseCase(productRepository)

	_ = createProduct
	_ = getProduct
	_ = listProducts

	// Customer UseCases

	createCustomer := customerusecase.NewCreateCustomerUseCase(customerRepository)
	getCustomer := customerusecase.NewGetCustomerUseCase(customerRepository)

	_ = createCustomer
	_ = getCustomer

	// User UseCases

	createUser := userusecase.NewCreateUserUseCase(userRepository)

	_ = createUser

	// Order UseCases

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

	// Cadastro de Produtos

	fmt.Println()
	fmt.Println(">>> Cadastro de Produtos")

	notebook := domain.NewProduct(
		"P001",
		"Notebook",
		3500.00,
		5,
	)

	mouse := domain.NewProduct(
		"P002",
		"Mouse",
		80.00,
		10,
	)

	teclado := domain.NewProduct(
		"P003",
		"Teclado",
		180.00,
		8,
	)

	if err := createProduct.Execute(notebook); err != nil {
		fmt.Println(err)
	}

	if err := createProduct.Execute(mouse); err != nil {
		fmt.Println(err)
	}

	if err := createProduct.Execute(teclado); err != nil {
		fmt.Println(err)
	}

	fmt.Println("Produtos cadastrados com sucesso.")

	fmt.Println()
	fmt.Println(">>> Lista de Produtos")

	products, err := listProducts.Execute()
	if err != nil {
		fmt.Println(err)
		return
	}

	for _, product := range products {

		fmt.Printf(
			"ID: %-5s Nome: %-10s Preço: R$ %7.2f Estoque: %d\n",
			product.ID,
			product.Name,
			product.Price,
			product.Stock,
		)
	}

	fmt.Println()
	fmt.Println(">>> Buscar Produto")

	product, err := getProduct.Execute("P001")
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf(
		"Produto encontrado: %s | R$ %.2f | Estoque: %d\n",
		product.Name,
		product.Price,
		product.Stock,
	)

}
