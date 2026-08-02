package customer

import (
	"desafio-go/internal/domain"
	"desafio-go/internal/repository"
	"desafio-go/internal/security"
)

type CreateCustomerUseCase struct {
	customerRepository repository.CustomerRepository
}

func NewCreateCustomerUseCase(
	customerRepository repository.CustomerRepository,
) *CreateCustomerUseCase {

	return &CreateCustomerUseCase{
		customerRepository: customerRepository,
	}
}

func (uc *CreateCustomerUseCase) Execute(
	customer *domain.Customer,
) error {

	hash, err := security.HashPassword(customer.Password)
	if err != nil {
		return err
	}

	customer.PasswordHash = hash
	customer.Password = ""

	if err := customer.Validate(); err != nil {
		return err
	}

	return uc.customerRepository.Save(customer)
}
