package rabbitmq

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"

	amqp "github.com/rabbitmq/amqp091-go"

	"product-service/internal/messaging"
)

const (
	ExchangeName = "saga.events"
	ExchangeType = "topic"
	QueueName    = "product-service.stock"
)

type Consumer struct {
	channel *amqp.Channel
	logger  *slog.Logger
}

func NewConsumer(
	connection *Connection,
	logger *slog.Logger,
) (*Consumer, error) {
	if connection == nil {
		return nil, fmt.Errorf("rabbitmq connection is required")
	}

	channel, err := connection.Channel()
	if err != nil {
		return nil, err
	}

	if err := channel.ExchangeDeclare(
		ExchangeName,
		ExchangeType,
		true,
		false,
		false,
		false,
		nil,
	); err != nil {
		_ = channel.Close()

		return nil, fmt.Errorf("declare rabbitmq exchange: %w", err)
	}

	_, err = channel.QueueDeclare(
		QueueName,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		_ = channel.Close()

		return nil, fmt.Errorf("declare rabbitmq queue: %w", err)
	}

	if err := channel.QueueBind(
		QueueName,
		"stock.reserve",
		ExchangeName,
		false,
		nil,
	); err != nil {
		_ = channel.Close()

		return nil, fmt.Errorf("bind stock.reserve: %w", err)
	}

	if err := channel.QueueBind(
		QueueName,
		"stock.release",
		ExchangeName,
		false,
		nil,
	); err != nil {
		_ = channel.Close()

		return nil, fmt.Errorf("bind stock.release: %w", err)
	}

	return &Consumer{
		channel: channel,
		logger:  logger,
	}, nil
}

func (c *Consumer) Consume(
	ctx context.Context,
	handler func(context.Context, messaging.Event) error,
) error {
	if c == nil || c.channel == nil {
		return fmt.Errorf("rabbitmq consumer is not initialized")
	}

	if handler == nil {
		return fmt.Errorf("rabbitmq consumer handler is required")
	}

	messages, err := c.channel.Consume(
		QueueName,
		"",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return fmt.Errorf("consume rabbitmq queue: %w", err)
	}

	for {
		select {
		case <-ctx.Done():
			return nil

		case delivery, ok := <-messages:
			if !ok {
				return fmt.Errorf("rabbitmq consumer channel closed")
			}

			var event messaging.Event

			if err := json.Unmarshal(delivery.Body, &event); err != nil {
				if c.logger != nil {
					c.logger.Error(
						"invalid rabbitmq event",
						"error", err,
					)
				}

				_ = delivery.Nack(false, false)
				continue
			}

			if err := handler(ctx, event); err != nil {
				if c.logger != nil {
					c.logger.ErrorContext(
						ctx,
						"rabbitmq event processing failed",
						"event_type", event.EventType,
						"message_id", event.MessageID,
						"correlation_id", event.CorrelationID,
						"saga_id", event.SagaID,
						"order_id", event.OrderID,
						"error", err,
					)
				}

				_ = delivery.Nack(false, true)
				continue
			}

			if err := delivery.Ack(false); err != nil && c.logger != nil {
				c.logger.ErrorContext(
					ctx,
					"rabbitmq message acknowledgement failed",
					"error", err,
				)
			}
		}
	}
}

func (c *Consumer) Close() error {
	if c == nil || c.channel == nil {
		return nil
	}

	return c.channel.Close()
}
