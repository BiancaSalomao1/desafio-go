package repository

import "product-service/internal/domain"

type ProductRepository interface {
	Save(product *domain.Product) error
	Update(product *domain.Product) error
	Delete(id string) error
	FindByID(id string) (*domain.Product, error)
	FindAll() ([]*domain.Product, error)
}
