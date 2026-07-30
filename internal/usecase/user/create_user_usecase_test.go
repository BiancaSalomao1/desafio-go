package user

/*
Testes do CreateUserUseCase.

Responsabilidades:
- validar criação de usuário;
- validar regras de negócio;
- validar geração do hash da senha;
- validar interação com o repositório.

Cenários:
- criação com sucesso;
- usuário inválido;
- erro do repositório;
- chamada do Save();
- geração do hash;
- validação do Mock.
*/

import (
	"errors"
	"testing"

	"desafio-go/internal/domain"
)

func TestCreateUserUseCase_Execute(t *testing.T) {

	t.Run("should create user successfully", func(t *testing.T) {

		repository := &userRepositoryStub{}

		useCase := NewCreateUserUseCase(repository)

		user := &domain.User{
			Name:         "Administrador",
			Email:        "admin@email.com",
			PasswordHash: "123456",
		}

		err := useCase.Execute(user)

		if err != nil {
			t.Fatalf("expected nil, got %v", err)
		}

		if user.PasswordHash == "123456" {
			t.Fatal("expected password to be hashed")
		}
	})

	t.Run("should return validation error", func(t *testing.T) {

		repository := &userRepositoryStub{}

		useCase := NewCreateUserUseCase(repository)

		user := &domain.User{}

		err := useCase.Execute(user)

		if err == nil {
			t.Fatal("expected validation error")
		}
	})

	t.Run("should return repository error", func(t *testing.T) {

		repository := &userRepositoryStub{
			saveError: errors.New("database error"),
		}

		useCase := NewCreateUserUseCase(repository)

		user := &domain.User{
			Name:         "Administrador",
			Email:        "admin@email.com",
			PasswordHash: "123456",
		}

		err := useCase.Execute(user)

		if err == nil {
			t.Fatal("expected repository error")
		}

		if err.Error() != "database error" {
			t.Fatalf("expected database error, got %v", err)
		}
	})

	t.Run("should call save once", func(t *testing.T) {

		repository := &userRepositorySpy{}

		useCase := NewCreateUserUseCase(repository)

		user := &domain.User{
			Name:         "Administrador",
			Email:        "admin@email.com",
			PasswordHash: "123456",
		}

		err := useCase.Execute(user)

		if err != nil {
			t.Fatalf("expected nil, got %v", err)
		}

		if !repository.saveCalled {
			t.Fatal("expected Save to be called")
		}

		if repository.saveCalls != 1 {
			t.Fatalf("expected 1 call, got %d", repository.saveCalls)
		}

		if repository.user != user {
			t.Fatal("expected same user instance")
		}
	})

	t.Run("should satisfy mock expectations", func(t *testing.T) {

		repository := &userRepositoryMock{
			expectedSaveCalls: 1,
		}

		useCase := NewCreateUserUseCase(repository)

		user := &domain.User{
			Name:         "Administrador",
			Email:        "admin@email.com",
			PasswordHash: "123456",
		}

		err := useCase.Execute(user)

		if err != nil {
			t.Fatalf("expected nil, got %v", err)
		}

		if err := repository.Verify(); err != nil {
			t.Fatal(err)
		}
	})
}
