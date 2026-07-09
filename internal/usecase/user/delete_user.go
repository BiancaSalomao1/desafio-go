package user

/*
struct DeleteUserUseCase
- excluir usuário.

Métodos:
- NewDeleteUserUseCase()
- Execute()
*/

import "desafio-go/internal/repository"

type DeleteUserUseCase struct {
	userRepository repository.UserRepository
}

func NewDeleteUserUseCase(userRepository repository.UserRepository) *DeleteUserUseCase {
	return &DeleteUserUseCase{
		userRepository: userRepository,
	}
}

func (uc *DeleteUserUseCase) Execute(id string) error {
	return uc.userRepository.Delete(id)
}
