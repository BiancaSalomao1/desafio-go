package customer

/*
Spy do CustomerRepository.

Responsabilidades:
- registrar chamadas ao repositório;
- armazenar parâmetros recebidos;
- contabilizar chamadas.
*/

import (
	"desafio-go/internal/domain"
)

type customerRepositorySpy struct {
	saveCalled bool
	saveCalls  int
	customer   *domain.Customer
}

func (r *customerRepositorySpy) Save(customer *domain.Customer) error {

	r.saveCalled = true
	r.saveCalls++
	r.customer = customer

	return nil
}

func (r *customerRepositorySpy) Update(customer *domain.Customer) error {
	return nil
}

func (r *customerRepositorySpy) Delete(id string) error {
	return nil
}

func (r *customerRepositorySpy) FindByID(id string) (*domain.Customer, error) {
	return nil, nil
}

func (r *customerRepositorySpy) FindAll() ([]*domain.Customer, error) {
	return nil, nil
}
