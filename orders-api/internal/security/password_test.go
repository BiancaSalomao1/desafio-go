package security

import "testing"

func TestHashPassword(t *testing.T) {

	password := "123456"

	hash, err := HashPassword(password)

	if err != nil {
		t.Fatalf("expected nil, got %v", err)
	}

	if hash == "" {
		t.Fatal("expected hash")
	}

	if hash == password {
		t.Fatal("password should not equal hash")
	}

	if !CheckPassword(hash, password) {
		t.Fatal("expected valid password")
	}
}

func TestCheckPassword(t *testing.T) {

	hash, err := HashPassword("123456")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	tests := []struct {
		name     string
		password string
		expected bool
	}{
		{
			name:     "valid password",
			password: "123456",
			expected: true,
		},
		{
			name:     "invalid password",
			password: "654321",
			expected: false,
		},
		{
			name:     "empty password",
			password: "",
			expected: false,
		},
	}

	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {

			result := CheckPassword(hash, tt.password)

			if result != tt.expected {
				t.Fatalf(
					"expected %v, got %v",
					tt.expected,
					result,
				)
			}
		})
	}
}
