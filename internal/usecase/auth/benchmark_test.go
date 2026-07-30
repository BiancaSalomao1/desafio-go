package auth

/*
Benchmark do LoginUseCase.

Responsabilidades:
- medir desempenho da autenticação.
*/

import (
	"testing"
	"time"

	"desafio-go/internal/domain"
	"desafio-go/internal/security"
)

func BenchmarkLoginUseCase_Execute(b *testing.B) {

	hash, _ := security.HashPassword("123456")

	repository := &userRepositoryStub{
		user: &domain.User{
			ID:           "1",
			Name:         "Administrador",
			Email:        "admin@email.com",
			PasswordHash: hash,
		},
	}

	useCase := NewLoginUseCase(
		repository,
		"secret",
		time.Hour,
	)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {

		_, _ = useCase.Execute(
			"admin@email.com",
			"123456",
		)
	}
}
