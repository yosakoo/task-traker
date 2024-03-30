package rabbitmq

import (
	"context"
	"fmt"
	"log"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Config struct {
	URL          string
	WaitTime     time.Duration
	Attempts     int
	Exchange     string
	ExchangeType string
	Queue        string 
}

type Connection struct {
	Config
	Connection *amqp.Connection
	Channel    *amqp.Channel
}

func New(cfg Config) (*Connection, error) {
	conn := &Connection{
		Config: cfg,
	}
	if err := conn.attemptConnect(); err != nil {
		return nil, fmt.Errorf("failed to connect to RabbitMQ: %w", err)
	}
	if err := conn.Channel.ExchangeDeclare(
		cfg.Exchange,     // name
		cfg.ExchangeType, // type
		true,             // durable
		false,            // auto-delete
		false,            // internal
		false,            // noWait
		nil,              // arguments
	); err != nil {
		return nil, fmt.Errorf("failed to exchange declare: %s", err)
	}

	if _, err := conn.Channel.QueueDeclare(
		cfg.Queue, // name
		true,      // Durable
		false,     // Auto-delete
		false,     // Exclusive
		false,     // No-wait
		nil,       // Arguments
	); err != nil {
		return nil, fmt.Errorf("failed to declare queue: %w", err)
	}

	if err := conn.Channel.QueueBind(cfg.Queue, "", cfg.Exchange, false, nil); err != nil {
		return nil, fmt.Errorf("failed to bind queue to exchange: %w", err)
	}

	return conn, nil
}

func (c *Connection) attemptConnect() error {
	var err error
	for i := c.Attempts; i > 0; i-- {
		if err = c.connect(); err == nil {
			break
		}

		log.Printf("RabbitMQ is trying to connect, attempts left: %d", i)
		time.Sleep(c.WaitTime)
	}

	if err != nil {
		return fmt.Errorf("attemptConnect: %w", err)
	}

	return nil
}

func (c *Connection) connect() error {
	var err error

	c.Connection, err = amqp.Dial(c.URL)
	if err != nil {
		return fmt.Errorf("amqp.Dial: %w", err)
	}

	c.Channel, err = c.Connection.Channel()
	if err != nil {
		return fmt.Errorf("Connection.Channel: %w", err)
	}

	return nil
}

func (c *Connection) PublishMessage(ctx context.Context, contentType string, body []byte) error {
	err := c.Channel.PublishWithContext(ctx,
		"",
		c.Config.Queue,
		false,
		false,
		amqp.Publishing{
			ContentType: contentType,
			Body:        body,
		},
	)
	if err != nil {
		return fmt.Errorf("failed to publish message: %w", err)
	}

	return nil
}

func (c *Connection) Close() error {
	if c.Connection == nil {
		return nil
	}

	if err := c.Channel.Close(); err != nil {
		return fmt.Errorf("failed to close channel: %w", err)
	}

	if err := c.Connection.Close(); err != nil {
		return fmt.Errorf("failed to close connection: %w", err)
	}

	return nil
}
