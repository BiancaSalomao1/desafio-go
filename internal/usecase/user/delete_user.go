package user

/*
struct DeleteUserUseCase

Responsabilidades:
- remover um usuário.

Campos:
- userRepository

Métodos:
- NewDeleteUserUseCase()
- Execute()
*/

import (
	"strings"

	"desafio-go/internal/domain"
	"desafio-go/internal/repository"
)

type DeleteUserUseCase struct {
	userRepository repository.UserRepository
}

func NewDeleteUserUseCase(
	userRepository repository.UserRepository,
) *DeleteUserUseCase {

	return &DeleteUserUseCase{
		userRepository: userRepository,
	}
}

func (uc *DeleteUserUseCase) Execute(
	id string,
) error {

	if _, err := uc.userRepository.FindByID(id); err != nil {
		return err
	}

	err := uc.userRepository.Delete(id)
	if err != nil {

		if strings.Contains(
			err.Error(),
			"SQLSTATE 23503",
		) {
			return domain.ErrUserInUse
		}

		return err
	}

	return nil
}
