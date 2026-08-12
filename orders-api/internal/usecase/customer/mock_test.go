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

	"desafio-go/orders-api/internal/domain"
)

type customerRepositoryMock struct {
	customer *domain.Customer

	expectedSaveCalls int
	currentSaveCalls  int

	expectedUpdateCalls int
	currentUpdateCalls  int

	saveError   error
	updateError error
}

func (r *customerRepositoryMock) Save(customer *domain.Customer) error {

	r.currentSaveCalls++

	return r.saveError
}

func (r *customerRepositoryMock) Update(customer *domain.Customer) error {

	r.currentUpdateCalls++

	return r.updateError
}

func (r *customerRepositoryMock) Delete(id string) error {
	return nil
}

func (r *customerRepositoryMock) FindByID(id string) (*domain.Customer, error) {
	return r.customer, nil
}

func (r *customerRepositoryMock) FindAll() ([]*domain.Customer, error) {
	return nil, nil
}

func (r *customerRepositoryMock) Verify() error {

	if r.currentSaveCalls != r.expectedSaveCalls {
		return errors.New("unexpected number of Save calls")
	}

	if r.currentUpdateCalls != r.expectedUpdateCalls {
		return errors.New("unexpected number of Update calls")
	}

	return nil
}
func (r *customerRepositoryMock) FindByEmail(
	email string,
) (*domain.Customer, error) {

	return r.customer, nil
}
