package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"

	"github.com/twmb/franz-go/pkg/kgo"

	"orders-api/internal/messaging"
)

type Producer struct {
	client *kgo.Client
	topic  string
	logger *slog.Logger
}

func NewProducer(
	brokers []string,
	topic string,
	logger *slog.Logger,
) (*Producer, error) {
	if len(brokers) == 0 {
		return nil, fmt.Errorf("kafka brokers are required")
	}

	if topic == "" {
		return nil, fmt.Errorf("kafka topic is required")
	}

	client, err := kgo.NewClient(
		kgo.SeedBrokers(brokers...),
	)
	if err != nil {
		return nil, fmt.Errorf("create kafka client: %w", err)
	}

	return &Producer{
		client: client,
		topic:  topic,
		logger: logger,
	}, nil
}

func (p *Producer) Publish(
	ctx context.Context,
	event messaging.Event,
) error {
	payload, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("marshal event: %w", err)
	}

	record := &kgo.Record{
		Topic: p.topic,
		Key:   []byte(event.OrderID),
		Value: payload,
	}

	if err := p.client.ProduceSync(ctx, record).FirstErr(); err != nil {
		if p.logger != nil {
			p.logger.ErrorContext(
				ctx,
				"kafka publish failed",
				"topic", p.topic,
				"event_type", event.EventType,
				"message_id", event.MessageID,
				"correlation_id", event.CorrelationID,
				"saga_id", event.SagaID,
				"error", err,
			)
		}

		return fmt.Errorf("publish kafka event: %w", err)
	}

	if p.logger != nil {
		p.logger.InfoContext(
			ctx,
			"kafka event published",
			"topic", p.topic,
			"event_type", event.EventType,
			"message_id", event.MessageID,
			"correlation_id", event.CorrelationID,
			"saga_id", event.SagaID,
		)
	}

	return nil
}

func (p *Producer) Close() {
	p.client.Close()
}
