package memory

/*
struct MemoryOrderRepository

Responsabilidades:

- armazenar pedidos em memória;
- implementar OrderRepository.

Métodos:
- NewMemoryOrderRepository()
- Save()
- Update()
- Delete()
- FindByID()
- FindAll()
*/

import (
	"desafio-go/internal/domain"
	"desafio-go/internal/repository"
)

type MemoryOrderRepository struct {
	orders map[string]*domain.Order
}

func NewMemoryOrderRepository() repository.OrderRepository {
	return &MemoryOrderRepository{
		orders: make(map[string]*domain.Order),
	}
}

func (r *MemoryOrderRepository) Save(order *domain.Order) error {

	if err := order.Validate(); err != nil {
		return err
	}

	r.orders[order.ID] = order

	return nil
}

func (r *MemoryOrderRepository) Update(order *domain.Order) error {

	if _, exists := r.orders[order.ID]; !exists {
		return domain.ErrOrderNotFound
	}

	if err := order.Validate(); err != nil {
		return err
	}

	r.orders[order.ID] = order

	return nil
}

func (r *MemoryOrderRepository) Delete(id string) error {

	if _, exists := r.orders[id]; !exists {
		return domain.ErrOrderNotFound
	}

	delete(r.orders, id)

	return nil
}

func (r *MemoryOrderRepository) FindByID(id string) (*domain.Order, error) {

	order, exists := r.orders[id]

	if !exists {
		return nil, domain.ErrOrderNotFound
	}

	return order, nil
}

func (r *MemoryOrderRepository) FindAll() ([]*domain.Order, error) {

	orders := make([]*domain.Order, 0, len(r.orders))

	for _, order := range r.orders {
		orders = append(orders, order)
	}

	return orders, nil
}
