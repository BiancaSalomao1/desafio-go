package messaging

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"time"

	"product-service/internal/usecase/product"

	"github.com/google/uuid"

	events "product-service/internal/messaging/events"
)

type EventPublisher interface {
	Publish(context.Context, Event) error
}

type StockHandler struct {
	logger              *slog.Logger
	reserveStockUseCase *product.ReserveStockUseCase
	releaseStockUseCase *product.ReleaseStockUseCase
	publisher           EventPublisher
}

func NewStockHandler(
	logger *slog.Logger,
	reserveStockUseCase *product.ReserveStockUseCase,
	releaseStockUseCase *product.ReleaseStockUseCase,
	publisher EventPublisher,
) *StockHandler {
	return &StockHandler{
		logger:              logger,
		reserveStockUseCase: reserveStockUseCase,
		releaseStockUseCase: releaseStockUseCase,
		publisher:           publisher,
	}
}

func (h *StockHandler) Handle(
	ctx context.Context,
	event Event,
) error {
	switch event.EventType {
	case "ReserveStock":
		return h.handleReserveStock(ctx, event)

	case "ReleaseStock":
		return h.handleReleaseStock(ctx, event)

	default:
		if h.logger != nil {
			h.logger.InfoContext(
				ctx,
				"rabbitmq event ignored",
				"event_type", event.EventType,
			)
		}

		return nil
	}
}

func (h *StockHandler) handleReserveStock(
	ctx context.Context,
	event Event,
) error {
	var data events.ReserveStockData

	if err := decodeEventData(event.Data, &data); err != nil {
		return fmt.Errorf("decode reserve stock event: %w", err)
	}

	if err := h.reserveStockUseCase.Execute(
		mapReserveStockItems(data.Items),
	); err != nil {
		if h.logger != nil {
			h.logger.ErrorContext(
				ctx,
				"stock reservation failed",
				"order_id", event.OrderID,
				"saga_id", event.SagaID,
				"correlation_id", event.CorrelationID,
				"error", err,
			)
		}

		return h.publishReservationFailed(
			ctx,
			event,
			err,
		)
	}

	if h.logger != nil {
		h.logger.InfoContext(
			ctx,
			"stock reserved",
			"order_id", event.OrderID,
			"saga_id", event.SagaID,
			"correlation_id", event.CorrelationID,
		)
	}

	return h.publisher.Publish(
		ctx,
		Event{
			MessageID:     uuid.NewString(),
			EventType:     "StockReserved",
			CorrelationID: event.CorrelationID,
			SagaID:        event.SagaID,
			OrderID:       event.OrderID,
			OccurredAt:    time.Now().UTC(),
			Data: events.StockReservedData{
				OrderID: event.OrderID,
			},
		},
	)
}

func (h *StockHandler) handleReleaseStock(
	ctx context.Context,
	event Event,
) error {
	var data events.ReleaseStockData

	if err := decodeEventData(event.Data, &data); err != nil {
		return fmt.Errorf("decode release stock event: %w", err)
	}

	if err := h.releaseStockUseCase.Execute(
		mapReleaseStockItems(data.Items),
	); err != nil {
		return fmt.Errorf("release stock: %w", err)
	}

	if h.logger != nil {
		h.logger.InfoContext(
			ctx,
			"stock released",
			"order_id", event.OrderID,
			"saga_id", event.SagaID,
			"correlation_id", event.CorrelationID,
		)
	}

	return nil
}

func (h *StockHandler) publishReservationFailed(
	ctx context.Context,
	event Event,
	err error,
) error {
	return h.publisher.Publish(
		ctx,
		Event{
			MessageID:     uuid.NewString(),
			EventType:     "StockReservationFailed",
			CorrelationID: event.CorrelationID,
			SagaID:        event.SagaID,
			OrderID:       event.OrderID,
			OccurredAt:    time.Now().UTC(),
			Data: events.StockReservationFailedData{
				OrderID: event.OrderID,
				Reason:  err.Error(),
			},
		},
	)
}

func decodeEventData(
	data any,
	target any,
) error {
	payload, err := json.Marshal(data)
	if err != nil {
		return err
	}

	return json.Unmarshal(payload, target)
}

func mapReserveStockItems(
	items []events.StockItem,
) []product.StockItem {
	result := make([]product.StockItem, 0, len(items))

	for _, item := range items {
		result = append(result, product.StockItem{
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
		})
	}

	return result
}

func mapReleaseStockItems(
	items []events.StockItem,
) []product.StockItem {
	result := make([]product.StockItem, 0, len(items))

	for _, item := range items {
		result = append(result, product.StockItem{
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
		})
	}

	return result
}
