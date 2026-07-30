package domain

import (
	"errors"
	"testing"
)

func TestNewCustomer(t *testing.T) {

	customer := NewCustomer(
		"1",
		"João",
		"joao@email.com",
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
				Name:  "João",
				Email: "joao@email.com",
			},
			expectedErr: nil,
		},
		{
			name: "empty name",
			customer: Customer{
				Name:  "",
				Email: "joao@email.com",
			},
			expectedErr: ErrCustomerInvalid,
		},
		{
			name: "empty email",
			customer: Customer{
				Name:  "João",
				Email: "",
			},
			expectedErr: ErrCustomerInvalid,
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
			Name:  "João",
			Email: "joao@email.com",
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
			Name:  "João",
			Email: "joao@email.com",
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
