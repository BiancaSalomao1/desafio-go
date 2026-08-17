package order

/*
struct GetOrderUseCase

Responsabilidades:
- buscar pedido por ID.

Métodos:
- NewGetOrderUseCase()
- Execute()
*/

import (
	"orders-api/internal/domain"
	"orders-api/internal/repository"
)

type GetOrderUseCase struct {
	orderRepository repository.OrderRepository
}

func NewGetOrderUseCase(
	orderRepository repository.OrderRepository,
) *GetOrderUseCase {

	return &GetOrderUseCase{
		orderRepository: orderRepository,
	}
}

func (uc *GetOrderUseCase) Execute(id string) (*domain.Order, error) {

	return uc.orderRepository.FindByID(id)
}
