package product

/*
Mock do ProductRepository.

Responsabilidades:
- validar chamadas esperadas;
- simular erros;
- verificar quantidade de execuções.
*/

import (
	"errors"

	"desafio-go/internal/domain"
)

type productRepositoryMock struct {
	expectedCalls int
	currentCalls  int

	saveError error
}

func (r *productRepositoryMock) Save(product *domain.Product) error {

	r.currentCalls++

	return r.saveError
}

func (r *productRepositoryMock) Update(product *domain.Product) error {
	return nil
}

func (r *productRepositoryMock) Delete(id string) error {
	return nil
}

func (r *productRepositoryMock) FindByID(id string) (*domain.Product, error) {
	return nil, nil
}

func (r *productRepositoryMock) FindAll() ([]*domain.Product, error) {
	return nil, nil
}

func (r *productRepositoryMock) Verify() error {

	if r.currentCalls != r.expectedCalls {
		return errors.New("unexpected number of Save calls")
	}

	return nil
}
