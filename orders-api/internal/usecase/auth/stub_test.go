package auth

/*
Stub do UserRepository.

Responsabilidades:
- simular o repositório;
- controlar o retorno do método FindByEmail().
*/

import (
	"context"
	"orders-api/internal/domain"
	"time"
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

type tokenStoreStub struct {
	saveError   error
	deleteError error
}

func (s *tokenStoreStub) Save(
	ctx context.Context,
	token string,
	ttl time.Duration,
) error {
	return s.saveError
}

func (s *tokenStoreStub) Exists(
	ctx context.Context,
	token string,
) (bool, error) {
	return true, nil
}

func (s *tokenStoreStub) Delete(
	ctx context.Context,
	token string,
) error {
	return s.deleteError
}
