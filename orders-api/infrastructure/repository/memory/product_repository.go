package memory

/*
struct MemoryProductRepository

Responsabilidades:

- armazenar produtos em memória;
- implementar ProductRepository.

Métodos:
- NewMemoryProductRepository()
- Save()
- Update()
- Delete()
- FindByID()
- FindAll()
*/

import (
	"orders-api/internal/domain"

	"orders-api/internal/repository"
)

type MemoryProductRepository struct {
	products map[string]*domain.Product
}

func NewMemoryProductRepository() repository.ProductRepository {
	return &MemoryProductRepository{
		products: make(map[string]*domain.Product),
	}
}

func (r *MemoryProductRepository) Save(product *domain.Product) error {

	if err := product.Validate(); err != nil {
		return err
	}

	r.products[product.ID] = product

	return nil
}

func (r *MemoryProductRepository) Update(product *domain.Product) error {

	if _, exists := r.products[product.ID]; !exists {
		return domain.ErrProductNotFound
	}

	if err := product.Validate(); err != nil {
		return err
	}

	r.products[product.ID] = product

	return nil
}

func (r *MemoryProductRepository) Delete(id string) error {

	if _, exists := r.products[id]; !exists {
		return domain.ErrProductNotFound
	}

	delete(r.products, id)

	return nil
}

func (r *MemoryProductRepository) FindByID(id string) (*domain.Product, error) {

	product, exists := r.products[id]

	if !exists {
		return nil, domain.ErrProductNotFound
	}

	return product, nil
}

func (r *MemoryProductRepository) FindAll() ([]*domain.Product, error) {

	products := make([]*domain.Product, 0, len(r.products))

	for _, product := range r.products {
		products = append(products, product)
	}

	return products, nil
}
