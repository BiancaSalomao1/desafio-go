package main

import (
	"fmt"

	"github.com/google/uuid"

	"desafio-go/config"
	"desafio-go/infrastructure/database"
	"desafio-go/infrastructure/repository/postgres"

	"desafio-go/internal/domain"
	productusecase "desafio-go/internal/usecase/product"
)

func main() {

	// ==========================================
	// Configuração
	// ==========================================

	cfg := config.Load()

	// ==========================================
	// Migrations
	// ==========================================

	if err := database.RunMigrations(cfg.DatabaseURL); err != nil {
		panic(err)
	}

	// ==========================================
	// Banco de Dados
	// ==========================================

	db, err := database.New(cfg)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	fmt.Println("===================================")
	fmt.Println(" SERVIÇO DE PEDIDOS")
	fmt.Println("===================================")
	fmt.Println("Conectado ao PostgreSQL com sucesso.")

	// ==========================================
	// Repository
	// ==========================================

	productRepository := postgres.NewProductRepository(db)

	// ==========================================
	// Use Cases
	// ==========================================

	createProduct := productusecase.NewCreateProductUseCase(productRepository)
	getProduct := productusecase.NewGetProductUseCase(productRepository)
	listProducts := productusecase.NewListProductsUseCase(productRepository)

	// ==========================================
	// Cadastro
	// ==========================================

	fmt.Println()
	fmt.Println(">>> Cadastro de Produtos")

	notebook := domain.NewProduct(
		uuid.NewString(),
		"Notebook",
		3500.00,
		5,
	)

	mouse := domain.NewProduct(
		uuid.NewString(),
		"Mouse",
		80.00,
		10,
	)

	teclado := domain.NewProduct(
		uuid.NewString(),
		"Teclado",
		180.00,
		8,
	)

	products := []*domain.Product{
		notebook,
		mouse,
		teclado,
	}
	for _, product := range products {

		if err := createProduct.Execute(product); err != nil {
			fmt.Println(err)
			continue
		}

		fmt.Printf("Produto %s cadastrado.\n", product.Name)
	}

	// ==========================================
	// Listagem
	// ==========================================

	fmt.Println()
	fmt.Println(">>> Lista de Produtos")

	list, err := listProducts.Execute()
	if err != nil {
		panic(err)
	}

	for _, product := range list {

		fmt.Printf(
			"ID: %-5s Nome: %-10s Preço: R$ %8.2f Estoque: %d\n",
			product.ID,
			product.Name,
			product.Price,
			product.Stock,
		)
	}

	// ==========================================
	// Busca
	// ==========================================

	fmt.Println()
	fmt.Println(">>> Buscar Produto")

	product, err := getProduct.Execute(notebook.ID)
	if err != nil {
		panic(err)
	}

	fmt.Printf(
		"Produto encontrado: %s | R$ %.2f | Estoque: %d\n",
		product.Name,
		product.Price,
		product.Stock,
	)

	fmt.Println()
	fmt.Println("Teste do ProductRepository finalizado com sucesso.")

}
