package product

/*
struct GetProductUseCase
- buscar produto por ID.

Métodos:
- NewGetProductUseCase()
- Execute()
*/

import (
	"desafio-go/internal/domain"
	"desafio-go/internal/repository"
)

type GetProductUseCase struct {
	productRepository repository.ProductRepository
}

func NewGetProductUseCase(productRepository repository.ProductRepository) *GetProductUseCase {
	return &GetProductUseCase{
		productRepository: productRepository,
	}
}

func (uc *GetProductUseCase) Execute(id string) (*domain.Product, error) {
	return uc.productRepository.FindByID(id)
}
