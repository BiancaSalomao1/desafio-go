package domain

import (
	"errors"
	"testing"

	"orders-api/internal/security"
)

func TestNewCustomer(t *testing.T) {

	customer := NewCustomer(
		"1",
		"João",
		"joao@email.com",
		"hash",
	)

	if customer == nil {
		t.Fatal("expected customer")
	}

	if customer.ID != "1" {
		t.Fatalf("expected id 1, got %s", customer.ID)
	}

	if customer.Name != "João" {
		t.Fatalf("expected João, got %s", customer.Name)
	}

	if customer.Email != "joao@email.com" {
		t.Fatalf("expected joao@email.com, got %s", customer.Email)
	}

	if customer.PasswordHash != "hash" {
		t.Fatalf("expected hash, got %s", customer.PasswordHash)
	}
}

func TestCustomer_Validate(t *testing.T) {

	tests := []struct {
		name        string
		customer    Customer
		expectedErr error
	}{
		{
			name: "valid customer",
			customer: Customer{
				Name:         "João",
				Email:        "joao@email.com",
				PasswordHash: "hash",
			},
			expectedErr: nil,
		},
		{
			name: "empty name",
			customer: Customer{
				Name:         "",
				Email:        "joao@email.com",
				PasswordHash: "hash",
			},
			expectedErr: ErrCustomerInvalid,
		},
		{
			name: "empty email",
			customer: Customer{
				Name:         "João",
				Email:        "",
				PasswordHash: "hash",
			},
			expectedErr: ErrCustomerInvalid,
		},
		{
			name: "empty password hash",
			customer: Customer{
				Name:         "João",
				Email:        "joao@email.com",
				PasswordHash: "",
			},
			expectedErr: ErrPasswordRequired,
		},
	}

	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {

			err := tt.customer.Validate()

			if !errors.Is(err, tt.expectedErr) {
				t.Fatalf(
					"expected %v, got %v",
					tt.expectedErr,
					err,
				)
			}
		})
	}
}

func TestCustomer_Update(t *testing.T) {

	t.Run("should update customer", func(t *testing.T) {

		customer := Customer{
			Name:         "João",
			Email:        "joao@email.com",
			PasswordHash: "hash",
		}

		err := customer.Update(
			"Maria",
			"maria@email.com",
		)

		if err != nil {
			t.Fatalf("expected nil, got %v", err)
		}

		if customer.Name != "Maria" {
			t.Fatalf(
				"expected Maria, got %s",
				customer.Name,
			)
		}

		if customer.Email != "maria@email.com" {
			t.Fatalf(
				"expected maria@email.com, got %s",
				customer.Email,
			)
		}
	})

	t.Run("should return validation error", func(t *testing.T) {

		customer := Customer{
			Name:         "João",
			Email:        "joao@email.com",
			PasswordHash: "hash",
		}

		err := customer.Update(
			"",
			"",
		)

		if !errors.Is(err, ErrCustomerInvalid) {
			t.Fatalf(
				"expected %v, got %v",
				ErrCustomerInvalid,
				err,
			)
		}
	})
}

func TestCustomer_CheckPassword(t *testing.T) {

	t.Run("should return true for correct password", func(t *testing.T) {

		password := "mysecretpassword"

		passwordHash, err := security.HashPassword(password)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		customer := Customer{
			PasswordHash: passwordHash,
		}

		if !customer.CheckPassword(password) {
			t.Fatal("expected true, got false")
		}
	})

	t.Run("should return false for incorrect password", func(t *testing.T) {

		password := "mysecretpassword"

		passwordHash, err := security.HashPassword(password)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		customer := Customer{
			PasswordHash: passwordHash,
		}

		if customer.CheckPassword("wrongpassword") {
			t.Fatal("expected false, got true")
		}
	})
}
