package rabbitmq

import (
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Connection struct {
	conn *amqp.Connection
}

func NewConnection(url string) (*Connection, error) {
	if url == "" {
		return nil, fmt.Errorf("rabbitmq url is required")
	}

	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, fmt.Errorf("connect to rabbitmq: %w", err)
	}

	return &Connection{
		conn: conn,
	}, nil
}

func (c *Connection) Channel() (*amqp.Channel, error) {
	if c == nil || c.conn == nil {
		return nil, fmt.Errorf("rabbitmq connection is not initialized")
	}

	channel, err := c.conn.Channel()
	if err != nil {
		return nil, fmt.Errorf("create rabbitmq channel: %w", err)
	}

	return channel, nil
}

func (c *Connection) Close() error {
	if c == nil || c.conn == nil {
		return nil
	}

	return c.conn.Close()
}
