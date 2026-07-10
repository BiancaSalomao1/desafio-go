package handler

import "desafio-go/internal/domain"

/*
interface CreateProductUseCase

Responsabilidades:
- criar um produto.

Métodos:
- Execute()
*/

type CreateProductUseCase interface {
	Execute(product *domain.Product) error
}

/*
interface GetProductUseCase

Responsabilidades:
- buscar um produto.

Métodos:
- Execute()
*/

type GetProductUseCase interface {
	Execute(id string) (*domain.Product, error)
}

/*
interface ListProductsUseCase

Responsabilidades:
- listar produtos.

Métodos:
- Execute()
*/

type ListProductsUseCase interface {
	Execute() ([]*domain.Product, error)
}
