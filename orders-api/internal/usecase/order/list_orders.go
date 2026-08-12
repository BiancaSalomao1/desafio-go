package order

/*
struct ListOrdersUseCase
- listar pedidos.

Métodos:
- NewListOrdersUseCase()
- Execute()
*/

import (
	"desafio-go/orders-api/internal/domain"
	"desafio-go/orders-api/internal/repository"
)

type ListOrdersUseCase struct {
	orderRepository repository.OrderRepository
}

func NewListOrdersUseCase(orderRepository repository.OrderRepository) *ListOrdersUseCase {
	return &ListOrdersUseCase{
		orderRepository: orderRepository,
	}
}

func (uc *ListOrdersUseCase) Execute(limit, offset int) ([]*domain.Order, error) {
	return uc.orderRepository.List(limit, offset)
}
