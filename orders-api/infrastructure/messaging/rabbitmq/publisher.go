package rabbitmq

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"

	amqp "github.com/rabbitmq/amqp091-go"

	"orders-api/internal/messaging"
)

type Publisher struct {
	channel *amqp.Channel
	logger  *slog.Logger
}

func NewPublisher(
	connection *Connection,
	logger *slog.Logger,
) (*Publisher, error) {
	if connection == nil {
		return nil, fmt.Errorf("rabbitmq connection is required")
	}

	channel, err := connection.Channel()
	if err != nil {
		return nil, err
	}

	err = channel.ExchangeDeclare(
		ExchangeName,
		ExchangeType,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		channel.Close()

		return nil, fmt.Errorf("declare rabbitmq exchange: %w", err)
	}

	return &Publisher{
		channel: channel,
		logger:  logger,
	}, nil
}

func (p *Publisher) Publish(
	ctx context.Context,
	event messaging.Event,
) error {
	if p == nil || p.channel == nil {
		return fmt.Errorf("rabbitmq publisher is not initialized")
	}

	payload, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("marshal event: %w", err)
	}

	routingKey := routingKey(event.EventType)

	err = p.channel.PublishWithContext(
		ctx,
		ExchangeName,
		routingKey,
		false,
		false,
		amqp.Publishing{
			ContentType:   "application/json",
			DeliveryMode:  amqp.Persistent,
			MessageId:     event.MessageID,
			CorrelationId: event.CorrelationID,
			Body:          payload,
		},
	)
	if err != nil {
		if p.logger != nil {
			p.logger.ErrorContext(
				ctx,
				"rabbitmq publish failed",
				"exchange", ExchangeName,
				"routing_key", routingKey,
				"event_type", event.EventType,
				"message_id", event.MessageID,
				"correlation_id", event.CorrelationID,
				"saga_id", event.SagaID,
				"order_id", event.OrderID,
				"error", err,
			)
		}

		return fmt.Errorf("publish rabbitmq event: %w", err)
	}

	if p.logger != nil {
		p.logger.InfoContext(
			ctx,
			"rabbitmq event published",
			"exchange", ExchangeName,
			"routing_key", routingKey,
			"event_type", event.EventType,
			"message_id", event.MessageID,
			"correlation_id", event.CorrelationID,
			"saga_id", event.SagaID,
			"order_id", event.OrderID,
		)
	}

	return nil
}

func (p *Publisher) Close() error {
	if p == nil || p.channel == nil {
		return nil
	}

	return p.channel.Close()
}

func routingKey(eventType string) string {
	switch eventType {
	case "OrderCreated":
		return "order.created"
	case "ReserveStock":
		return "stock.reserve"
	case "StockReserved":
		return "stock.reserved"
	case "StockReservationFailed":
		return "stock.reservation_failed"
	case "ReleaseStock":
		return "stock.release"
	case "StockReleased":
		return "stock.released"
	default:
		return eventType
	}
}
