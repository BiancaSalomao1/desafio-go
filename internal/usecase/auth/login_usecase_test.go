package auth

/*
Testes do LoginUseCase.

Responsabilidades:
- validar autenticação;
- validar credenciais;
- validar interação com o repositório.

Cenários:
- login realizado com sucesso;
- usuário não encontrado;
- senha inválida;
- chamada do FindByEmail();
- validação do Mock.
*/

import (
	"errors"
	"testing"
	"time"

	"desafio-go/internal/domain"
	"desafio-go/internal/security"
)

func TestLoginUseCase_Execute(t *testing.T) {

	newUser := func() *domain.User {

		hash, err := security.HashPassword("123456")
		if err != nil {
			t.Fatalf("failed to hash password: %v", err)
		}

		return &domain.User{
			ID:           "1",
			Name:         "Administrador",
			Email:        "admin@email.com",
			PasswordHash: hash,
		}
	}

	t.Run("should login successfully", func(t *testing.T) {

		repository := &userRepositoryStub{
			user: newUser(),
		}

		useCase := NewLoginUseCase(
			repository,
			"secret",
			time.Hour,
		)

		token, err := useCase.Execute(
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

		useCase := NewLoginUseCase(
			repository,
			"secret",
			time.Hour,
		)

		_, err := useCase.Execute(
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

		useCase := NewLoginUseCase(
			repository,
			"secret",
			time.Hour,
		)

		_, err := useCase.Execute(
			"admin@email.com",
			"senha-incorreta",
		)

		if err == nil {
			t.Fatal("expected invalid credentials")
		}

		if err.Error() != "invalid credentials" {
			t.Fatalf("expected invalid credentials, got %v", err)
		}
	})

	t.Run("should call FindByEmail once", func(t *testing.T) {

		repository := &userRepositorySpy{
			user: newUser(),
		}

		useCase := NewLoginUseCase(
			repository,
			"secret",
			time.Hour,
		)

		_, err := useCase.Execute(
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

		useCase := NewLoginUseCase(
			repository,
			"secret",
			time.Hour,
		)

		_, err := useCase.Execute(
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
