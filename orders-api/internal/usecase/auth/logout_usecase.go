package auth

import (
	"context"

	"orders-api/internal/security"
)

type LogoutUseCase struct {
	tokenStore security.TokenStore
}

func NewLogoutUseCase(
	tokenStore security.TokenStore,
) *LogoutUseCase {
	return &LogoutUseCase{
		tokenStore: tokenStore,
	}
}

func (uc *LogoutUseCase) Execute(
	ctx context.Context,
	token string,
) error {
	return uc.tokenStore.Delete(
		ctx,
		token,
	)
}
