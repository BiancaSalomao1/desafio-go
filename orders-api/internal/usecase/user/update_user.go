package user

/*
struct UpdateUserUseCase

Responsabilidades:
- localizar um usuário;
- atualizar seus dados;
- persistir a alteração.

Campos:
- userRepository

Métodos:
- NewUpdateUserUseCase()
- Execute()
*/

import (
	"desafio-go/orders-api/internal/domain"
	"desafio-go/orders-api/internal/repository"
)

type UpdateUserUseCase struct {
	userRepository repository.UserRepository
}

func NewUpdateUserUseCase(
	userRepository repository.UserRepository,
) *UpdateUserUseCase {

	return &UpdateUserUseCase{
		userRepository: userRepository,
	}
}

func (uc *UpdateUserUseCase) Execute(
	user *domain.User,
) error {

	currentUser, err := uc.userRepository.FindByID(user.ID)
	if err != nil {
		return err
	}

	if err := currentUser.Update(
		user.Name,
		user.Email,
	); err != nil {
		return err
	}

	return uc.userRepository.Update(currentUser)
}
