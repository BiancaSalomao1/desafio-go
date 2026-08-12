package auth

/*
Stub do UserRepository.

Responsabilidades:
- simular o repositório;
- controlar o retorno do método FindByEmail().
*/

import (
	"desafio-go/orders-api/internal/domain"
)

type userRepositoryStub struct {
	user             *domain.User
	findByEmailError error
}

func (r *userRepositoryStub) Save(user *domain.User) error {
	return nil
}

func (r *userRepositoryStub) Update(user *domain.User) error {
	return nil
}

func (r *userRepositoryStub) Delete(id string) error {
	return nil
}

func (r *userRepositoryStub) FindByID(id string) (*domain.User, error) {
	return nil, nil
}

func (r *userRepositoryStub) FindByEmail(email string) (*domain.User, error) {

	if r.findByEmailError != nil {
		return nil, r.findByEmailError
	}

	return r.user, nil
}

func (r *userRepositoryStub) FindAll() ([]*domain.User, error) {
	return nil, nil
}
