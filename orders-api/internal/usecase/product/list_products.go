package product

/*
struct ListProductsUseCase
- listar produtos.

Métodos:
- NewListProductsUseCase()
- Execute()
*/

import (
	"orders-api/internal/domain"
	"orders-api/internal/repository"
)

type ListProductsUseCase struct {
	productRepository repository.ProductRepository
}

func NewListProductsUseCase(productRepository repository.ProductRepository) *ListProductsUseCase {
	return &ListProductsUseCase{
		productRepository: productRepository,
	}
}

func (uc *ListProductsUseCase) Execute() ([]*domain.Product, error) {
	return uc.productRepository.FindAll()
}
