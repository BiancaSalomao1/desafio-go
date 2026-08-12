package auth

import (
	"testing"
	"time"

	"desafio-go/orders-api/internal/domain"
	"desafio-go/orders-api/internal/security"
)

func BenchmarkLoginUseCase_Execute(b *testing.B) {

	hash, _ := security.HashPassword("123456")

	repository := &userRepositoryStub{
		user: domain.NewUser(
			"1",
			"Administrador",
			"admin@email.com",
			hash,
		),
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
