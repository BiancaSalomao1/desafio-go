package product

/*
struct DeleteProductUseCase

Responsabilidades:
- remover um produto.

Campos:
- productRepository

Métodos:
- NewDeleteProductUseCase()
- Execute()
*/

import (
	"strings"

	"orders-api/internal/domain"
	"orders-api/internal/repository"
)

type DeleteProductUseCase struct {
	productRepository repository.ProductRepository
}

func NewDeleteProductUseCase(
	productRepository repository.ProductRepository,
) *DeleteProductUseCase {

	return &DeleteProductUseCase{
		productRepository: productRepository,
	}
}

func (uc *DeleteProductUseCase) Execute(
	id string,
) error {

	if _, err := uc.productRepository.FindByID(id); err != nil {
		return err
	}

	err := uc.productRepository.Delete(id)
	if err != nil {

		if strings.Contains(
			err.Error(),
			"SQLSTATE 23503",
		) {
			return domain.ErrProductInUse
		}

		return err
	}

	return nil
}
