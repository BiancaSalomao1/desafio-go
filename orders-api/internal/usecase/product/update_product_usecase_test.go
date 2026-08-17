package product

/*
Testes do UpdateProductUseCase.

Responsabilidades:
- validar atualização de produto;
- validar regras de negócio;
- validar interação com o repositório.

Cenários:
- atualização com sucesso;
- produto inexistente;
- produto inválido;
- erro do repositório;
- chamada do Update();
- validação do Mock.
*/

import (
	"errors"
	"testing"

	"orders-api/internal/domain"
)

func TestUpdateProductUseCase_Execute(t *testing.T) {

	t.Run("should update product successfully", func(t *testing.T) {

		repository := &productRepositoryStub{
			product: &domain.Product{
				ID:    "1",
				Name:  "Mouse",
				Price: 100,
				Stock: 10,
			},
		}

		useCase := NewUpdateProductUseCase(repository)

		product := &domain.Product{
			ID:    "1",
			Name:  "Mouse Gamer",
			Price: 150,
			Stock: 20,
		}

		err := useCase.Execute(product)

		if err != nil {
			t.Fatalf("expected nil, got %v", err)
		}

		if repository.product.Name != "Mouse Gamer" {
			t.Fatal("expected product name updated")
		}

		if repository.product.Price != 150 {
			t.Fatal("expected product price updated")
		}

		if repository.product.Stock != 20 {
			t.Fatal("expected product stock updated")
		}
	})

	t.Run("should return repository find error", func(t *testing.T) {

		repository := &productRepositoryStub{
			findError: errors.New("product not found"),
		}

		useCase := NewUpdateProductUseCase(repository)

		err := useCase.Execute(&domain.Product{
			ID: "1",
		})

		if err == nil {
			t.Fatal("expected repository error")
		}

		if err.Error() != "product not found" {
			t.Fatalf("expected product not found, got %v", err)
		}
	})

	t.Run("should return validation error", func(t *testing.T) {

		repository := &productRepositoryStub{
			product: &domain.Product{
				ID:    "1",
				Name:  "Mouse",
				Price: 100,
				Stock: 10,
			},
		}

		useCase := NewUpdateProductUseCase(repository)

		err := useCase.Execute(&domain.Product{
			ID:    "1",
			Name:  "",
			Price: 150,
			Stock: 10,
		})

		if err == nil {
			t.Fatal("expected validation error")
		}
	})

	t.Run("should return update error", func(t *testing.T) {

		repository := &productRepositoryStub{
			product: &domain.Product{
				ID:    "1",
				Name:  "Mouse",
				Price: 100,
				Stock: 10,
			},
			updateError: errors.New("database error"),
		}

		useCase := NewUpdateProductUseCase(repository)

		err := useCase.Execute(&domain.Product{
			ID:    "1",
			Name:  "Mouse Gamer",
			Price: 150,
			Stock: 20,
		})

		if err == nil {
			t.Fatal("expected update error")
		}

		if err.Error() != "database error" {
			t.Fatalf("expected database error, got %v", err)
		}
	})

	t.Run("should call update once", func(t *testing.T) {

		repository := &productRepositorySpy{
			product: &domain.Product{
				ID:    "1",
				Name:  "Mouse",
				Price: 100,
				Stock: 10,
			},
		}

		useCase := NewUpdateProductUseCase(repository)

		err := useCase.Execute(&domain.Product{
			ID:    "1",
			Name:  "Mouse Gamer",
			Price: 150,
			Stock: 20,
		})

		if err != nil {
			t.Fatalf("expected nil, got %v", err)
		}

		if !repository.updateCalled {
			t.Fatal("expected Update to be called")
		}

		if repository.updateCalls != 1 {
			t.Fatalf("expected 1 call, got %d", repository.updateCalls)
		}
	})

	t.Run("should satisfy mock expectations", func(t *testing.T) {

		repository := &productRepositoryMock{
			product: &domain.Product{
				ID:    "1",
				Name:  "Mouse",
				Price: 100,
				Stock: 10,
			},
			expectedUpdateCalls: 1,
		}

		useCase := NewUpdateProductUseCase(repository)

		err := useCase.Execute(&domain.Product{
			ID:    "1",
			Name:  "Mouse Gamer",
			Price: 150,
			Stock: 20,
		})

		if err != nil {
			t.Fatalf("expected nil, got %v", err)
		}

		if err := repository.Verify(); err != nil {
			t.Fatal(err)
		}
	})
}
