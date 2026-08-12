package user

import (
	"errors"
	"testing"

	"desafio-go/orders-api/internal/domain"
)

func TestUpdateUserUseCase_Execute(t *testing.T) {

	t.Run("should update user successfully", func(t *testing.T) {

		repository := &userRepositoryStub{
			user: &domain.User{
				ID:           "1",
				Name:         "Administrador",
				Email:        "admin@email.com",
				PasswordHash: "hash",
			},
		}

		useCase := NewUpdateUserUseCase(repository)

		err := useCase.Execute(&domain.User{
			ID:           "1",
			Name:         "Novo Nome",
			Email:        "novo@email.com",
			PasswordHash: "hash",
		})

		if err != nil {
			t.Fatalf("expected nil, got %v", err)
		}

		if repository.user.Name != "Novo Nome" {
			t.Fatal("expected updated name")
		}

		if repository.user.Email != "novo@email.com" {
			t.Fatal("expected updated email")
		}
	})

	t.Run("should return repository find error", func(t *testing.T) {

		repository := &userRepositoryStub{
			findError: errors.New("user not found"),
		}

		useCase := NewUpdateUserUseCase(repository)

		err := useCase.Execute(&domain.User{ID: "1"})

		if err == nil {
			t.Fatal("expected repository error")
		}
	})

	t.Run("should return validation error", func(t *testing.T) {

		repository := &userRepositoryStub{
			user: &domain.User{
				ID:           "1",
				Name:         "Administrador",
				Email:        "admin@email.com",
				PasswordHash: "hash",
			},
		}

		useCase := NewUpdateUserUseCase(repository)

		err := useCase.Execute(&domain.User{
			ID:           "1",
			Name:         "",
			Email:        "",
			PasswordHash: "hash",
		})

		if err == nil {
			t.Fatal("expected validation error")
		}
	})

	t.Run("should return update error", func(t *testing.T) {

		repository := &userRepositoryStub{
			user: &domain.User{
				ID:           "1",
				Name:         "Administrador",
				Email:        "admin@email.com",
				PasswordHash: "hash",
			},
			updateError: errors.New("database error"),
		}

		useCase := NewUpdateUserUseCase(repository)

		err := useCase.Execute(&domain.User{
			ID:           "1",
			Name:         "Novo Nome",
			Email:        "novo@email.com",
			PasswordHash: "hash",
		})

		if err == nil {
			t.Fatal("expected update error")
		}
	})

	t.Run("should call update once", func(t *testing.T) {

		repository := &userRepositorySpy{
			user: &domain.User{
				ID:           "1",
				Name:         "Administrador",
				Email:        "admin@email.com",
				PasswordHash: "hash",
			},
		}

		useCase := NewUpdateUserUseCase(repository)

		_ = useCase.Execute(&domain.User{
			ID:           "1",
			Name:         "Novo Nome",
			Email:        "novo@email.com",
			PasswordHash: "hash",
		})

		if !repository.updateCalled {
			t.Fatal("expected Update to be called")
		}

		if repository.updateCalls != 1 {
			t.Fatalf("expected 1 call, got %d", repository.updateCalls)
		}
	})

	t.Run("should satisfy mock expectations", func(t *testing.T) {

		repository := &userRepositoryMock{
			user: &domain.User{
				ID:           "1",
				Name:         "Administrador",
				Email:        "admin@email.com",
				PasswordHash: "hash",
			},
			expectedUpdateCalls: 1,
		}

		useCase := NewUpdateUserUseCase(repository)

		err := useCase.Execute(&domain.User{
			ID:           "1",
			Name:         "Novo Nome",
			Email:        "novo@email.com",
			PasswordHash: "hash",
		})

		if err != nil {
			t.Fatalf("expected nil, got %v", err)
		}

		if err := repository.Verify(); err != nil {
			t.Fatal(err)
		}
	})
}
