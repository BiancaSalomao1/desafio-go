package handler

/*
interface CreateUserUseCase

Responsabilidades:
- criar um usuário.

Métodos:
- Execute()
*/

import (
	"desafio-go/internal/domain"
)

type CreateUserUseCase interface {
	Execute(user *domain.User) error
}

/*
interface GetUserUseCase

Responsabilidades:
- buscar um usuário por ID.

Métodos:
- Execute()
*/

type GetUserUseCase interface {
	Execute(id string) (*domain.User, error)
}

/*
interface ListUsersUseCase

Responsabilidades:
- listar usuários.

Métodos:
- Execute()
*/

type ListUsersUseCase interface {
	Execute() ([]*domain.User, error)
}
