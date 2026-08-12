package user

/*
Mock do UserRepository.

Responsabilidades:
- validar chamadas esperadas;
- simular erros;
- verificar quantidade de execuções.
*/

import (
	"errors"

	"desafio-go/orders-api/internal/domain"
)

type userRepositoryMock struct {
	user *domain.User

	expectedSaveCalls int
	currentSaveCalls  int

	expectedUpdateCalls int
	currentUpdateCalls  int

	saveError   error
	updateError error
}

func (r *userRepositoryMock) Save(user *domain.User) error {

	r.currentSaveCalls++

	return r.saveError
}

func (r *userRepositoryMock) Update(user *domain.User) error {

	r.currentUpdateCalls++

	return r.updateError
}

func (r *userRepositoryMock) Delete(id string) error {
	return nil
}

func (r *userRepositoryMock) FindByID(id string) (*domain.User, error) {
	return r.user, nil
}

func (r *userRepositoryMock) FindByEmail(email string) (*domain.User, error) {
	return nil, nil
}

func (r *userRepositoryMock) FindAll() ([]*domain.User, error) {
	return nil, nil
}

func (r *userRepositoryMock) Verify() error {

	if r.currentSaveCalls != r.expectedSaveCalls {
		return errors.New("unexpected number of Save calls")
	}

	if r.currentUpdateCalls != r.expectedUpdateCalls {
		return errors.New("unexpected number of Update calls")
	}

	return nil
}
