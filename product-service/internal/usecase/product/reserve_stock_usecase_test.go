package product

import (
	"errors"
	"testing"

	"product-service/internal/domain"
)

type reserveStockRepositoryStub struct {
	products map[string]*domain.Product
	updates  int
}

func (r *reserveStockRepositoryStub) Save(product *domain.Product) error {
	r.products[product.ID] = product
	return nil
}

func (r *reserveStockRepositoryStub) Update(product *domain.Product) error {
	r.products[product.ID] = product
	r.updates++
	return nil
}

func (r *reserveStockRepositoryStub) Delete(id string) error {
	delete(r.products, id)
	return nil
}

func (r *reserveStockRepositoryStub) FindByID(id string) (*domain.Product, error) {
	product, ok := r.products[id]
	if !ok {
		return nil, domain.ErrProductNotFound
	}

	return product, nil
}

func (r *reserveStockRepositoryStub) FindAll() ([]*domain.Product, error) {
	result := make([]*domain.Product, 0, len(r.products))

	for _, product := range r.products {
		result = append(result, product)
	}

	return result, nil
}

func TestReserveStockUseCase_Execute(t *testing.T) {
	product := domain.NewProduct(
		"product-1",
		"Notebook",
		3500,
		10,
	)

	repository := &reserveStockRepositoryStub{
		products: map[string]*domain.Product{
			product.ID: product,
		},
	}

	useCase := NewReserveStockUseCase(repository)

	err := useCase.Execute([]StockItem{
		{
			ProductID: "product-1",
			Quantity:  2,
		},
	})

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if product.Stock != 8 {
		t.Fatalf("expected stock 8, got %d", product.Stock)
	}

	if repository.updates != 1 {
		t.Fatalf("expected 1 update, got %d", repository.updates)
	}
}

func TestReserveStockUseCase_InsufficientStock(t *testing.T) {
	product := domain.NewProduct(
		"product-1",
		"Notebook",
		3500,
		2,
	)

	repository := &reserveStockRepositoryStub{
		products: map[string]*domain.Product{
			product.ID: product,
		},
	}

	useCase := NewReserveStockUseCase(repository)

	err := useCase.Execute([]StockItem{
		{
			ProductID: "product-1",
			Quantity:  3,
		},
	})

	if !errors.Is(err, domain.ErrInsufficientStock) {
		t.Fatalf(
			"expected ErrInsufficientStock, got %v",
			err,
		)
	}

	if product.Stock != 2 {
		t.Fatalf(
			"expected stock to remain 2, got %d",
			product.Stock,
		)
	}

	if repository.updates != 0 {
		t.Fatalf(
			"expected no update, got %d",
			repository.updates,
		)
	}
}

func TestReserveStockUseCase_InvalidQuantity(t *testing.T) {
	product := domain.NewProduct(
		"product-1",
		"Notebook",
		3500,
		10,
	)

	repository := &reserveStockRepositoryStub{
		products: map[string]*domain.Product{
			product.ID: product,
		},
	}

	useCase := NewReserveStockUseCase(repository)

	err := useCase.Execute([]StockItem{
		{
			ProductID: "product-1",
			Quantity:  0,
		},
	})

	if !errors.Is(err, domain.ErrInvalidQuantity) {
		t.Fatalf(
			"expected ErrInvalidQuantity, got %v",
			err,
		)
	}

	if product.Stock != 10 {
		t.Fatalf(
			"expected stock to remain 10, got %d",
			product.Stock,
		)
	}
}
