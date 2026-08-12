package user

/*
Spy do UserRepository.

Responsabilidades:
- registrar chamadas ao repositório;
- armazenar parâmetros recebidos;
- contabilizar chamadas.
*/

import (
	"desafio-go/orders-api/internal/domain"
)

type userRepositorySpy struct {
	saveCalled bool
	saveCalls  int

	updateCalled bool
	updateCalls  int

	user *domain.User
}

func (r *userRepositorySpy) Save(user *domain.User) error {

	r.saveCalled = true
	r.saveCalls++
	r.user = user

	return nil
}

func (r *userRepositorySpy) Update(user *domain.User) error {

	r.updateCalled = true
	r.updateCalls++
	r.user = user

	return nil
}

func (r *userRepositorySpy) Delete(id string) error {
	return nil
}

func (r *userRepositorySpy) FindByID(id string) (*domain.User, error) {
	return r.user, nil
}

func (r *userRepositorySpy) FindByEmail(email string) (*domain.User, error) {
	return nil, nil
}

func (r *userRepositorySpy) FindAll() ([]*domain.User, error) {
	return nil, nil
}
