package customer

/*
struct DeleteCustomerUseCase
- excluir cliente.

Métodos:
- NewDeleteCustomerUseCase()
- Execute()
*/

import "desafio-go/internal/repository"

type DeleteCustomerUseCase struct {
	customerRepository repository.CustomerRepository
}

func NewDeleteCustomerUseCase(customerRepository repository.CustomerRepository) *DeleteCustomerUseCase {
	return &DeleteCustomerUseCase{
		customerRepository: customerRepository,
	}
}

func (uc *DeleteCustomerUseCase) Execute(id string) error {
	return uc.customerRepository.Delete(id)
}
