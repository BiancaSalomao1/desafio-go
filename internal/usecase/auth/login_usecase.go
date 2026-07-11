package auth

/*
struct LoginUseCase

Responsabilidades:
- autenticar usuário;
- gerar JWT.

Métodos:
- Execute()
*/

import (
	"errors"
	"time"

	"desafio-go/internal/security"

	"desafio-go/internal/repository"
)

type LoginUseCase struct {
	userRepository repository.UserRepository

	jwtSecret string
	jwtTTL    time.Duration
}

func NewLoginUseCase(
	userRepository repository.UserRepository,
	jwtSecret string,
	jwtTTL time.Duration,
) *LoginUseCase {

	return &LoginUseCase{
		userRepository: userRepository,
		jwtSecret:      jwtSecret,
		jwtTTL:         jwtTTL,
	}
}

func (uc *LoginUseCase) Execute(
	email string,
	password string,
) (string, error) {

	user, err := uc.userRepository.FindByEmail(email)
	if err != nil {
		return "", err
	}

	if !user.CheckPassword(password) {
		return "", errors.New("invalid credentials")
	}

	token, err := security.GenerateToken(
		user.ID,
		user.Email,
		uc.jwtSecret,
		uc.jwtTTL,
	)

	if err != nil {
		return "", err
	}

	return token, nil
}
