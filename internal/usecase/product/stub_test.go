package product

/*
Stub do ProductRepository.

Responsabilidades:
- simular o repositório;
- controlar o retorno do método Save().
*/

import (
	"desafio-go/internal/domain"
)

type productRepositoryStub struct {
	saveError error
}

func (r *productRepositoryStub) Save(product *domain.Product) error {
	return r.saveError
}

func (r *productRepositoryStub) Update(product *domain.Product) error {
	return nil
}

func (r *productRepositoryStub) Delete(id string) error {
	return nil
}

func (r *productRepositoryStub) FindByID(id string) (*domain.Product, error) {
	return nil, nil
}

func (r *productRepositoryStub) FindAll() ([]*domain.Product, error) {
	return nil, nil
}
