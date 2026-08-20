package handler

import "context"

/*
interface LoginUseCase

Responsabilidades:
- autenticar usuário.

Métodos:
- Execute()
*/

type LoginUseCase interface {
	Execute(
		ctx context.Context,
		email string,
		password string,
	) (string, error)
}
