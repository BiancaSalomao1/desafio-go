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

	"desafio-go/orders-api/internal/domain"
)

type productRepositoryMock struct {
	product *domain.Product

	expectedSaveCalls int
	currentSaveCalls  int

	expectedUpdateCalls int
	currentUpdateCalls  int

	saveError   error
	updateError error
}

func (r *productRepositoryMock) Save(product *domain.Product) error {
	r.currentSaveCalls++
	return r.saveError
}

func (r *productRepositoryMock) Update(product *domain.Product) error {
	r.currentUpdateCalls++
	return r.updateError
}

func (r *productRepositoryMock) Delete(id string) error {
	return nil
}

func (r *productRepositoryMock) FindByID(id string) (*domain.Product, error) {
	return r.product, nil
}

func (r *productRepositoryMock) FindAll() ([]*domain.Product, error) {
	return nil, nil
}

func (r *productRepositoryMock) Verify() error {

	if r.currentSaveCalls != r.expectedSaveCalls {
		return errors.New("unexpected number of Save calls")
	}

	if r.currentUpdateCalls != r.expectedUpdateCalls {
		return errors.New("unexpected number of Update calls")
	}

	return nil
}
