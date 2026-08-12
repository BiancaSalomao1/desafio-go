package product

/*
struct CreateProductUseCase
- criar produto;
- validar dados;
- salvar produto.

Métodos:
- NewCreateProductUseCase()
- Execute()
*/

import (
	"desafio-go/orders-api/internal/domain"
	"desafio-go/orders-api/internal/repository"
)

type CreateProductUseCase struct {
	productRepository repository.ProductRepository
}

func NewCreateProductUseCase(productRepository repository.ProductRepository) *CreateProductUseCase {
	return &CreateProductUseCase{
		productRepository: productRepository,
	}
}

func (uc *CreateProductUseCase) Execute(product *domain.Product) error {

	if err := product.Validate(); err != nil {
		return err
	}

	return uc.productRepository.Save(product)
}
