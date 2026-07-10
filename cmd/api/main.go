package main

import (
	"fmt"

	"desafio-go/config"
	"desafio-go/infrastructure/database"
	"desafio-go/infrastructure/repository/postgres"
)

func main() {

	cfg := config.Load()

	db, err := database.New(cfg)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	fmt.Println("===================================")
	fmt.Println(" SERVIÇO DE PEDIDOS")
	fmt.Println("===================================")

	fmt.Println("Conexão com PostgreSQL estabelecida.")

	factory := postgres.NewRepositoryFactory()

	productRepository := factory.Product(db.Pool)
	customerRepository := factory.Customer(db.Pool)
	userRepository := factory.User(db.Pool)
	orderRepository := factory.Order(db.Pool)

	fmt.Printf("ProductRepository  : %T\n", productRepository)
	fmt.Printf("CustomerRepository : %T\n", customerRepository)
	fmt.Printf("UserRepository     : %T\n", userRepository)
	fmt.Printf("OrderRepository    : %T\n", orderRepository)

	fmt.Println()
	fmt.Println("Infraestrutura carregada com sucesso.")
}
