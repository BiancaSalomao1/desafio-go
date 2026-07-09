package customer

/*
struct ListCustomersUseCase
- listar clientes.

Métodos:
- NewListCustomersUseCase()
- Execute()
*/

import (
	"desafio-go/internal/domain"
	"desafio-go/internal/repository"
)

type ListCustomersUseCase struct {
	customerRepository repository.CustomerRepository
}

func NewListCustomersUseCase(customerRepository repository.CustomerRepository) *ListCustomersUseCase {
	return &ListCustomersUseCase{
		customerRepository: customerRepository,
	}
}

func (uc *ListCustomersUseCase) Execute() ([]*domain.Customer, error) {
	return uc.customerRepository.FindAll()
}
