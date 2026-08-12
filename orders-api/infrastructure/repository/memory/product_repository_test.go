package memory

import (
	"testing"

	"desafio-go/orders-api/internal/domain"
)

func newProductRepository() *MemoryProductRepository {
	return NewMemoryProductRepository().(*MemoryProductRepository)
}

func newValidProduct() *domain.Product {
	return domain.NewProduct(
		"PROD001",
		"Notebook",
		5000,
		10,
	)
}

func TestNewMemoryProductRepository(t *testing.T) {
	repo := NewMemoryProductRepository()

	if repo == nil {
		t.Fatal("expected repository")
	}
}

func TestSave(t *testing.T) {
	repo := newProductRepository()

	product := newValidProduct()

	err := repo.Save(product)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(repo.products) != 1 {
		t.Fatalf("expected 1 product, got %d", len(repo.products))
	}
}

func TestSave_InvalidProduct(t *testing.T) {
	repo := newProductRepository()

	product := domain.NewProduct(
		"",
		"",
		0,
		0,
	)

	err := repo.Save(product)

	if err == nil {
		t.Fatal("expected validation error")
	}

	if len(repo.products) != 0 {
		t.Fatal("invalid product should not be saved")
	}
}

func TestUpdate(t *testing.T) {
	repo := newProductRepository()

	product := newValidProduct()

	_ = repo.Save(product)

	product.Name = "Notebook Gamer"

	err := repo.Update(product)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	saved, _ := repo.FindByID(product.ID)

	if saved.Name != "Notebook Gamer" {
		t.Fatal("product was not updated")
	}
}

func TestUpdate_ProductNotFound(t *testing.T) {
	repo := newProductRepository()

	product := newValidProduct()

	err := repo.Update(product)

	if err != domain.ErrProductNotFound {
		t.Fatalf("expected %v, got %v",
			domain.ErrProductNotFound,
			err,
		)
	}
}

func TestUpdate_InvalidProduct(t *testing.T) {
	repo := newProductRepository()

	product := newValidProduct()

	_ = repo.Save(product)

	product.Name = ""

	err := repo.Update(product)

	if err == nil {
		t.Fatal("expected validation error")
	}
}

func TestDelete(t *testing.T) {
	repo := newProductRepository()

	product := newValidProduct()

	_ = repo.Save(product)

	err := repo.Delete(product.ID)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(repo.products) != 0 {
		t.Fatal("product should have been removed")
	}
}

func TestDelete_ProductNotFound(t *testing.T) {
	repo := newProductRepository()

	err := repo.Delete("INVALID")

	if err != domain.ErrProductNotFound {
		t.Fatalf("expected %v, got %v",
			domain.ErrProductNotFound,
			err,
		)
	}
}

func TestFindByID(t *testing.T) {
	repo := newProductRepository()

	product := newValidProduct()

	_ = repo.Save(product)

	found, err := repo.FindByID(product.ID)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if found.ID != product.ID {
		t.Fatal("wrong product returned")
	}
}

func TestFindByID_NotFound(t *testing.T) {
	repo := newProductRepository()

	_, err := repo.FindByID("INVALID")

	if err != domain.ErrProductNotFound {
		t.Fatalf("expected %v, got %v",
			domain.ErrProductNotFound,
			err,
		)
	}
}

func TestFindAll(t *testing.T) {
	repo := newProductRepository()

	_ = repo.Save(domain.NewProduct(
		"1",
		"Notebook",
		100,
		1,
	))

	_ = repo.Save(domain.NewProduct(
		"2",
		"Mouse",
		50,
		2,
	))

	products, err := repo.FindAll()

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(products) != 2 {
		t.Fatalf("expected 2 products, got %d", len(products))
	}
}

func TestFindAll_Empty(t *testing.T) {
	repo := newProductRepository()

	products, err := repo.FindAll()

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(products) != 0 {
		t.Fatalf("expected empty slice, got %d", len(products))
	}
}
