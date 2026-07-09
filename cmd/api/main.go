package main

import (
	"fmt"

	"desafio-go/config"

	"desafio-go/infrastructure/repository/memory"

	"desafio-go/internal/domain"
	customerusecase "desafio-go/internal/usecase/customer"
	orderusecase "desafio-go/internal/usecase/order"
	productusecase "desafio-go/internal/usecase/product"
	userusecase "desafio-go/internal/usecase/user"
)

func main() {

	cfg := config.Load()

	fmt.Println(cfg.AppName)
	fmt.Println(cfg.DatabaseURL)

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

	fmt.Println()
	fmt.Println(">>> Criação de Pedido")

	order := domain.NewOrder(
		"PED-001",
		customer.ID,
	)

	err = order.AddItem(
		*domain.NewOrderItem(
			"P001",
			"",
			0,
			1,
		),
	)
	if err != nil {
		fmt.Println(err)
		return
	}

	err = order.AddItem(
		*domain.NewOrderItem(
			"P002",
			"",
			0,
			2,
		),
	)
	if err != nil {
		fmt.Println(err)
		return
	}

	if err := createOrder.Execute(order); err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("Pedido criado com sucesso.")

	fmt.Println()

	savedOrder, err := getOrder.Execute("PED-001")
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("Pedido........: %s\n", savedOrder.ID)
	fmt.Printf("Cliente.......: %s\n", savedOrder.CustomerID)
	fmt.Printf("Status........: %s\n", savedOrder.Status)
	fmt.Printf("Itens.........: %d\n", len(savedOrder.Items))
	fmt.Printf("Total.........: %.2f\n", savedOrder.Total())

	fmt.Println()
	fmt.Println(">>> Estoque após pedido")

	products, _ = listProducts.Execute()

	for _, product := range products {

		fmt.Printf(
			"%s - Estoque: %d\n",
			product.Name,
			product.Stock,
		)
	}

	fmt.Println()
	fmt.Println(">>> Pagamento")

	if err := payOrder.Execute(order.ID); err != nil {
		fmt.Println(err)
		return
	}

	savedOrder, _ = getOrder.Execute(order.ID)

	fmt.Printf("Novo status: %s\n", savedOrder.Status)

	fmt.Println()
	fmt.Println(">>> Cancelando pedido pago")

	if err := cancelOrder.Execute(order.ID); err != nil {
		fmt.Println(err)
	}

	fmt.Println()
	fmt.Println(">>> Pedido com estoque insuficiente")

	order2 := domain.NewOrder(
		"PED-002",
		customer.ID,
	)

	order2.AddItem(
		*domain.NewOrderItem(
			"P001",
			"",
			0,
			100,
		),
	)

	if err := createOrder.Execute(order2); err != nil {
		fmt.Println(err)
	}

	fmt.Println()
	fmt.Println(">>> Pedido vazio")

	order3 := domain.NewOrder(
		"PED-003",
		customer.ID,
	)

	if err := createOrder.Execute(order3); err != nil {
		fmt.Println(err)
	}

	fmt.Println()
	fmt.Println(">>> Cliente inválido")

	order4 := domain.NewOrder(
		"PED-004",
		"",
	)

	order4.AddItem(
		*domain.NewOrderItem(
			"P002",
			"",
			0,
			1,
		),
	)

	if err := createOrder.Execute(order4); err != nil {
		fmt.Println(err)
	}
}
