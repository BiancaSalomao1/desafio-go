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
	listCustomers := customerusecase.NewListCustomersUseCase(customerRepository)

	_ = createCustomer
	_ = getCustomer
	_ = listCustomers

	// User UseCases

	createUser := userusecase.NewCreateUserUseCase(userRepository)
	getUser := userusecase.NewGetUserUseCase(userRepository)
	listUsers := userusecase.NewListUsersUseCase(userRepository)

	_ = createUser
	_ = getUser
	_ = listUsers
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

	// Cadastro de Usuário

	fmt.Println()
	fmt.Println(">>> Cadastro de Usuário")

	user := domain.NewUser(
		"U001",
		"Administrador",
		"admin@email.com",
		"123456",
	)

	if err := createUser.Execute(user); err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("Usuário cadastrado com sucesso.")

	fmt.Println()
	fmt.Println(">>> Lista de Usuários")

	users, err := listUsers.Execute()
	if err != nil {
		fmt.Println(err)
		return
	}

	for _, user := range users {

		fmt.Printf(
			"ID: %-5s Nome: %-15s Email: %s\n",
			user.ID,
			user.Name,
			user.Email,
		)
	}
	fmt.Println()
	fmt.Println(">>> Buscar Usuário")

	savedUser, err := getUser.Execute("U001")
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf(
		"Usuário: %s | %s\n",
		savedUser.Name,
		savedUser.Email,
	)

	//----------------------------------
	// Cadastro de Cliente
	//----------------------------------

	fmt.Println()
	fmt.Println(">>> Cadastro de Cliente")

	customer := domain.NewCustomer(
		"C001",
		"Ana",
		"ana@email.com",
	)

	if err := createCustomer.Execute(customer); err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("Cliente cadastrado com sucesso.")

	fmt.Println()
	fmt.Println(">>> Lista de Clientes")

	customers, err := listCustomers.Execute()
	if err != nil {
		fmt.Println(err)
		return
	}

	for _, customer := range customers {

		fmt.Printf(
			"ID: %-5s Nome: %-10s Email: %s\n",
			customer.ID,
			customer.Name,
			customer.Email,
		)
	}
	fmt.Println()
	fmt.Println(">>> Buscar Cliente")

	savedCustomer, err := getCustomer.Execute("C001")
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf(
		"Cliente: %s | %s\n",
		savedCustomer.Name,
		savedCustomer.Email,
	)

}
