package auth

import (
	"context"
	"testing"
	"time"

	"orders-api/internal/domain"
	"orders-api/internal/security"
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

	tokenStore := &tokenStoreStub{}

	useCase := NewLoginUseCase(
		repository,
		tokenStore,
		"secret",
		time.Hour,
	)

	ctx := context.Background()

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, _ = useCase.Execute(
			ctx,
			"admin@email.com",
			"123456",
		)
	}
}
