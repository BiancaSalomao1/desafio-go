package product

/*
Stub do ProductRepository.

Responsabilidades:
- simular o repositório;
- controlar o retorno do método Save().
*/

import (
	"desafio-go/orders-api/internal/domain"
)

type productRepositoryStub struct {
	product *domain.Product

	saveError   error
	findError   error
	updateError error
}

func (r *productRepositoryStub) Save(product *domain.Product) error {
	return r.saveError
}

func (r *productRepositoryStub) Update(product *domain.Product) error {
	return r.updateError
}

func (r *productRepositoryStub) Delete(id string) error {
	return nil
}

func (r *productRepositoryStub) FindByID(id string) (*domain.Product, error) {

	if r.findError != nil {
		return nil, r.findError
	}

	return r.product, nil
}

func (r *productRepositoryStub) FindAll() ([]*domain.Product, error) {
	return nil, nil
}
