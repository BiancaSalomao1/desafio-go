package customer

/*
struct DeleteCustomerUseCase

Responsabilidades:
- remover um cliente.

Campos:
- customerRepository

Métodos:
- NewDeleteCustomerUseCase()
- Execute()
*/

import (
	"strings"

	"orders-api/internal/domain"
	"orders-api/internal/repository"
)

type DeleteCustomerUseCase struct {
	customerRepository repository.CustomerRepository
}

func NewDeleteCustomerUseCase(
	customerRepository repository.CustomerRepository,
) *DeleteCustomerUseCase {

	return &DeleteCustomerUseCase{
		customerRepository: customerRepository,
	}
}

func (uc *DeleteCustomerUseCase) Execute(
	id string,
) error {

	if _, err := uc.customerRepository.FindByID(id); err != nil {
		return err
	}

	err := uc.customerRepository.Delete(id)
	if err != nil {

		if strings.Contains(
			err.Error(),
			"SQLSTATE 23503",
		) {
			return domain.ErrCustomerInUse
		}

		return err
	}

	return nil
}
