package customer

/*
struct CreateCustomerUseCase

Responsabilidades:
- criar cliente;
- validar dados;
- gerar hash da senha;
- salvar cliente;
- registrar eventos relevantes da operação.

Campos:
- customerRepository

Métodos:
- NewCreateCustomerUseCase()
- Execute()
*/

import (
	"log/slog"

	"orders-api/internal/domain"
	"orders-api/internal/repository"
	"orders-api/internal/security"
)

const customerServiceName = "orders-api"

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

	slog.Info(
		"customer creation started",
		"service", customerServiceName,
		"operation", "create_customer",
		"customer_id", customer.ID,
	)

	// Quando a senha ainda está em texto puro,
	// gera o hash antes da validação e persistência.
	if customer.Password != "" {

		hash, err := security.HashPassword(customer.Password)
		if err != nil {

			slog.Error(
				"customer creation failed",
				"service", customerServiceName,
				"operation", "create_customer",
				"result", "error",
				"customer_id", customer.ID,
				"error", err,
			)

			return err
		}

		customer.PasswordHash = hash
		customer.Password = ""
	}

	if err := customer.Validate(); err != nil {

		slog.Error(
			"customer creation failed",
			"service", customerServiceName,
			"operation", "create_customer",
			"result", "error",
			"customer_id", customer.ID,
			"error", err,
		)

		return err
	}

	if err := uc.customerRepository.Save(customer); err != nil {

		slog.Error(
			"customer creation failed",
			"service", customerServiceName,
			"operation", "create_customer",
			"result", "error",
			"customer_id", customer.ID,
			"error", err,
		)

		return err
	}

	slog.Info(
		"customer created",
		"service", customerServiceName,
		"operation", "create_customer",
		"result", "success",
		"customer_id", customer.ID,
	)

	return nil
}
