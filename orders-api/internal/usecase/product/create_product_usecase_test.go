package product

/*
Testes do CreateProductUseCase.

Responsabilidades:
- validar criação de produto;
- validar regras de negócio;
- validar interação com o repositório.

Cenários:
- criação com sucesso;
- produto inválido;
- erro do repositório;
- chamada do Save();
- validação do Mock.
*/

import (
	"errors"
	"testing"

	"desafio-go/orders-api/internal/domain"
)

func TestCreateProductUseCase_Execute(t *testing.T) {

	t.Run("should create product successfully", func(t *testing.T) {

		repository := &productRepositoryStub{}

		useCase := NewCreateProductUseCase(repository)

		product := &domain.Product{
			Name:  "Notebook",
			Price: 3500,
			Stock: 10,
		}

		err := useCase.Execute(product)

		if err != nil {
			t.Fatalf("expected nil, got %v", err)
		}
	})

	t.Run("should return validation error", func(t *testing.T) {

		repository := &productRepositoryStub{}

		useCase := NewCreateProductUseCase(repository)

		product := &domain.Product{}

		err := useCase.Execute(product)

		if err == nil {
			t.Fatal("expected validation error")
		}
	})

	t.Run("should return repository error", func(t *testing.T) {

		repository := &productRepositoryStub{
			saveError: errors.New("database error"),
		}

		useCase := NewCreateProductUseCase(repository)

		product := &domain.Product{
			Name:  "Notebook",
			Price: 3500,
			Stock: 10,
		}

		err := useCase.Execute(product)

		if err == nil {
			t.Fatal("expected repository error")
		}

		if err.Error() != "database error" {
			t.Fatalf("expected database error, got %v", err)
		}
	})

	t.Run("should call save once", func(t *testing.T) {

		repository := &productRepositorySpy{}

		useCase := NewCreateProductUseCase(repository)

		product := &domain.Product{
			Name:  "Notebook",
			Price: 3500,
			Stock: 10,
		}

		err := useCase.Execute(product)

		if err != nil {
			t.Fatalf("expected nil, got %v", err)
		}

		if !repository.saveCalled {
			t.Fatal("expected Save to be called")
		}

		if repository.saveCalls != 1 {
			t.Fatalf("expected 1 call, got %d", repository.saveCalls)
		}

		if repository.product != product {
			t.Fatal("expected same product instance")
		}
	})

	t.Run("should satisfy mock expectations", func(t *testing.T) {

		repository := &productRepositoryMock{
			expectedSaveCalls: 1,
		}

		useCase := NewCreateProductUseCase(repository)

		product := &domain.Product{
			Name:  "Notebook",
			Price: 3500,
			Stock: 10,
		}

		err := useCase.Execute(product)

		if err != nil {
			t.Fatalf("expected nil, got %v", err)
		}

		err = repository.Verify()

		if err != nil {
			t.Fatal(err)
		}
	})
}
