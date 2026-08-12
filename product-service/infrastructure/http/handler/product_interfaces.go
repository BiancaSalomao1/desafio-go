package handler

import "product-service/internal/domain"

type CreateProductUseCase interface {
	Execute(product *domain.Product) error
}

type GetProductUseCase interface {
	Execute(id string) (*domain.Product, error)
}

type ListProductsUseCase interface {
	Execute() ([]*domain.Product, error)
}

type UpdateProductUseCase interface {
	Execute(product *domain.Product) error
}

type DeleteProductUseCase interface {
	Execute(id string) error
}
