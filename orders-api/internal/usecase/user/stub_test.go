package user

/*
Stub do UserRepository.

Responsabilidades:
- simular o repositório;
- controlar o retorno do método Save().
*/

import (
	"orders-api/internal/domain"
)

type userRepositoryStub struct {
	user *domain.User

	saveError   error
	findError   error
	updateError error
}

func (r *userRepositoryStub) Save(user *domain.User) error {
	return r.saveError
}

func (r *userRepositoryStub) Update(user *domain.User) error {
	return r.updateError
}

func (r *userRepositoryStub) Delete(id string) error {
	return nil
}

func (r *userRepositoryStub) FindByID(id string) (*domain.User, error) {

	if r.findError != nil {
		return nil, r.findError
	}

	return r.user, nil
}

func (r *userRepositoryStub) FindByEmail(email string) (*domain.User, error) {
	return nil, nil
}

func (r *userRepositoryStub) FindAll() ([]*domain.User, error) {
	return nil, nil
}
