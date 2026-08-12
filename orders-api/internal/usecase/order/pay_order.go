package order

/*
struct PayOrderUseCase

Responsabilidades:
- localizar pedido;
- alterar status para PAID;
- salvar alteração.

Métodos:
- NewPayOrderUseCase()
- Execute()
*/

import "desafio-go/orders-api/internal/repository"

type PayOrderUseCase struct {
	orderRepository repository.OrderRepository
}

func NewPayOrderUseCase(
	orderRepository repository.OrderRepository,
) *PayOrderUseCase {

	return &PayOrderUseCase{
		orderRepository: orderRepository,
	}
}

func (uc *PayOrderUseCase) Execute(id string) error {

	order, err := uc.orderRepository.FindByID(id)
	if err != nil {
		return err
	}

	if err := order.Pay(); err != nil {
		return err
	}

	return uc.orderRepository.Update(order)
}
