package handler

import (
	"desafio-go/orders-api/internal/domain"
)

type CreateUserUseCase interface {
	Execute(user *domain.User) error
}

type GetUserUseCase interface {
	Execute(id string) (*domain.User, error)
}

type ListUsersUseCase interface {
	Execute() ([]*domain.User, error)
}

type UpdateUserUseCase interface {
	Execute(user *domain.User) error
}

type DeleteUserUseCase interface {
	Execute(id string) error
}
