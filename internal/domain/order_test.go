package domain

import (
	"errors"
	"testing"
)

func validOrderItem() OrderItem {
	return OrderItem{
		ProductID: "p1",
		Name:      "Notebook",
		Price:     100,
		Quantity:  2,
	}
}

func TestNewOrder(t *testing.T) {

	order := NewOrder("1", "c1")

	if order == nil {
		t.Fatal("expected order")
	}

	if order.ID != "1" {
		t.Fatalf("expected id 1, got %s", order.ID)
	}

	if order.CustomerID != "c1" {
		t.Fatalf("expected customer c1, got %s", order.CustomerID)
	}

	if order.Status != OrderStatusPending {
		t.Fatalf("expected pending, got %v", order.Status)
	}

	if len(order.Items) != 0 {
		t.Fatal("expected empty items")
	}
}

func TestOrder_AddItem(t *testing.T) {

	t.Run("should add item", func(t *testing.T) {

		order := NewOrder("1", "c1")

		item := validOrderItem()

		err := order.AddItem(item)

		if err != nil {
			t.Fatalf("expected nil, got %v", err)
		}

		if len(order.Items) != 1 {
			t.Fatalf("expected 1 item, got %d", len(order.Items))
		}
	})

	t.Run("should return validation error", func(t *testing.T) {

		order := NewOrder("1", "c1")

		item := OrderItem{
			ProductID: "",
			Quantity:  1,
		}

		err := order.AddItem(item)

		if !errors.Is(err, ErrProductInvalid) {
			t.Fatalf("expected %v, got %v", ErrProductInvalid, err)
		}
	})
}

func TestOrder_RemoveItem(t *testing.T) {

	t.Run("should remove item", func(t *testing.T) {

		order := NewOrder("1", "c1")

		order.Items = append(order.Items, validOrderItem())

		err := order.RemoveItem("p1")

		if err != nil {
			t.Fatalf("expected nil, got %v", err)
		}

		if len(order.Items) != 0 {
			t.Fatal("expected empty order")
		}
	})

	t.Run("should return product not found", func(t *testing.T) {

		order := NewOrder("1", "c1")

		err := order.RemoveItem("x")

		if !errors.Is(err, ErrProductNotFound) {
			t.Fatalf("expected %v, got %v", ErrProductNotFound, err)
		}
	})
}

func TestOrder_Total(t *testing.T) {

	order := NewOrder("1", "c1")

	order.Items = append(order.Items,
		OrderItem{
			ProductID: "1",
			Price:     100,
			Quantity:  2,
		},
		OrderItem{
			ProductID: "2",
			Price:     50,
			Quantity:  3,
		},
	)

	total := order.Total()

	if total != 350 {
		t.Fatalf("expected 350, got %.2f", total)
	}
}

func TestOrder_Pay(t *testing.T) {

	t.Run("should pay order", func(t *testing.T) {

		order := NewOrder("1", "c1")

		err := order.Pay()

		if err != nil {
			t.Fatalf("expected nil, got %v", err)
		}

		if order.Status != OrderStatusPaid {
			t.Fatalf("expected paid, got %v", order.Status)
		}
	})

	t.Run("should return invalid status", func(t *testing.T) {

		order := NewOrder("1", "c1")
		order.Status = OrderStatusCanceled

		err := order.Pay()

		if !errors.Is(err, ErrOrderStatusInvalid) {
			t.Fatalf("expected %v, got %v", ErrOrderStatusInvalid, err)
		}
	})
}

func TestOrder_Cancel(t *testing.T) {

	t.Run("should cancel order", func(t *testing.T) {

		order := NewOrder("1", "c1")

		err := order.Cancel()

		if err != nil {
			t.Fatalf("expected nil, got %v", err)
		}

		if order.Status != OrderStatusCanceled {
			t.Fatalf("expected canceled, got %v", order.Status)
		}
	})

	t.Run("should return invalid status", func(t *testing.T) {

		order := NewOrder("1", "c1")
		order.Status = OrderStatusPaid

		err := order.Cancel()

		if !errors.Is(err, ErrOrderStatusInvalid) {
			t.Fatalf("expected %v, got %v", ErrOrderStatusInvalid, err)
		}
	})
}

func TestOrder_Validate(t *testing.T) {

	tests := []struct {
		name        string
		order       Order
		expectedErr error
	}{
		{
			name: "valid order",
			order: Order{
				CustomerID: "c1",
				Items: []OrderItem{
					validOrderItem(),
				},
			},
			expectedErr: nil,
		},
		{
			name: "empty customer",
			order: Order{
				CustomerID: "",
				Items: []OrderItem{
					validOrderItem(),
				},
			},
			expectedErr: ErrCustomerInvalid,
		},
		{
			name: "empty items",
			order: Order{
				CustomerID: "c1",
				Items:      []OrderItem{},
			},
			expectedErr: ErrEmptyOrder,
		},
	}

	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {

			err := tt.order.Validate()

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
