package handler

/*
interface CreateOrderUseCase

Responsabilidades:
- criar um pedido.

Métodos:
- Execute()
*/

import (
	"desafio-go/orders-api/internal/domain"
)

type CreateOrderUseCase interface {
	Execute(order *domain.Order) error
}

/*
interface GetOrderUseCase

Responsabilidades:
- buscar um pedido por ID.

Métodos:
- Execute()
*/

type GetOrderUseCase interface {
	Execute(id string) (*domain.Order, error)
}

/*
interface PayOrderUseCase

Responsabilidades:
- realizar o pagamento de um pedido.

Métodos:
- Execute()
*/

type PayOrderUseCase interface {
	Execute(id string) error
}

/*
interface CancelOrderUseCase

Responsabilidades:
- cancelar um pedido.

Métodos:
- Execute()
*/

type CancelOrderUseCase interface {
	Execute(id string) error
}

type ListOrdersUseCase interface {
	Execute(limit int, offset int) ([]*domain.Order, error)
}
