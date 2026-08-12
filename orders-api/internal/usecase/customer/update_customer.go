package customer

/*
struct UpdateCustomerUseCase

Responsabilidades:
- localizar um cliente;
- atualizar seus dados;
- persistir a alteração.

Campos:
- customerRepository

Métodos:
- NewUpdateCustomerUseCase()
- Execute()
*/

import (
	"desafio-go/orders-api/internal/domain"
	"desafio-go/orders-api/internal/repository"
)

type UpdateCustomerUseCase struct {
	customerRepository repository.CustomerRepository
}

func NewUpdateCustomerUseCase(
	customerRepository repository.CustomerRepository,
) *UpdateCustomerUseCase {

	return &UpdateCustomerUseCase{
		customerRepository: customerRepository,
	}
}

func (uc *UpdateCustomerUseCase) Execute(
	customer *domain.Customer,
) error {

	currentCustomer, err := uc.customerRepository.FindByID(customer.ID)
	if err != nil {
		return err
	}

	if err := currentCustomer.Update(
		customer.Name,
		customer.Email,
	); err != nil {
		return err
	}

	return uc.customerRepository.Update(currentCustomer)
}
