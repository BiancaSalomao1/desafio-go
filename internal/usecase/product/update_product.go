package product

/*
struct UpdateProductUseCase
- atualizar produto.

Métodos:
- NewUpdateProductUseCase()
- Execute()
*/

import (
	"desafio-go/internal/domain"
	"desafio-go/internal/repository"
)

type UpdateProductUseCase struct {
	productRepository repository.ProductRepository
}

func NewUpdateProductUseCase(productRepository repository.ProductRepository) *UpdateProductUseCase {
	return &UpdateProductUseCase{
		productRepository: productRepository,
	}
}

func (uc *UpdateProductUseCase) Execute(product *domain.Product) error {

	if err := product.Validate(); err != nil {
		return err
	}

	return uc.productRepository.Update(product)
}
