package handler

/*
interface LoginUseCase

Responsabilidades:
- autenticar usuário.

Métodos:
- Execute()
*/

type LoginUseCase interface {
	Execute(
		email string,
		password string,
	) (string, error)
}
