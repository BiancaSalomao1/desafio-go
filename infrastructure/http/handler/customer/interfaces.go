package handler

/*
interface CreateCustomerUseCase

Responsabilidades:
- criar um cliente.

Métodos:
- Execute()
*/

import (
	"desafio-go/internal/domain"
)

type CreateCustomerUseCase interface {
	Execute(customer *domain.Customer) error
}

/*
interface GetCustomerUseCase

Responsabilidades:
- buscar um cliente por ID.

Métodos:
- Execute()
*/

type GetCustomerUseCase interface {
	Execute(id string) (*domain.Customer, error)
}

/*
interface ListCustomersUseCase

Responsabilidades:
- listar clientes.

Métodos:
- Execute()
*/

type ListCustomersUseCase interface {
	Execute() ([]*domain.Customer, error)
}
