package order

/*
FakeProductRepository.

Responsabilidades:
- simular o ProductRepository;
- controlar retornos;
- registrar chamadas.

Métodos:
- Save()
- Update()
- Delete()
- FindByID()
- FindAll()
*/

import (
	"orders-api/internal/domain"
)

type FakeProductRepository struct {
	product *domain.Product

	findError   error
	updateError error
	saveError   error

	findCalls   int
	updateCalls int
	saveCalls   int

	spy  *OrderSpy
	mock *OrderMock
}

func (r *FakeProductRepository) Save(product *domain.Product) error {

	r.saveCalls++

	return r.saveError
}

func (r *FakeProductRepository) Update(product *domain.Product) error {

	r.updateCalls++

	if r.spy != nil {
		r.spy.Add("product.update")
	}

	if r.mock != nil {
		r.mock.ProductUpdate++
	}

	return r.updateError
}

func (r *FakeProductRepository) Delete(id string) error {
	return nil
}

func (r *FakeProductRepository) FindByID(id string) (*domain.Product, error) {

	r.findCalls++

	if r.spy != nil {
		r.spy.Add("product.find")
	}

	if r.mock != nil {
		r.mock.ProductFind++
	}

	if r.findError != nil {
		return nil, r.findError
	}

	return r.product, nil
}

func (r *FakeProductRepository) FindAll() ([]*domain.Product, error) {
	return nil, nil
}
