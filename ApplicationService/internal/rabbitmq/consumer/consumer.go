package consumer

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"log/slog"
	"main/internal/dto/ApplicationDto"
	"main/internal/dto/RabbitDto"
	"main/internal/lib/liblogger"

	"github.com/google/uuid"
	"github.com/mitchellh/mapstructure"

	amqp "github.com/rabbitmq/amqp091-go"
)

const (
	APPLICATIONS_TABLE = "applications"
)

// command
const (
	GET    = ""
	CREATE = "create"
	UPDATE = "update"
	DELETE = "delete"
)

type ApplicationService interface {
	CreateApplication(ctx context.Context, applicationDTO ApplicationDto.CreateApplicationDTO) (uuid.UUID, error)
	UpdateApplication(ctx context.Context, id string, statusDTO ApplicationDto.UpdateApplicationDTO) error
	DeleteApplication(ctx context.Context, id string) error
	DeleteByFilter(ctx context.Context, deleteDTO ApplicationDto.DeleteApplicationDTO) error
}

type RabbitConsumer struct {
	log                *slog.Logger
	rabbitChannel      *amqp.Channel
	applicationService ApplicationService
}

func New(log *slog.Logger, channel *amqp.Channel, applicationService ApplicationService) *RabbitConsumer {
	clog := log.With(
		slog.String("op", "RabbitConsumer"),
	)

	return &RabbitConsumer{
		log:                clog,
		rabbitChannel:      channel,
		applicationService: applicationService,
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

				rabbitDTO := RabbitDto.RabbitDTO{}

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

func (c *RabbitConsumer) handler(ctx context.Context, rabbitDTO RabbitDto.RabbitDTO) error {

	switch rabbitDTO.Method {
	case CREATE:
		err := c.create(ctx, rabbitDTO.Data.Table, rabbitDTO.Data.Attributes)
		if err != nil {
			c.log.Error("failed create data", liblogger.Err(err))
			return fmt.Errorf("failed create data")
		}
		c.log.Debug("success create")
	case UPDATE:
		id, ok := rabbitDTO.Data.SearchAttributes["id"].(string)
		if !ok {
			c.log.Error("failed id does not exist", slog.String("id", id))
			return fmt.Errorf("failed update: ID does not exist")
		}
		c.log.Debug("id", slog.String("id", id))

		err := c.update(ctx, rabbitDTO.Data.Table, rabbitDTO.Data.Attributes, id)
		if err != nil {
			c.log.Error("failed update data", liblogger.Err(err))
			return fmt.Errorf("failed upadate data")
		}
		c.log.Debug("success update")

	case DELETE:
		err := c.delete(ctx, rabbitDTO.Data.Table, rabbitDTO.Data.SearchAttributes)
		if err != nil {
			c.log.Error("failed delete data", liblogger.Err(err))
			return fmt.Errorf("failed delete data")
		}
		c.log.Debug("success delete")
	}
	return nil
}

func (c *RabbitConsumer) create(ctx context.Context, tableName string, data map[string]any) error {
	switch tableName {
	case APPLICATIONS_TABLE:
		appDTO := ApplicationDto.CreateApplicationDTO{}
		if err := mapstructure.Decode(data, &appDTO); err != nil {
			return fmt.Errorf("failed parse json: %w", err)
		}

		_, err := c.applicationService.CreateApplication(ctx, appDTO)
		if err != nil {
			return fmt.Errorf("%w", err)
		}

	default:
		return fmt.Errorf("unexpected table: %s", tableName)
	}

	return nil
}
func (c *RabbitConsumer) update(ctx context.Context, tableName string, data map[string]any, id string) error {
	switch tableName {
	case APPLICATIONS_TABLE:
		appDTO := ApplicationDto.UpdateApplicationDTO{}
		if err := mapstructure.WeakDecode(data, &appDTO); err != nil {
			return fmt.Errorf("failed parse json: %w", err)
		}

		log.Println("id", id)

		err := c.applicationService.UpdateApplication(ctx, id, appDTO)
		if err != nil {
			return fmt.Errorf("%w", err)
		}
	default:
		return fmt.Errorf("unexpected table: %s", tableName)
	}

	return nil
}

func (c *RabbitConsumer) delete(ctx context.Context, tableName string, data map[string]any) error {
	switch tableName {
	case APPLICATIONS_TABLE:
		appDTO := ApplicationDto.DeleteApplicationDTO{}
		if err := mapstructure.WeakDecode(data, &appDTO); err != nil {
			return fmt.Errorf("failed parse to DTO: %w", err)
		}

		err := c.applicationService.DeleteByFilter(ctx, appDTO)
		if err != nil {
			return fmt.Errorf("%w", err)
		}

	default:
		return fmt.Errorf("unexpected table: %s", tableName)
	}

	return nil
}
