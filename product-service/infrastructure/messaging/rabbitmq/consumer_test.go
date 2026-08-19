package rabbitmq

import (
	"context"
	"testing"
	"time"

	"product-service/internal/messaging"
)

func TestNewConsumer(t *testing.T) {
	connection, err := NewConnection(
		"amqp://app:app@localhost:5672/",
	)
	if err != nil {
		t.Fatalf("connect rabbitmq: %v", err)
	}

	defer connection.Close()

	consumer, err := NewConsumer(connection, nil)
	if err != nil {
		t.Fatalf("create consumer: %v", err)
	}

	defer consumer.Close()

	if consumer.channel == nil {
		t.Fatal("expected rabbitmq channel")
	}
}

func TestConsumerContextCancellation(t *testing.T) {
	connection, err := NewConnection(
		"amqp://app:app@localhost:5672/",
	)
	if err != nil {
		t.Fatalf("connect rabbitmq: %v", err)
	}

	defer connection.Close()

	consumer, err := NewConsumer(connection, nil)
	if err != nil {
		t.Fatalf("create consumer: %v", err)
	}

	defer consumer.Close()

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	done := make(chan error, 1)

	go func() {
		done <- consumer.Consume(
			ctx,
			func(
				context.Context,
				messaging.Event,
			) error {
				return nil
			},
		)
	}()

	select {
	case err := <-done:
		if err != nil {
			t.Fatalf("expected nil error, got %v", err)
		}

	case <-time.After(time.Second):
		t.Fatal("consumer did not stop after context cancellation")
	}
}
