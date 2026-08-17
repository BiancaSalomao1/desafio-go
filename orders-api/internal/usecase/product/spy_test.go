package product

/*
Spy do ProductRepository.

Responsabilidades:
- registrar chamadas ao repositório;
- armazenar parâmetros recebidos;
- contabilizar chamadas.
*/

import (
	"orders-api/internal/domain"
)

type productRepositorySpy struct {
	saveCalled bool
	saveCalls  int

	updateCalled bool
	updateCalls  int

	product *domain.Product
}

func (r *productRepositorySpy) Save(product *domain.Product) error {
	r.saveCalled = true
	r.saveCalls++
	r.product = product
	return nil
}

func (r *productRepositorySpy) Update(product *domain.Product) error {
	r.updateCalled = true
	r.updateCalls++
	r.product = product
	return nil
}

func (r *productRepositorySpy) Delete(id string) error {
	return nil
}

func (r *productRepositorySpy) FindByID(id string) (*domain.Product, error) {
	return r.product, nil
}

func (r *productRepositorySpy) FindAll() ([]*domain.Product, error) {
	return nil, nil
}
