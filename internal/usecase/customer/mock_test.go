package customer

/*
Mock do CustomerRepository.

Responsabilidades:
- validar chamadas esperadas;
- simular erros;
- verificar quantidade de execuções.
*/

import (
	"errors"

	"desafio-go/internal/domain"
)

type customerRepositoryMock struct {
	expectedCalls int
	currentCalls  int

	saveError error
}

func (r *customerRepositoryMock) Save(customer *domain.Customer) error {

	r.currentCalls++

	return r.saveError
}

func (r *customerRepositoryMock) Update(customer *domain.Customer) error {
	return nil
}

func (r *customerRepositoryMock) Delete(id string) error {
	return nil
}

func (r *customerRepositoryMock) FindByID(id string) (*domain.Customer, error) {
	return nil, nil
}

func (r *customerRepositoryMock) FindAll() ([]*domain.Customer, error) {
	return nil, nil
}

func (r *customerRepositoryMock) Verify() error {

	if r.currentCalls != r.expectedCalls {
		return errors.New("unexpected number of Save calls")
	}

	return nil
}
