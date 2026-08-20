package auth

/*
struct LoginUseCase

Responsabilidades:
- autenticar usuário;
- gerar JWT;
- registrar token ativo.

Métodos:
- Execute()
*/

import (
	"context"
	"time"

	"orders-api/internal/domain"
	"orders-api/internal/repository"
	"orders-api/internal/security"
)

type LoginUseCase struct {
	userRepository repository.UserRepository
	tokenStore     security.TokenStore

	jwtSecret string
	jwtTTL    time.Duration
}

func NewLoginUseCase(
	userRepository repository.UserRepository,
	tokenStore security.TokenStore,
	jwtSecret string,
	jwtTTL time.Duration,
) *LoginUseCase {
	return &LoginUseCase{
		userRepository: userRepository,
		tokenStore:     tokenStore,
		jwtSecret:      jwtSecret,
		jwtTTL:         jwtTTL,
	}
}

func (uc *LoginUseCase) Execute(
	ctx context.Context,
	email string,
	password string,
) (string, error) {

	user, err := uc.userRepository.FindByEmail(email)
	if err != nil {
		return "", err
	}

	if !user.CheckPassword(password) {
		return "", domain.ErrInvalidCredentials
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

	if err := uc.tokenStore.Save(
		ctx,
		token,
		uc.jwtTTL,
	); err != nil {
		return "", err
	}

	return token, nil
}
