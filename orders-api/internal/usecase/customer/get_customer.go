package customer

/*
struct GetCustomerUseCase
- buscar cliente por ID.

Métodos:
- NewGetCustomerUseCase()
- Execute()
*/

import (
	"desafio-go/orders-api/internal/domain"
	"desafio-go/orders-api/internal/repository"
)

type GetCustomerUseCase struct {
	customerRepository repository.CustomerRepository
}

func NewGetCustomerUseCase(customerRepository repository.CustomerRepository) *GetCustomerUseCase {
	return &GetCustomerUseCase{
		customerRepository: customerRepository,
	}
}

func (uc *GetCustomerUseCase) Execute(id string) (*domain.Customer, error) {
	return uc.customerRepository.FindByID(id)
}
