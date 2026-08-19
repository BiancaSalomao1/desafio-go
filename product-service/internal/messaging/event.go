package messaging

import "time"

type Event struct {
	MessageID     string    `json:"message_id"`
	EventType     string    `json:"event_type"`
	CorrelationID string    `json:"correlation_id"`
	SagaID        string    `json:"saga_id"`
	OrderID       string    `json:"order_id"`
	OccurredAt    time.Time `json:"occurred_at"`
	Data          any       `json:"data,omitempty"`
}
