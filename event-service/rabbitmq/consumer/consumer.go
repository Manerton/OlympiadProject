package consumer

import (
	"context"
	"encoding/json"
	"log/slog"
	"main/internal/dto/rabbit_dto"
	"main/internal/lib/liblogger"

	amqp "github.com/rabbitmq/amqp091-go"
)

// command
const (
	GET    = ""
	CREATE = "create"
	UPDATE = "update"
	DELETE = "delete"
)

type EventService interface {
}

type RabbitConsumer struct {
	log           *slog.Logger
	rabbitChannel *amqp.Channel
	eventService  EventService
}

func New(log *slog.Logger, channel *amqp.Channel, eventService EventService) *RabbitConsumer {
	return &RabbitConsumer{
		log:           log,
		rabbitChannel: channel,
		eventService:  eventService,
	}
}

func (c *RabbitConsumer) Start(ctx context.Context, queueName string) {

	msgs, err := c.rabbitChannel.Consume(
		queueName,
		"",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		c.log.Error("Failed consume queue", slog.String("queue name", queueName), liblogger.Err(err))
		return
	}

	go func() {
		for {
			select {
			case <-ctx.Done():
				c.log.Info("init cancel consumer")
				if err := c.rabbitChannel.Cancel("", false); err != nil {
					c.log.Error("failed to cancel consumer", liblogger.Err(err))
				}
				return
			case msg, ok := <-msgs:
				if !ok {
					c.log.Info("message channel closed")
					return
				}

				rabbitDTO := rabbit_dto.RabbitDTO{}

				if err := json.Unmarshal(msg.Body, &rabbitDTO); err != nil {
					c.log.Error("invalid message format", liblogger.Err(err))
					msg.Nack(false, false)
					continue
				}

				err = c.handler(ctx, rabbitDTO)
				if err != nil {
					c.log.Error("task failed", liblogger.Err(err))
					msg.Nack(false, false)
					continue
				}

				msg.Ack(false)
			}
		}
	}()
}

func (c *RabbitConsumer) handler(ctx context.Context, rabbitDTO rabbit_dto.RabbitDTO) error {

	return nil

}
