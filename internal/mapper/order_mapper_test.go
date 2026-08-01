package mapper

import (
	"testing"

	domain "desafio-go/internal/domain"
	orderdto "desafio-go/internal/dto/order"
)

func TestToOrder(t *testing.T) {

	request := orderdto.CreateOrderRequest{
		CustomerID: "CUS001",
		Items: []orderdto.CreateOrderItemRequest{
			{
				ProductID: "PROD001",
				Quantity:  2,
			},
			{
				ProductID: "PROD002",
				Quantity:  3,
			},
		},
	}

	order := ToOrder(request)

	if order.ID == "" {
		t.Fatal("expected generated id")
	}

	if order.CustomerID != request.CustomerID {
		t.Fatalf("expected %s, got %s",
			request.CustomerID,
			order.CustomerID,
		)
	}

	if len(order.Items) != 2 {
		t.Fatalf("expected 2 items, got %d",
			len(order.Items),
		)
	}

	if order.Items[0].ProductID != "PROD001" {
		t.Fatal("unexpected first item")
	}

	if order.Items[0].Quantity != 2 {
		t.Fatal("unexpected quantity")
	}

	if order.Items[1].ProductID != "PROD002" {
		t.Fatal("unexpected second item")
	}
}
func TestToOrderResponse(t *testing.T) {
	order := domain.NewOrder("ORD001", "CUS001")

	item := domain.NewOrderItem(
		"PROD001",
		"Notebook",
		5000,
		2,
	)

	order.AddItem(*item)

	response := ToOrderResponse(order)

	if response.ID != order.ID {
		t.Fatal("id mismatch")
	}

	if response.CustomerID != order.CustomerID {
		t.Fatal("customer mismatch")
	}

	if response.Status != string(order.Status) {
		t.Fatal("status mismatch")
	}

	if len(response.Items) != 1 {
		t.Fatal("expected one item")
	}

	if response.Items[0].ProductID != "PROD001" {
		t.Fatal("product mismatch")
	}

	if response.Items[0].Subtotal != 10000 {
		t.Fatal("subtotal mismatch")
	}

	if response.Total != 10000 {
		t.Fatal("total mismatch")
	}
}

func TestToOrderResponseList(t *testing.T) {
	order1 := domain.NewOrder("ORD001", "CUS001")
	order2 := domain.NewOrder("ORD002", "CUS002")

	response := ToOrderResponseList([]*domain.Order{
		order1,
		order2,
	})

	if len(response) != 2 {
		t.Fatalf("expected 2 orders, got %d", len(response))
	}

	if response[0].ID != "ORD001" {
		t.Fatal("unexpected first order")
	}

	if response[1].ID != "ORD002" {
		t.Fatal("unexpected second order")
	}
}

func TestToOrderResponseList_Empty(t *testing.T) {
	response := ToOrderResponseList([]*domain.Order{})

	if len(response) != 0 {
		t.Fatalf("expected empty list, got %d", len(response))
	}
}
