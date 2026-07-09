package customer

/*
struct UpdateCustomerUseCase
- atualizar cliente.

Métodos:
- NewUpdateCustomerUseCase()
- Execute()
*/

import (
	"desafio-go/internal/domain"
	"desafio-go/internal/repository"
)

type UpdateCustomerUseCase struct {
	customerRepository repository.CustomerRepository
}

func NewUpdateCustomerUseCase(customerRepository repository.CustomerRepository) *UpdateCustomerUseCase {
	return &UpdateCustomerUseCase{
		customerRepository: customerRepository,
	}
}

func (uc *UpdateCustomerUseCase) Execute(customer *domain.Customer) error {

	if err := customer.Validate(); err != nil {
		return err
	}

	return uc.customerRepository.Update(customer)
}
