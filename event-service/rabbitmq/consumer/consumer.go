package consumer

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"main/internal/dto/event_dto"
	"main/internal/dto/rabbit_dto"
	"main/internal/lib/liblogger"
	"main/rabbitmq"
	"time"

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
	CreateEvent(ctx context.Context, eventDTO event_dto.CreateEventDTORequest) (uuid.UUID, error)
	UpdateEvent(ctx context.Context, id string, eventDTO event_dto.UpdateEventDTORequest) error
	DeleteEvent(ctx context.Context, id string) error
}

type RabbitConsumer struct {
	log          *slog.Logger
	eventService EventService

	connectionManager *rabbitmq.ConnectionManager
}

func New(log *slog.Logger, mananger *rabbitmq.ConnectionManager, eventService EventService) *RabbitConsumer {
	clog := log.With(
		slog.String("owner", "RabbitConsumer"),
	)

	return &RabbitConsumer{
		log:               clog,
		eventService:      eventService,
		connectionManager: mananger,
	}
}

func (c *RabbitConsumer) Start(ctx context.Context, queueName string) {
	go func() {
		for {
			select {
			case <-ctx.Done():
				c.log.Info("consumer stopped by context")
				return
			default:
				if err := c.consumeLoop(ctx, queueName); err != nil {
					c.log.Error("consume loop error", liblogger.Err(err))
					time.Sleep(5 * time.Second) // задержка перед переподключением
				}
			}
		}
	}()

}

func (c *RabbitConsumer) consumeLoop(ctx context.Context, queueName string) error {
	rabbitConnect, err := c.connectionManager.GetConnection()
	if err != nil {
		return fmt.Errorf("amqp dial failed: %w", err)
	}
	c.log.Info("RabbitMQ connected success")

	closeConnectErrChan := make(chan *amqp.Error)
	rabbitConnect.NotifyClose(closeConnectErrChan)

	rabbitChannel, err := rabbitConnect.Channel()
	if err != nil {
		c.log.Error("failed create channel for RabbitMQ")
		return err
	}

	// Обработка закрытия канала
	closeErrChan := make(chan *amqp.Error)
	rabbitChannel.NotifyClose(closeErrChan)

	// Если очереди не существует то зависает тут
	msgs, err := rabbitChannel.Consume(
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
		return err
	}

	for {
		select {
		case <-ctx.Done():
			_ = rabbitChannel.Cancel("", false)
			return nil
		case <-closeErrChan:
			c.log.Error("rabbit channel closed")
			return fmt.Errorf("channel closed")
		case <-closeConnectErrChan:
			c.log.Error("rabbit connection closed")
			return fmt.Errorf("connection closed: %w", err)
		case msg, ok := <-msgs:
			if !ok {
				return fmt.Errorf("messages channel closed: %v", ok)
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
}

func (c *RabbitConsumer) handler(ctx context.Context, rabbitDTO rabbit_dto.RabbitDTO) error {
	id := ""
	if rabbitDTO.Method == UPDATE || rabbitDTO.Method == DELETE {
		var ok bool
		id, ok = rabbitDTO.Data.SearchAttributes["id"].(string)
		if !ok {
			c.log.Error("failed id does not exist", slog.String("id", id))
			return fmt.Errorf("failed search: ID does not exist")
		}
		c.log.Debug("id", slog.String("id", id))
	}

	switch rabbitDTO.Method {
	case CREATE:
		err := c.create(ctx, rabbitDTO.Data.Table, rabbitDTO.Data.Attributes)
		if err != nil {
			c.log.Error("failed create", liblogger.Err(err))
			return fmt.Errorf("failed create data")
		}
		c.log.Debug("success create", rabbitDTO.Data.Table, rabbitDTO.Data.Attributes)
	case UPDATE:
		err := c.update(ctx, rabbitDTO.Data.Table, rabbitDTO.Data.Attributes, id)
		if err != nil {
			c.log.Error("failed update", liblogger.Err(err))
			return fmt.Errorf("failed update data")
		}
		c.log.Debug("success update", rabbitDTO.Data.Table, rabbitDTO.Data.Attributes)
	case DELETE:
		err := c.delete(ctx, rabbitDTO.Data.Table, id)
		if err != nil {
			c.log.Error("failed delete", liblogger.Err(err))
			return fmt.Errorf("failed delete data")
		}
		c.log.Debug("success delete", rabbitDTO.Data.Table, rabbitDTO.Data.Attributes)
	default:
		c.log.Error("unexpected method", slog.String("method", rabbitDTO.Method))
		return fmt.Errorf("unexpected method %s", rabbitDTO.Method)
	}

	return nil
}

func MapToStructViaJSON(data map[string]any, out any) error {
	bytes, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("failed to marshal map to JSON: %w", err)
	}

	if err := json.Unmarshal(bytes, out); err != nil {
		return fmt.Errorf("failed to unmarshal JSON to struct: %w", err)
	}

	return nil
}

func (c *RabbitConsumer) create(ctx context.Context, tableName string, data map[string]any) error {

	switch tableName {
	case EVENT_TABLE:
		eventDto := event_dto.CreateEventDTORequest{}
		if err := MapToStructViaJSON(data, &eventDto); err != nil {
			return fmt.Errorf("failed convert data to dto: %w", err)
		}
		_, err := c.eventService.CreateEvent(ctx, eventDto)
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
		eventDTO := event_dto.UpdateEventDTORequest{}
		if err := MapToStructViaJSON(data, eventDTO); err != nil {
			return fmt.Errorf("failed convert data to dto")
		}
		err := c.eventService.UpdateEvent(ctx, id, eventDTO)
		if err != nil {
			return fmt.Errorf("falled update event: %w", err)
		}

	default:
		return fmt.Errorf("unexpected table: %s", tableName)
	}

	return nil
}

func (c *RabbitConsumer) delete(ctx context.Context, tableName string, id string) error {
	switch tableName {
	case EVENT_TABLE:
		err := c.eventService.DeleteEvent(ctx, id)
		if err != nil {
			return fmt.Errorf("failed delete event: %w", err)
		}
	default:
		return fmt.Errorf("unexpected table: %s", tableName)
	}

	return nil
}
