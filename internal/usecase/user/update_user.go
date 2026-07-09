package user

/*
struct UpdateUserUseCase
- atualizar usuário.

Métodos:
- NewUpdateUserUseCase()
- Execute()
*/

import (
	"desafio-go/internal/domain"
	"desafio-go/internal/repository"
)

type UpdateUserUseCase struct {
	userRepository repository.UserRepository
}

func NewUpdateUserUseCase(userRepository repository.UserRepository) *UpdateUserUseCase {
	return &UpdateUserUseCase{
		userRepository: userRepository,
	}
}

func (uc *UpdateUserUseCase) Execute(user *domain.User) error {

	if err := user.Validate(); err != nil {
		return err
	}

	return uc.userRepository.Update(user)
}
