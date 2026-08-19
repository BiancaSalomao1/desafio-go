package rabbitmq

import (
	"context"
	"log/slog"

	"orders-api/internal/messaging"
	orderusecase "orders-api/internal/usecase/order"
)

type OrderEventHandler struct {
	stockReservedUseCase          *orderusecase.HandleStockReservedUseCase
	stockReservationFailedUseCase *orderusecase.HandleStockReservationFailedUseCase
	logger                        *slog.Logger
}

func NewOrderEventHandler(
	stockReservedUseCase *orderusecase.HandleStockReservedUseCase,
	stockReservationFailedUseCase *orderusecase.HandleStockReservationFailedUseCase,
	logger *slog.Logger,
) *OrderEventHandler {
	return &OrderEventHandler{
		stockReservedUseCase:          stockReservedUseCase,
		stockReservationFailedUseCase: stockReservationFailedUseCase,
		logger:                        logger,
	}
}

func (h *OrderEventHandler) Handle(
	ctx context.Context,
	event messaging.Event,
) error {
	switch event.EventType {
	case "StockReserved":
		if h.logger != nil {
			h.logger.InfoContext(
				ctx,
				"stock reservation confirmed",
				"order_id", event.OrderID,
				"saga_id", event.SagaID,
				"correlation_id", event.CorrelationID,
			)
		}

		return h.stockReservedUseCase.Execute(event.OrderID)

	case "StockReservationFailed":
		if h.logger != nil {
			h.logger.WarnContext(
				ctx,
				"stock reservation failed",
				"order_id", event.OrderID,
				"saga_id", event.SagaID,
				"correlation_id", event.CorrelationID,
			)
		}

		return h.stockReservationFailedUseCase.Execute(event.OrderID)

	default:
		return nil
	}
}
