package security

import (
	"testing"
	"time"
)

func TestGenerateToken(t *testing.T) {

	token, err := GenerateToken(
		"1",
		"admin@email.com",
		"secret",
		time.Hour,
	)

	if err != nil {
		t.Fatalf("expected nil, got %v", err)
	}

	if token == "" {
		t.Fatal("expected token")
	}
}

func TestValidateToken(t *testing.T) {

	t.Run("valid token", func(t *testing.T) {

		token, _ := GenerateToken(
			"1",
			"admin@email.com",
			"secret",
			time.Hour,
		)

		claims, err := ValidateToken(
			token,
			"secret",
		)

		if err != nil {
			t.Fatalf("expected nil, got %v", err)
		}

		if claims["sub"] != "1" {
			t.Fatalf("expected user id 1")
		}

		if claims["email"] != "admin@email.com" {
			t.Fatalf("expected email")
		}
	})

	t.Run("invalid secret", func(t *testing.T) {

		token, _ := GenerateToken(
			"1",
			"admin@email.com",
			"secret",
			time.Hour,
		)

		_, err := ValidateToken(
			token,
			"wrong-secret",
		)

		if err == nil {
			t.Fatal("expected error")
		}
	})

	t.Run("invalid token", func(t *testing.T) {

		_, err := ValidateToken(
			"invalid.token",
			"secret",
		)

		if err == nil {
			t.Fatal("expected error")
		}
	})

	t.Run("expired token", func(t *testing.T) {

		token, _ := GenerateToken(
			"1",
			"admin@email.com",
			"secret",
			-time.Minute,
		)

		_, err := ValidateToken(
			token,
			"secret",
		)

		if err == nil {
			t.Fatal("expected expiration error")
		}
	})
}
