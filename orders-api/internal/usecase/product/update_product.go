package product

/*
struct UpdateProductUseCase

Responsabilidades:
- atualizar um produto.

Campos:
- productRepository

Métodos:
- NewUpdateProductUseCase()
- Execute()
*/

import (
	"orders-api/internal/domain"
	"orders-api/internal/repository"
)

type UpdateProductUseCase struct {
	productRepository repository.ProductRepository
}

func NewUpdateProductUseCase(
	productRepository repository.ProductRepository,
) *UpdateProductUseCase {

	return &UpdateProductUseCase{
		productRepository: productRepository,
	}
}

func (uc *UpdateProductUseCase) Execute(
	product *domain.Product,
) error {

	currentProduct, err := uc.productRepository.FindByID(product.ID)
	if err != nil {
		return err
	}

	if err := currentProduct.Update(
		product.Name,
		product.Price,
		product.Stock,
	); err != nil {
		return err
	}

	return uc.productRepository.Update(currentProduct)
}
