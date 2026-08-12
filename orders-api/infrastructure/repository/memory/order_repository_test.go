package memory

import (
	"testing"

	"desafio-go/orders-api/internal/domain"
)

func newOrderRepository() *MemoryOrderRepository {
	return NewMemoryOrderRepository().(*MemoryOrderRepository)
}

func newValidOrder() *domain.Order {
	order := domain.NewOrder(
		"ORD001",
		"CUS001",
	)

	order.AddItem(
		*domain.NewOrderItem(
			"PROD001",
			"Notebook",
			5000,
			2,
		),
	)

	return order
}

func TestNewMemoryOrderRepository(t *testing.T) {
	repo := NewMemoryOrderRepository()

	if repo == nil {
		t.Fatal("expected repository")
	}
}

func TestOrderSave(t *testing.T) {
	repo := newOrderRepository()

	order := newValidOrder()

	err := repo.Save(order)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(repo.orders) != 1 {
		t.Fatalf("expected 1 order, got %d", len(repo.orders))
	}
}

func TestOrderSave_InvalidOrder(t *testing.T) {
	repo := newOrderRepository()

	order := domain.NewOrder(
		"",
		"",
	)

	err := repo.Save(order)

	if err == nil {
		t.Fatal("expected validation error")
	}
}

func TestOrderUpdate(t *testing.T) {
	repo := newOrderRepository()

	order := newValidOrder()

	_ = repo.Save(order)

	order.Status = domain.OrderStatusPaid

	err := repo.Update(order)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	saved, _ := repo.FindByID(order.ID)

	if saved.Status != domain.OrderStatusPaid {
		t.Fatal("order was not updated")
	}
}

func TestOrderUpdate_NotFound(t *testing.T) {
	repo := newOrderRepository()

	order := newValidOrder()

	err := repo.Update(order)

	if err != domain.ErrOrderNotFound {
		t.Fatalf("expected %v, got %v",
			domain.ErrOrderNotFound,
			err,
		)
	}
}

func TestOrderUpdate_InvalidOrder(t *testing.T) {
	repo := newOrderRepository()

	order := newValidOrder()

	_ = repo.Save(order)

	order.Items = nil

	err := repo.Update(order)

	if err == nil {
		t.Fatal("expected validation error")
	}
}

func TestOrderDelete(t *testing.T) {
	repo := newOrderRepository()

	order := newValidOrder()

	_ = repo.Save(order)

	err := repo.Delete(order.ID)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(repo.orders) != 0 {
		t.Fatal("order should have been removed")
	}
}

func TestOrderDelete_NotFound(t *testing.T) {
	repo := newOrderRepository()

	err := repo.Delete("INVALID")

	if err != domain.ErrOrderNotFound {
		t.Fatalf("expected %v, got %v",
			domain.ErrOrderNotFound,
			err,
		)
	}
}

func TestOrderFindByID(t *testing.T) {
	repo := newOrderRepository()

	order := newValidOrder()

	_ = repo.Save(order)

	found, err := repo.FindByID(order.ID)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if found.ID != order.ID {
		t.Fatal("wrong order returned")
	}
}

func TestOrderFindByID_NotFound(t *testing.T) {
	repo := newOrderRepository()

	_, err := repo.FindByID("INVALID")

	if err != domain.ErrOrderNotFound {
		t.Fatalf("expected %v, got %v",
			domain.ErrOrderNotFound,
			err,
		)
	}
}

func TestOrderFindAll(t *testing.T) {
	repo := newOrderRepository()

	_ = repo.Save(newValidOrder())

	order2 := domain.NewOrder(
		"ORD002",
		"CUS002",
	)

	order2.AddItem(
		*domain.NewOrderItem(
			"PROD002",
			"Mouse",
			100,
			1,
		),
	)

	_ = repo.Save(order2)

	orders, err := repo.FindAll()

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(orders) != 2 {
		t.Fatalf("expected 2 orders, got %d", len(orders))
	}
}

func TestOrderFindAll_Empty(t *testing.T) {
	repo := newOrderRepository()

	orders, err := repo.FindAll()

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(orders) != 0 {
		t.Fatalf("expected empty slice, got %d", len(orders))
	}
}

func TestOrderList(t *testing.T) {
	repo := newOrderRepository()

	for i := 1; i <= 5; i++ {
		order := domain.NewOrder(
			string(rune('0'+i)),
			"CUS001",
		)

		order.AddItem(
			*domain.NewOrderItem(
				"PROD001",
				"Notebook",
				100,
				1,
			),
		)

		_ = repo.Save(order)
	}

	orders, err := repo.List(2, 1)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(orders) != 2 {
		t.Fatalf("expected 2 orders, got %d", len(orders))
	}
}

func TestOrderList_OffsetGreaterThanLength(t *testing.T) {
	repo := newOrderRepository()

	orders, err := repo.List(10, 100)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(orders) != 0 {
		t.Fatalf("expected empty slice, got %d", len(orders))
	}
}
