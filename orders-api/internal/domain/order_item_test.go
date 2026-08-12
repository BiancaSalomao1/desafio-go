package domain

import (
	"errors"
	"testing"
)

func TestNewOrderItem(t *testing.T) {

	item := NewOrderItem(
		"p1",
		"Notebook",
		3500,
		2,
	)

	if item == nil {
		t.Fatal("expected order item")
	}

	if item.ProductID != "p1" {
		t.Fatalf("expected product id p1, got %s", item.ProductID)
	}

	if item.Name != "Notebook" {
		t.Fatalf("expected Notebook, got %s", item.Name)
	}

	if item.Price != 3500 {
		t.Fatalf("expected 3500, got %f", item.Price)
	}

	if item.Quantity != 2 {
		t.Fatalf("expected quantity 2, got %d", item.Quantity)
	}
}

func TestOrderItem_Validate(t *testing.T) {

	tests := []struct {
		name        string
		item        OrderItem
		expectedErr error
	}{
		{
			name: "valid item",
			item: OrderItem{
				ProductID: "p1",
				Name:      "Notebook",
				Price:     3500,
				Quantity:  2,
			},
			expectedErr: nil,
		},
		{
			name: "empty product id",
			item: OrderItem{
				ProductID: "",
				Name:      "Notebook",
				Price:     3500,
				Quantity:  2,
			},
			expectedErr: ErrProductInvalid,
		},
		{
			name: "invalid quantity",
			item: OrderItem{
				ProductID: "p1",
				Name:      "Notebook",
				Price:     3500,
				Quantity:  0,
			},
			expectedErr: ErrInvalidQuantity,
		},
	}

	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {

			err := tt.item.Validate()

			if !errors.Is(err, tt.expectedErr) {
				t.Fatalf(
					"expected %v, got %v",
					tt.expectedErr,
					err,
				)
			}
		})
	}
}

func TestOrderItem_Subtotal(t *testing.T) {

	tests := []struct {
		name     string
		item     OrderItem
		expected float64
	}{
		{
			name: "quantity one",
			item: OrderItem{
				Price:    10,
				Quantity: 1,
			},
			expected: 10,
		},
		{
			name: "quantity many",
			item: OrderItem{
				Price:    25.5,
				Quantity: 4,
			},
			expected: 102,
		},
		{
			name: "zero price",
			item: OrderItem{
				Price:    0,
				Quantity: 5,
			},
			expected: 0,
		},
	}

	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {

			result := tt.item.Subtotal()

			if result != tt.expected {
				t.Fatalf(
					"expected %.2f, got %.2f",
					tt.expected,
					result,
				)
			}
		})
	}
}
