package handler

import (
	"desafio-go/internal/domain"
)

type CreateCustomerUseCase interface {
	Execute(customer *domain.Customer) error
}

type GetCustomerUseCase interface {
	Execute(id string) (*domain.Customer, error)
}

type ListCustomersUseCase interface {
	Execute() ([]*domain.Customer, error)
}

type UpdateCustomerUseCase interface {
	Execute(*domain.Customer) error
}

type DeleteCustomerUseCase interface {
	Execute(string) error
}
