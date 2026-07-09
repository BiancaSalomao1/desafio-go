package customer

/*
struct CreateCustomerUseCase
- criar cliente;
- validar dados;
- salvar cliente.

Métodos:
- NewCreateCustomerUseCase()
- Execute()
*/

import (
	"desafio-go/internal/domain"
	"desafio-go/internal/repository"
)

type CreateCustomerUseCase struct {
	customerRepository repository.CustomerRepository
}

func NewCreateCustomerUseCase(customerRepository repository.CustomerRepository) *CreateCustomerUseCase {
	return &CreateCustomerUseCase{
		customerRepository: customerRepository,
	}
}

func (uc *CreateCustomerUseCase) Execute(customer *domain.Customer) error {

	if err := customer.Validate(); err != nil {
		return err
	}

	return uc.customerRepository.Save(customer)
}
