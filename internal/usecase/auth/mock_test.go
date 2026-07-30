package auth

/*
Mock do UserRepository.

Responsabilidades:
- validar chamadas esperadas;
- simular erros;
- verificar quantidade de execuções.
*/

import (
	"errors"

	"desafio-go/internal/domain"
)

type userRepositoryMock struct {
	expectedCalls int
	currentCalls  int

	user             *domain.User
	findByEmailError error
}

func (r *userRepositoryMock) Save(user *domain.User) error {
	return nil
}

func (r *userRepositoryMock) Update(user *domain.User) error {
	return nil
}

func (r *userRepositoryMock) Delete(id string) error {
	return nil
}

func (r *userRepositoryMock) FindByID(id string) (*domain.User, error) {
	return nil, nil
}

func (r *userRepositoryMock) FindByEmail(email string) (*domain.User, error) {

	r.currentCalls++

	if r.findByEmailError != nil {
		return nil, r.findByEmailError
	}

	return r.user, nil
}

func (r *userRepositoryMock) FindAll() ([]*domain.User, error) {
	return nil, nil
}

func (r *userRepositoryMock) Verify() error {

	if r.currentCalls != r.expectedCalls {
		return errors.New("unexpected number of FindByEmail calls")
	}

	return nil
}
