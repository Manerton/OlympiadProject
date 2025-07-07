package consumer

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"main/internal/dto/rabbit_dto"
	"main/internal/lib/liblogger"

	"github.com/google/uuid"
	amqp "github.com/rabbitmq/amqp091-go"
)

const (
	EVENT_TABLE = "event"
)

// command
const (
	GET    = ""
	CREATE = "create"
	UPDATE = "update"
	DELETE = "delete"
)

type EventService interface {
	Create(ctx context.Context) (uuid.UUID, error)
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

	id := ""
	if rabbitDTO.Method == UPDATE || rabbitDTO.Method == DELETE {
		var ok bool
		id, ok = rabbitDTO.Data.SearchAttributes["id"].(string)
		if !ok {
			c.log.Error("failed id does not exist", slog.String("id", id))
			return fmt.Errorf("failed update: ID does not exist")
		}
		c.log.Debug("id", slog.String("id", id))
	}

	switch rabbitDTO.Method {
	case CREATE:
		err := c.create(ctx, rabbitDTO.Data.Table, rabbitDTO.Data.Attributes)
		if err != nil {
			c.log.Error("failed create data", liblogger.Err(err))
			return fmt.Errorf("failed create data")
		}
	case UPDATE:
		err := c.update(ctx, rabbitDTO.Data.Table, rabbitDTO.Data.Attributes, id)
		if err != nil {
			c.log.Error("failed update data", liblogger.Err(err))
			return fmt.Errorf("failed upadate data")
		}

	case DELETE:
		err := c.delete(ctx, rabbitDTO.Data.Table, id)
		if err != nil {
			c.log.Error("failed delete data", liblogger.Err(err))
			return fmt.Errorf("failed delete data")
		}
	}
	return nil
}

func (c *RabbitConsumer) create(ctx context.Context, tableName string, data map[string]any) error {

	switch tableName {
	case EVENT_TABLE:

		_, err := c.eventService.Create(ctx)
		if err != nil {
			return fmt.Errorf("failed create event: %w", err)
		}
	default:
		return fmt.Errorf("unexpected table: %s", tableName)
	}

	return nil
}

func (c *RabbitConsumer) update(ctx context.Context, tableName string, data map[string]any, id string) error {
	switch tableName {
	case EVENT_TABLE:

	default:
		return fmt.Errorf("unexpected table: %s", tableName)
	}

	return nil
}

func (c *RabbitConsumer) delete(ctx context.Context, tableName string, id string) error {
	switch tableName {
	case EVENT_TABLE:

	default:
		return fmt.Errorf("unexpected table: %s", tableName)
	}

	return nil
}
