package domain

import (
	"errors"
	"testing"

	"desafio-go/internal/security"
)

func TestNewUser(t *testing.T) {

	user := NewUser(
		"1",
		"Administrador",
		"admin@email.com",
		"hash",
	)

	if user == nil {
		t.Fatal("expected user")
	}

	if user.ID != "1" {
		t.Fatalf("expected id 1, got %s", user.ID)
	}

	if user.Name != "Administrador" {
		t.Fatalf("expected Administrador, got %s", user.Name)
	}

	if user.Email != "admin@email.com" {
		t.Fatalf("expected admin@email.com, got %s", user.Email)
	}

	if user.PasswordHash != "hash" {
		t.Fatal("expected password hash")
	}
}

func TestUser_Validate(t *testing.T) {

	tests := []struct {
		name        string
		user        User
		expectedErr error
	}{
		{
			name: "valid user",
			user: User{
				Name:         "Administrador",
				Email:        "admin@email.com",
				PasswordHash: "hash",
			},
			expectedErr: nil,
		},
		{
			name: "empty name",
			user: User{
				Name:         "",
				Email:        "admin@email.com",
				PasswordHash: "hash",
			},
			expectedErr: ErrUserNameRequired,
		},
		{
			name: "empty email",
			user: User{
				Name:         "Administrador",
				Email:        "",
				PasswordHash: "hash",
			},
			expectedErr: ErrUserEmailRequired,
		},
		{
			name: "empty password",
			user: User{
				Name:         "Administrador",
				Email:        "admin@email.com",
				PasswordHash: "",
			},
			expectedErr: ErrPasswordRequired,
		},
	}

	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {

			err := tt.user.Validate()

			if !errors.Is(err, tt.expectedErr) {
				t.Fatalf("expected %v, got %v", tt.expectedErr, err)
			}
		})
	}
}

func TestUser_CheckPassword(t *testing.T) {

	hash, err := security.HashPassword("123456")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	user := User{
		PasswordHash: hash,
	}

	t.Run("valid password", func(t *testing.T) {

		if !user.CheckPassword("123456") {
			t.Fatal("expected password to be valid")
		}
	})

	t.Run("invalid password", func(t *testing.T) {

		if user.CheckPassword("654321") {
			t.Fatal("expected password to be invalid")
		}
	})
}

func TestUser_Update(t *testing.T) {

	t.Run("should update user", func(t *testing.T) {

		user := User{
			Name:         "Administrador",
			Email:        "admin@email.com",
			PasswordHash: "hash",
		}

		err := user.Update(
			"Novo Nome",
			"novo@email.com",
		)

		if err != nil {
			t.Fatalf("expected nil, got %v", err)
		}

		if user.Name != "Novo Nome" {
			t.Fatalf("expected Novo Nome, got %s", user.Name)
		}

		if user.Email != "novo@email.com" {
			t.Fatalf("expected novo@email.com, got %s", user.Email)
		}
	})

	t.Run("should return validation error", func(t *testing.T) {

		user := User{
			Name:         "Administrador",
			Email:        "admin@email.com",
			PasswordHash: "hash",
		}

		err := user.Update(
			"",
			"",
		)

		if !errors.Is(err, ErrUserNameRequired) {
			t.Fatalf("expected %v, got %v", ErrUserNameRequired, err)
		}
	})
}
