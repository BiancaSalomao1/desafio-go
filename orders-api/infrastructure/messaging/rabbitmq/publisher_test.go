package rabbitmq

import (
	"testing"

	"orders-api/infrastructure/logging"
)

func TestNewPublisher(t *testing.T) {
	connection, err := NewConnection("amqp://app:app@localhost:5672/")
	if err != nil {
		t.Fatalf("connect rabbitmq: %v", err)
	}

	defer connection.Close()

	logger := logging.New()

	publisher, err := NewPublisher(connection, logger)
	if err != nil {
		t.Fatalf("create publisher: %v", err)
	}

	defer publisher.Close()

	if publisher.channel == nil {
		t.Fatal("expected rabbitmq channel")
	}
}
