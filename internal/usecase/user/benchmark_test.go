package user

/*
Benchmark do CreateUserUseCase.

Responsabilidades:
- medir desempenho da criação de usuários.
*/

import (
	"testing"

	"desafio-go/internal/domain"
)

func BenchmarkCreateUserUseCase_Execute(b *testing.B) {

	repository := &userRepositoryStub{}

	useCase := NewCreateUserUseCase(repository)

	user := &domain.User{
		Name:         "Administrador",
		Email:        "admin@email.com",
		PasswordHash: "123456",
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = useCase.Execute(user)
	}
}
func BenchmarkUpdateUserUseCase_Execute(b *testing.B) {

	repository := &userRepositoryStub{
		user: &domain.User{
			ID:           "1",
			Name:         "Administrador",
			Email:        "admin@email.com",
			PasswordHash: "hash",
		},
	}

	useCase := NewUpdateUserUseCase(repository)

	user := &domain.User{
		ID:           "1",
		Name:         "Novo Nome",
		Email:        "novo@email.com",
		PasswordHash: "hash",
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = useCase.Execute(user)
	}
}
