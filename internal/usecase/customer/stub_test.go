package customer

/*
Stub do CustomerRepository.

Responsabilidades:
- simular o repositório;
- controlar o retorno do método Save().
*/

import (
	"desafio-go/internal/domain"
)

type customerRepositoryStub struct {
	saveError error
}

func (r *customerRepositoryStub) Save(customer *domain.Customer) error {
	return r.saveError
}

func (r *customerRepositoryStub) Update(customer *domain.Customer) error {
	return nil
}

func (r *customerRepositoryStub) Delete(id string) error {
	return nil
}

func (r *customerRepositoryStub) FindByID(id string) (*domain.Customer, error) {
	return nil, nil
}

func (r *customerRepositoryStub) FindAll() ([]*domain.Customer, error) {
	return nil, nil
}
