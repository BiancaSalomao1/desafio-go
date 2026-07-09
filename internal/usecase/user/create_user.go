package user

/*
struct CreateUserUseCase
- criar usuário;
- validar dados;
- salvar usuário.

Métodos:
- NewCreateUserUseCase()
- Execute()
*/

import (
	"desafio-go/internal/domain"
	"desafio-go/internal/repository"
)

type CreateUserUseCase struct {
	userRepository repository.UserRepository
}

func NewCreateUserUseCase(userRepository repository.UserRepository) *CreateUserUseCase {
	return &CreateUserUseCase{
		userRepository: userRepository,
	}
}

func (uc *CreateUserUseCase) Execute(user *domain.User) error {

	if err := user.Validate(); err != nil {
		return err
	}

	return uc.userRepository.Save(user)
}
