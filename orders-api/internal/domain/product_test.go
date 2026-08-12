package domain

import (
	"errors"
	"testing"
)

func TestNewProduct(t *testing.T) {

	product := NewProduct(
		"1",
		"Notebook",
		3500,
		10,
	)

	if product == nil {
		t.Fatal("expected product")
	}

	if product.ID != "1" {
		t.Fatalf("expected id 1, got %s", product.ID)
	}

	if product.Name != "Notebook" {
		t.Fatalf("expected Notebook, got %s", product.Name)
	}

	if product.Price != 3500 {
		t.Fatalf("expected 3500, got %.2f", product.Price)
	}

	if product.Stock != 10 {
		t.Fatalf("expected stock 10, got %d", product.Stock)
	}
}

func TestProduct_Validate(t *testing.T) {

	tests := []struct {
		name        string
		product     Product
		expectedErr error
	}{
		{
			name: "valid product",
			product: Product{
				Name:  "Notebook",
				Price: 3500,
				Stock: 10,
			},
			expectedErr: nil,
		},
		{
			name: "empty name",
			product: Product{
				Name:  "",
				Price: 3500,
				Stock: 10,
			},
			expectedErr: ErrProductNameRequired,
		},
		{
			name: "invalid price",
			product: Product{
				Name:  "Notebook",
				Price: 0,
				Stock: 10,
			},
			expectedErr: ErrInvalidPrice,
		},
		{
			name: "negative stock",
			product: Product{
				Name:  "Notebook",
				Price: 3500,
				Stock: -1,
			},
			expectedErr: ErrInvalidStock,
		},
	}

	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {

			err := tt.product.Validate()

			if !errors.Is(err, tt.expectedErr) {
				t.Fatalf("expected %v, got %v", tt.expectedErr, err)
			}
		})
	}
}

func TestProduct_ReduceStock(t *testing.T) {

	t.Run("should reduce stock", func(t *testing.T) {

		product := Product{
			Stock: 10,
		}

		err := product.ReduceStock(3)

		if err != nil {
			t.Fatalf("expected nil, got %v", err)
		}

		if product.Stock != 7 {
			t.Fatalf("expected stock 7, got %d", product.Stock)
		}
	})

	t.Run("should return invalid quantity", func(t *testing.T) {

		product := Product{
			Stock: 10,
		}

		err := product.ReduceStock(0)

		if !errors.Is(err, ErrInvalidQuantity) {
			t.Fatalf("expected %v, got %v", ErrInvalidQuantity, err)
		}
	})

	t.Run("should return insufficient stock", func(t *testing.T) {

		product := Product{
			Stock: 5,
		}

		err := product.ReduceStock(6)

		if !errors.Is(err, ErrInsufficientStock) {
			t.Fatalf("expected %v, got %v", ErrInsufficientStock, err)
		}
	})
}

func TestProduct_IncreaseStock(t *testing.T) {

	t.Run("should increase stock", func(t *testing.T) {

		product := Product{
			Stock: 5,
		}

		err := product.IncreaseStock(3)

		if err != nil {
			t.Fatalf("expected nil, got %v", err)
		}

		if product.Stock != 8 {
			t.Fatalf("expected stock 8, got %d", product.Stock)
		}
	})

	t.Run("should return invalid quantity", func(t *testing.T) {

		product := Product{
			Stock: 5,
		}

		err := product.IncreaseStock(0)

		if !errors.Is(err, ErrInvalidQuantity) {
			t.Fatalf("expected %v, got %v", ErrInvalidQuantity, err)
		}
	})
}

func TestProduct_Update(t *testing.T) {

	t.Run("should update product", func(t *testing.T) {

		product := Product{
			Name:  "Mouse",
			Price: 100,
			Stock: 5,
		}

		err := product.Update(
			"Mouse Gamer",
			150,
			10,
		)

		if err != nil {
			t.Fatalf("expected nil, got %v", err)
		}

		if product.Name != "Mouse Gamer" {
			t.Fatalf("expected Mouse Gamer, got %s", product.Name)
		}

		if product.Price != 150 {
			t.Fatalf("expected 150, got %.2f", product.Price)
		}

		if product.Stock != 10 {
			t.Fatalf("expected stock 10, got %d", product.Stock)
		}
	})

	t.Run("should return validation error", func(t *testing.T) {

		product := Product{
			Name:  "Mouse",
			Price: 100,
			Stock: 5,
		}

		err := product.Update(
			"",
			150,
			10,
		)

		if !errors.Is(err, ErrProductNameRequired) {
			t.Fatalf("expected %v, got %v", ErrProductNameRequired, err)
		}
	})
}
