package auth

import (
	"context"
	"errors"
	"testing"
	"time"

	"orders-api/internal/domain"
	"orders-api/internal/security"
)

func TestLoginUseCase_Execute(t *testing.T) {

	newUser := func() *domain.User {

		hash, err := security.HashPassword("123456")
		if err != nil {
			t.Fatalf("failed to hash password: %v", err)
		}

		return domain.NewUser(
			"1",
			"Administrador",
			"admin@email.com",
			hash,
		)
	}

	t.Run("should login successfully", func(t *testing.T) {

		repository := &userRepositoryStub{
			user: newUser(),
		}

		tokenStore := &tokenStoreStub{}

		useCase := NewLoginUseCase(
			repository,
			tokenStore,
			"secret",
			time.Hour,
		)

		token, err := useCase.Execute(
			context.Background(),
			"admin@email.com",
			"123456",
		)

		if err != nil {
			t.Fatalf("expected nil, got %v", err)
		}

		if token == "" {
			t.Fatal("expected token")
		}
	})

	t.Run("should return repository error", func(t *testing.T) {

		repository := &userRepositoryStub{
			findByEmailError: errors.New("user not found"),
		}

		tokenStore := &tokenStoreStub{}

		useCase := NewLoginUseCase(
			repository,
			tokenStore,
			"secret",
			time.Hour,
		)

		_, err := useCase.Execute(
			context.Background(),
			"admin@email.com",
			"123456",
		)

		if err == nil {
			t.Fatal("expected repository error")
		}

		if err.Error() != "user not found" {
			t.Fatalf("expected user not found, got %v", err)
		}
	})

	t.Run("should return invalid credentials", func(t *testing.T) {

		repository := &userRepositoryStub{
			user: newUser(),
		}

		tokenStore := &tokenStoreStub{}

		useCase := NewLoginUseCase(
			repository,
			tokenStore,
			"secret",
			time.Hour,
		)

		_, err := useCase.Execute(
			context.Background(),
			"admin@email.com",
			"wrong-password",
		)

		if !errors.Is(err, domain.ErrInvalidCredentials) {
			t.Fatalf(
				"expected %v, got %v",
				domain.ErrInvalidCredentials,
				err,
			)
		}
	})

	t.Run("should call FindByEmail once", func(t *testing.T) {

		repository := &userRepositorySpy{
			user: newUser(),
		}

		tokenStore := &tokenStoreStub{}

		useCase := NewLoginUseCase(
			repository,
			tokenStore,
			"secret",
			time.Hour,
		)

		_, err := useCase.Execute(
			context.Background(),
			"admin@email.com",
			"123456",
		)

		if err != nil {
			t.Fatalf("expected nil, got %v", err)
		}

		if !repository.findByEmailCalled {
			t.Fatal("expected FindByEmail to be called")
		}

		if repository.findByEmailCalls != 1 {
			t.Fatalf(
				"expected 1 call, got %d",
				repository.findByEmailCalls,
			)
		}

		if repository.email != "admin@email.com" {
			t.Fatalf(
				"expected admin@email.com, got %s",
				repository.email,
			)
		}
	})

	t.Run("should satisfy mock expectations", func(t *testing.T) {

		repository := &userRepositoryMock{
			expectedCalls: 1,
			user:          newUser(),
		}

		tokenStore := &tokenStoreStub{}

		useCase := NewLoginUseCase(
			repository,
			tokenStore,
			"secret",
			time.Hour,
		)

		_, err := useCase.Execute(
			context.Background(),
			"admin@email.com",
			"123456",
		)

		if err != nil {
			t.Fatalf("expected nil, got %v", err)
		}

		if err := repository.Verify(); err != nil {
			t.Fatal(err)
		}
	})
}
