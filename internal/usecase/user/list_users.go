package user

/*
struct ListUsersUseCase
- listar usuários.

Métodos:
- NewListUsersUseCase()
- Execute()
*/

import (
	"desafio-go/internal/domain"
	"desafio-go/internal/repository"
)

type ListUsersUseCase struct {
	userRepository repository.UserRepository
}

func NewListUsersUseCase(userRepository repository.UserRepository) *ListUsersUseCase {
	return &ListUsersUseCase{
		userRepository: userRepository,
	}
}

func (uc *ListUsersUseCase) Execute() ([]*domain.User, error) {
	return uc.userRepository.FindAll()
}
