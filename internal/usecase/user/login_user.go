package user

/*
struct LoginUserUseCase
- autenticar usuário.

Métodos:
- NewLoginUserUseCase()
- Execute()
*/

import (
	"desafio-go/internal/domain"
	"desafio-go/internal/repository"
)

type LoginUserUseCase struct {
	userRepository repository.UserRepository
}

func NewLoginUserUseCase(userRepository repository.UserRepository) *LoginUserUseCase {
	return &LoginUserUseCase{
		userRepository: userRepository,
	}
}

func (uc *LoginUserUseCase) Execute(email string) (*domain.User, error) {
	return uc.userRepository.FindByEmail(email)
}
