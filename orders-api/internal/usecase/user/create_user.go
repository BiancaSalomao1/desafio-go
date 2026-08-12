package user

/*
struct CreateUserUseCase

Responsabilidades:
- criar usuário;
- validar dados;
- gerar hash da senha;
- salvar usuário.

Métodos:
- NewCreateUserUseCase()
- Execute()
*/

import (
	"desafio-go/orders-api/internal/domain"
	"desafio-go/orders-api/internal/repository"
	"desafio-go/orders-api/internal/security"
)

type CreateUserUseCase struct {
	userRepository repository.UserRepository
}

func NewCreateUserUseCase(
	userRepository repository.UserRepository,
) *CreateUserUseCase {

	return &CreateUserUseCase{
		userRepository: userRepository,
	}
}

func (uc *CreateUserUseCase) Execute(
	user *domain.User,
) error {

	if err := user.Validate(); err != nil {
		return err
	}

	passwordHash, err := security.HashPassword(
		user.PasswordHash,
	)
	if err != nil {
		return err
	}

	user.PasswordHash = passwordHash

	return uc.userRepository.Save(user)
}
