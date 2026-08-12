package user

/*
struct GetUserUseCase
- buscar usuário por ID.

Métodos:
- NewGetUserUseCase()
- Execute()
*/

import (
	"desafio-go/orders-api/internal/domain"
	"desafio-go/orders-api/internal/repository"
)

type GetUserUseCase struct {
	userRepository repository.UserRepository
}

func NewGetUserUseCase(userRepository repository.UserRepository) *GetUserUseCase {
	return &GetUserUseCase{
		userRepository: userRepository,
	}
}

func (uc *GetUserUseCase) Execute(id string) (*domain.User, error) {
	return uc.userRepository.FindByID(id)
}
