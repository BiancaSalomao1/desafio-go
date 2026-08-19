package product

import (
	"testing"

	"product-service/internal/domain"
)

type releaseStockRepositoryStub struct {
	products map[string]*domain.Product
	updates  int
}

func (r *releaseStockRepositoryStub) Save(product *domain.Product) error {
	r.products[product.ID] = product
	return nil
}

func (r *releaseStockRepositoryStub) Update(product *domain.Product) error {
	r.products[product.ID] = product
	r.updates++
	return nil
}

func (r *releaseStockRepositoryStub) Delete(id string) error {
	delete(r.products, id)
	return nil
}

func (r *releaseStockRepositoryStub) FindByID(id string) (*domain.Product, error) {
	product, ok := r.products[id]
	if !ok {
		return nil, domain.ErrProductNotFound
	}

	return product, nil
}

func (r *releaseStockRepositoryStub) FindAll() ([]*domain.Product, error) {
	result := make([]*domain.Product, 0, len(r.products))

	for _, product := range r.products {
		result = append(result, product)
	}

	return result, nil
}

func TestReleaseStockUseCase_Execute(t *testing.T) {
	product := domain.NewProduct(
		"product-1",
		"Notebook",
		3500,
		8,
	)

	repository := &releaseStockRepositoryStub{
		products: map[string]*domain.Product{
			product.ID: product,
		},
	}

	useCase := NewReleaseStockUseCase(repository)

	err := useCase.Execute([]StockItem{
		{
			ProductID: "product-1",
			Quantity:  2,
		},
	})

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if product.Stock != 10 {
		t.Fatalf(
			"expected stock 10, got %d",
			product.Stock,
		)
	}

	if repository.updates != 1 {
		t.Fatalf(
			"expected 1 update, got %d",
			repository.updates,
		)
	}
}
