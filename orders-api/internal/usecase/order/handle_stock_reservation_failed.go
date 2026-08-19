package order

import (
	"orders-api/internal/repository"
)

type HandleStockReservationFailedUseCase struct {
	orderRepository repository.OrderRepository
}

func NewHandleStockReservationFailedUseCase(
	orderRepository repository.OrderRepository,
) *HandleStockReservationFailedUseCase {
	return &HandleStockReservationFailedUseCase{
		orderRepository: orderRepository,
	}
}

func (uc *HandleStockReservationFailedUseCase) Execute(
	orderID string,
) error {
	order, err := uc.orderRepository.FindByID(orderID)
	if err != nil {
		return err
	}

	if err := order.Cancel(); err != nil {
		return err
	}

	return uc.orderRepository.Update(order)
}
