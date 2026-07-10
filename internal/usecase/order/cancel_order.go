package order

/*
struct CancelOrderUseCase

Responsabilidades:
- cancelar pedido;
- devolver estoque;
- atualizar pedido.

Métodos:
- NewCancelOrderUseCase()
- Execute()
*/

import "desafio-go/internal/repository"

type CancelOrderUseCase struct {
	orderRepository   repository.OrderRepository
	productRepository repository.ProductRepository
}

func NewCancelOrderUseCase(
	orderRepository repository.OrderRepository,
	productRepository repository.ProductRepository,
) *CancelOrderUseCase {

	return &CancelOrderUseCase{
		orderRepository:   orderRepository,
		productRepository: productRepository,
	}
}

func (uc *CancelOrderUseCase) Execute(id string) error {

	order, err := uc.orderRepository.FindByID(id)
	if err != nil {
		return err
	}

	if err := order.Cancel(); err != nil {
		return err
	}

	for _, item := range order.Items {

		product, err := uc.productRepository.FindByID(item.ProductID)
		if err != nil {
			return err
		}

		product.IncreaseStock(item.Quantity)

		if err := uc.productRepository.Update(product); err != nil {
			return err
		}
	}

	return uc.orderRepository.Update(order)
}
