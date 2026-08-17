package customer

import (
	"errors"
	"testing"

	"orders-api/internal/domain"
)

func TestUpdateCustomerUseCase_Execute(t *testing.T) {

	t.Run("should update customer successfully", func(t *testing.T) {

		repository := &customerRepositoryStub{
			customer: domain.NewCustomer(
				"1",
				"João",
				"joao@email.com",
				"hash",
			),
		}

		useCase := NewUpdateCustomerUseCase(repository)

		err := useCase.Execute(
			domain.NewCustomer(
				"1",
				"Maria",
				"maria@email.com",
				"",
			),
		)

		if err != nil {
			t.Fatalf("expected nil, got %v", err)
		}

		if repository.customer.Name != "Maria" {
			t.Fatal("expected updated name")
		}

		if repository.customer.Email != "maria@email.com" {
			t.Fatal("expected updated email")
		}

		// O hash deve ser preservado
		if repository.customer.PasswordHash != "hash" {
			t.Fatal("password hash should be preserved")
		}
	})

	t.Run("should return repository find error", func(t *testing.T) {

		repository := &customerRepositoryStub{
			findError: errors.New("customer not found"),
		}

		useCase := NewUpdateCustomerUseCase(repository)

		err := useCase.Execute(
			domain.NewCustomer(
				"1",
				"Maria",
				"maria@email.com",
				"",
			),
		)

		if err == nil {
			t.Fatal("expected repository error")
		}
	})

	t.Run("should return validation error", func(t *testing.T) {

		repository := &customerRepositoryStub{
			customer: domain.NewCustomer(
				"1",
				"João",
				"joao@email.com",
				"hash",
			),
		}

		useCase := NewUpdateCustomerUseCase(repository)

		err := useCase.Execute(
			domain.NewCustomer(
				"1",
				"",
				"",
				"",
			),
		)

		if !errors.Is(err, domain.ErrCustomerInvalid) {
			t.Fatalf(
				"expected %v, got %v",
				domain.ErrCustomerInvalid,
				err,
			)
		}
	})

	t.Run("should return update error", func(t *testing.T) {

		repository := &customerRepositoryStub{
			customer: domain.NewCustomer(
				"1",
				"João",
				"joao@email.com",
				"hash",
			),
			updateError: errors.New("database error"),
		}

		useCase := NewUpdateCustomerUseCase(repository)

		err := useCase.Execute(
			domain.NewCustomer(
				"1",
				"Maria",
				"maria@email.com",
				"",
			),
		)

		if err == nil {
			t.Fatal("expected update error")
		}

		if err.Error() != "database error" {
			t.Fatalf("expected database error, got %v", err)
		}
	})

	t.Run("should call update once", func(t *testing.T) {

		repository := &customerRepositorySpy{
			customer: domain.NewCustomer(
				"1",
				"João",
				"joao@email.com",
				"hash",
			),
		}

		useCase := NewUpdateCustomerUseCase(repository)

		_ = useCase.Execute(
			domain.NewCustomer(
				"1",
				"Maria",
				"maria@email.com",
				"",
			),
		)

		if !repository.updateCalled {
			t.Fatal("expected Update to be called")
		}

		if repository.updateCalls != 1 {
			t.Fatalf(
				"expected 1 call, got %d",
				repository.updateCalls,
			)
		}
	})

	t.Run("should satisfy mock expectations", func(t *testing.T) {

		repository := &customerRepositoryMock{
			customer: domain.NewCustomer(
				"1",
				"João",
				"joao@email.com",
				"hash",
			),
			expectedUpdateCalls: 1,
		}

		useCase := NewUpdateCustomerUseCase(repository)

		err := useCase.Execute(
			domain.NewCustomer(
				"1",
				"Maria",
				"maria@email.com",
				"",
			),
		)

		if err != nil {
			t.Fatalf("expected nil, got %v", err)
		}

		if err := repository.Verify(); err != nil {
			t.Fatal(err)
		}
	})
}
