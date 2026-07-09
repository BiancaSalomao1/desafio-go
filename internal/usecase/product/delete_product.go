package product

/*
struct DeleteProductUseCase
- excluir produto.

Métodos:
- NewDeleteProductUseCase()
- Execute()
*/

import "desafio-go/internal/repository"

type DeleteProductUseCase struct {
	productRepository repository.ProductRepository
}

func NewDeleteProductUseCase(productRepository repository.ProductRepository) *DeleteProductUseCase {
	return &DeleteProductUseCase{
		productRepository: productRepository,
	}
}

func (uc *DeleteProductUseCase) Execute(id string) error {
	return uc.productRepository.Delete(id)
}
