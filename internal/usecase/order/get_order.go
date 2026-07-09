package order

/*
struct GetOrderUseCase
- buscar pedido por ID.

Métodos:
- NewGetOrderUseCase()
- Execute()
*/

import (
	"desafio-go/internal/domain"
	"desafio-go/internal/repository"
)

type GetOrderUseCase struct {
	orderRepository repository.OrderRepository
}

func NewGetOrderUseCase(orderRepository repository.OrderRepository) *GetOrderUseCase {
	return &GetOrderUseCase{
		orderRepository: orderRepository,
	}
}

func (uc *GetOrderUseCase) Execute(id string) (*domain.Order, error) {
	return uc.orderRepository.FindByID(id)
}
