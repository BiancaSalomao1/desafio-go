package customer

/*
Testes do CreateCustomerUseCase.

Responsabilidades:
- validar criação de cliente;
- validar regras de negócio;
- validar interação com o repositório.

Cenários:
- criação com sucesso;
- cliente inválido;
- erro do repositório;
- chamada do Save();
- validação do Mock.
*/

import (
	"errors"
	"testing"

	"desafio-go/internal/domain"
	"desafio-go/internal/security"
)

func validCustomer(t *testing.T) *domain.Customer {
	t.Helper()

	hash, err := security.HashPassword("123456")
	if err != nil {
		t.Fatalf("failed to hash password: %v", err)
	}

	return domain.NewCustomer(
		"1",
		"João da Silva",
		"joao@email.com",
		hash,
	)
}

func TestCreateCustomerUseCase_Execute(t *testing.T) {

	t.Run("should create customer successfully", func(t *testing.T) {

		repository := &customerRepositoryStub{}

		useCase := NewCreateCustomerUseCase(repository)

		customer := validCustomer(t)

		err := useCase.Execute(customer)

		if err != nil {
			t.Fatalf("expected nil, got %v", err)
		}

	})

	t.Run("should return validation error", func(t *testing.T) {

		repository := &customerRepositoryStub{}

		useCase := NewCreateCustomerUseCase(repository)

		customer := &domain.Customer{}

		err := useCase.Execute(customer)

		if err == nil {
			t.Fatal("expected validation error")
		}

	})

	t.Run("should return repository error", func(t *testing.T) {

		repository := &customerRepositoryStub{
			saveError: errors.New("database error"),
		}

		useCase := NewCreateCustomerUseCase(repository)

		customer := validCustomer(t)

		err := useCase.Execute(customer)

		if err == nil {
			t.Fatal("expected repository error")
		}

		if err.Error() != "database error" {
			t.Fatalf(
				"expected database error, got %v",
				err,
			)
		}
	})

	t.Run("should call save once", func(t *testing.T) {

		repository := &customerRepositorySpy{}

		useCase := NewCreateCustomerUseCase(repository)

		customer := validCustomer(t)

		err := useCase.Execute(customer)

		if err != nil {
			t.Fatalf("expected nil, got %v", err)
		}

		if !repository.saveCalled {
			t.Fatal("expected Save to be called")
		}

		if repository.saveCalls != 1 {
			t.Fatalf(
				"expected 1 call, got %d",
				repository.saveCalls,
			)
		}
	})

	t.Run("should satisfy mock expectations", func(t *testing.T) {

		repository := &customerRepositoryMock{
			expectedSaveCalls: 1,
		}

		useCase := NewCreateCustomerUseCase(repository)

		customer := validCustomer(t)

		err := useCase.Execute(customer)

		if err != nil {
			t.Fatalf("expected nil, got %v", err)
		}

		if err := repository.Verify(); err != nil {
			t.Fatal(err)
		}
	})
}
