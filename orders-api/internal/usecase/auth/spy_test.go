package auth

/*
Spy do UserRepository.

Responsabilidades:
- registrar chamadas ao repositório;
- armazenar parâmetros recebidos;
- contabilizar chamadas.
*/

import (
	"context"
	"orders-api/internal/domain"
	"time"
)

type userRepositorySpy struct {
	findByEmailCalled bool
	findByEmailCalls  int
	email             string

	user *domain.User
}

func (r *userRepositorySpy) Save(user *domain.User) error {
	return nil
}

func (r *userRepositorySpy) Update(user *domain.User) error {
	return nil
}

func (r *userRepositorySpy) Delete(id string) error {
	return nil
}

func (r *userRepositorySpy) FindByID(id string) (*domain.User, error) {
	return nil, nil
}

func (r *userRepositorySpy) FindByEmail(email string) (*domain.User, error) {

	r.findByEmailCalled = true
	r.findByEmailCalls++
	r.email = email

	return r.user, nil
}

func (r *userRepositorySpy) FindAll() ([]*domain.User, error) {
	return nil, nil
}

type tokenStoreSpy struct {
	deleteCalled bool
	deleteCalls  int
	token        string

	deleteError error
}

func (s *tokenStoreSpy) Save(
	ctx context.Context,
	token string,
	ttl time.Duration,
) error {
	return nil
}

func (s *tokenStoreSpy) Exists(
	ctx context.Context,
	token string,
) (bool, error) {
	return true, nil
}

func (s *tokenStoreSpy) Delete(
	ctx context.Context,
	token string,
) error {
	s.deleteCalled = true
	s.deleteCalls++
	s.token = token

	return s.deleteError
}
