package customer

import (
	"errors"
	"testing"

	"desafio-go/internal/domain"
)

func TestUpdateCustomerUseCase_Execute(t *testing.T) {

	t.Run("should update customer successfully", func(t *testing.T) {

		repository := &customerRepositoryStub{
			customer: &domain.Customer{
				ID:    "1",
				Name:  "João",
				Email: "joao@email.com",
			},
		}

		useCase := NewUpdateCustomerUseCase(repository)

		err := useCase.Execute(&domain.Customer{
			ID:    "1",
			Name:  "Maria",
			Email: "maria@email.com",
		})

		if err != nil {
			t.Fatalf("expected nil, got %v", err)
		}

		if repository.customer.Name != "Maria" {
			t.Fatal("expected updated name")
		}

		if repository.customer.Email != "maria@email.com" {
			t.Fatal("expected updated email")
		}
	})

	t.Run("should return repository find error", func(t *testing.T) {

		repository := &customerRepositoryStub{
			findError: errors.New("customer not found"),
		}

		useCase := NewUpdateCustomerUseCase(repository)

		err := useCase.Execute(&domain.Customer{ID: "1"})

		if err == nil {
			t.Fatal("expected repository error")
		}
	})

	t.Run("should return validation error", func(t *testing.T) {

		repository := &customerRepositoryStub{
			customer: &domain.Customer{
				ID:    "1",
				Name:  "João",
				Email: "joao@email.com",
			},
		}

		useCase := NewUpdateCustomerUseCase(repository)

		err := useCase.Execute(&domain.Customer{
			ID:    "1",
			Name:  "",
			Email: "",
		})

		if err == nil {
			t.Fatal("expected validation error")
		}
	})

	t.Run("should return update error", func(t *testing.T) {

		repository := &customerRepositoryStub{
			customer: &domain.Customer{
				ID:    "1",
				Name:  "João",
				Email: "joao@email.com",
			},
			updateError: errors.New("database error"),
		}

		useCase := NewUpdateCustomerUseCase(repository)

		err := useCase.Execute(&domain.Customer{
			ID:    "1",
			Name:  "Maria",
			Email: "maria@email.com",
		})

		if err == nil {
			t.Fatal("expected update error")
		}
	})

	t.Run("should call update once", func(t *testing.T) {

		repository := &customerRepositorySpy{
			customer: &domain.Customer{
				ID:    "1",
				Name:  "João",
				Email: "joao@email.com",
			},
		}

		useCase := NewUpdateCustomerUseCase(repository)

		_ = useCase.Execute(&domain.Customer{
			ID:    "1",
			Name:  "Maria",
			Email: "maria@email.com",
		})

		if !repository.updateCalled {
			t.Fatal("expected Update to be called")
		}

		if repository.updateCalls != 1 {
			t.Fatalf("expected 1 call, got %d", repository.updateCalls)
		}
	})

	t.Run("should satisfy mock expectations", func(t *testing.T) {

		repository := &customerRepositoryMock{
			customer: &domain.Customer{
				ID:    "1",
				Name:  "João",
				Email: "joao@email.com",
			},
			expectedUpdateCalls: 1,
		}

		useCase := NewUpdateCustomerUseCase(repository)

		err := useCase.Execute(&domain.Customer{
			ID:    "1",
			Name:  "Maria",
			Email: "maria@email.com",
		})

		if err != nil {
			t.Fatalf("expected nil, got %v", err)
		}

		if err := repository.Verify(); err != nil {
			t.Fatal(err)
		}
	})
}
