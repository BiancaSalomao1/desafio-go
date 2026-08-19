package order

import (
	"orders-api/internal/domain"
	"orders-api/internal/repository"
)

// HandleStockReservedUseCase trata a confirmação de reserva de estoque.
type HandleStockReservedUseCase struct {
	orderRepository repository.OrderRepository
}

// NewHandleStockReservedUseCase cria o caso de uso.
func NewHandleStockReservedUseCase(
	orderRepository repository.OrderRepository,
) *HandleStockReservedUseCase {
	return &HandleStockReservedUseCase{
		orderRepository: orderRepository,
	}
}

// Execute processa StockReserved.
// Por enquanto, a reserva de estoque não altera o status do pedido.
// O pedido continua PENDING até o pagamento.
func (uc *HandleStockReservedUseCase) Execute(
	orderID string,
) error {
	order, err := uc.orderRepository.FindByID(orderID)
	if err != nil {
		return err
	}

	if order.Status != domain.OrderStatusPending {
		return domain.ErrOrderStatusInvalid
	}

	return nil
}
