package product

/*
Spy do ProductRepository.

Responsabilidades:
- registrar chamadas ao repositório;
- armazenar parâmetros recebidos;
- contabilizar chamadas.
*/

import (
	"desafio-go/internal/domain"
)

type productRepositorySpy struct {
	saveCalled bool
	saveCalls  int
	product    *domain.Product
}

func (r *productRepositorySpy) Save(product *domain.Product) error {

	r.saveCalled = true
	r.saveCalls++
	r.product = product

	return nil
}

func (r *productRepositorySpy) Update(product *domain.Product) error {
	return nil
}

func (r *productRepositorySpy) Delete(id string) error {
	return nil
}

func (r *productRepositorySpy) FindByID(id string) (*domain.Product, error) {
	return nil, nil
}

func (r *productRepositorySpy) FindAll() ([]*domain.Product, error) {
	return nil, nil
}
