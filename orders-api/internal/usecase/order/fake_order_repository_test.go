package order

/*
FakeOrderRepository.

Responsabilidades:
- simular o OrderRepository;
- controlar retornos;
- registrar chamadas.

Métodos:
- Save()
- Update()
- Delete()
- FindByID()
- FindAll()
- List()
*/

import (
	"desafio-go/orders-api/internal/domain"
)

type FakeOrderRepository struct {
	order *domain.Order

	saveError error

	saveCalls int

	spy  *OrderSpy
	mock *OrderMock
}

func (r *FakeOrderRepository) Save(order *domain.Order) error {

	r.saveCalls++
	r.order = order

	if r.spy != nil {
		r.spy.Add("order.save")
	}

	if r.mock != nil {
		r.mock.OrderSave++
	}

	return r.saveError
}

func (r *FakeOrderRepository) Update(order *domain.Order) error {
	return nil
}

func (r *FakeOrderRepository) Delete(id string) error {
	return nil
}

func (r *FakeOrderRepository) FindByID(id string) (*domain.Order, error) {
	return nil, nil
}

func (r *FakeOrderRepository) FindAll() ([]*domain.Order, error) {
	return nil, nil
}

func (r *FakeOrderRepository) List(limit, offset int) ([]*domain.Order, error) {
	return nil, nil
}
