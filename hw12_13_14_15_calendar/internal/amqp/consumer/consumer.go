package simpleconsumer

import (
	"context"
	"fmt"

	"github.com/Fuchsoria/go_hw/hw12_13_14_15_calendar/internal/logger"
	"github.com/streadway/amqp"
)

type RMQConnection interface {
	Channel() (*amqp.Channel, error)
}

type Consumer struct {
	name   string
	conn   RMQConnection
	logger *logger.Logger
}

type AMQPMessage struct {
	EventID    string `json:"event_id"`
	EventTitle string `json:"event_title"`
	EventDate  int64  `json:"event_date"`
	User       string `json:"user"`
}

func New(name string, conn RMQConnection, logg *logger.Logger) *Consumer {
	return &Consumer{
		name:   name,
		conn:   conn,
		logger: logg,
	}
}

type Message struct {
	Ctx  context.Context
	Data []byte
}

func (c *Consumer) Consume(ctx context.Context, queue string) (<-chan Message, error) {
	messages := make(chan Message)

	ch, err := c.conn.Channel()
	if err != nil {
		return nil, fmt.Errorf("open channel: %w", err)
	}

	go func() {
		<-ctx.Done()
		if err := ch.Close(); err != nil {
			c.logger.Error(fmt.Errorf("cannot close, %w", err).Error())
		}
	}()

	deliveries, err := ch.Consume(queue, c.name, false, false, false, false, nil)
	if err != nil {
		return nil, fmt.Errorf("start consuming: %w", err)
	}

	go func() {
		defer func() {
			close(messages)
			c.logger.Info("close messages channel")
		}()

		for {
			select {
			case <-ctx.Done():
				return
			case del := <-deliveries:
				if err := del.Ack(false); err != nil {
					c.logger.Error(fmt.Errorf("cannot deliver message, %w", err).Error())
				}

				msg := Message{
					Ctx:  context.TODO(),
					Data: del.Body,
				}

				select {
				case <-ctx.Done():
					return
				case messages <- msg:
				}
			}
		}
	}()

	return messages, nil
}
