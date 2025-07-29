package consumer

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"main/internal/dto/juryAssignmentsDto"
	rabbit_dto "main/internal/dto/rabbit"
	"main/internal/lib/liblogger"
	"sync"
	"time"

	"github.com/google/uuid"
	amqp "github.com/rabbitmq/amqp091-go"
)

const (
	JURY_ASSIGNMENTS = "jury_assignments"
)

// command
const (
	GET    = ""
	CREATE = "create"
	UPDATE = "update"
	DELETE = "delete"
)

type JureAssignmentsService interface {
	Create(ctx context.Context, juryAssignment juryAssignmentsDto.CreateJuryAssignmentsDTO) (uuid.UUID, error)
	Update(ctx context.Context, id string, juryAssignment juryAssignmentsDto.UpdateJuryAssignmentsDTO) error
	DeleteByFields(ctx context.Context, juryAssignment juryAssignmentsDto.JuryAssignmentsResponseDTO) error
}

type RabbitConsumer struct {
	log               *slog.Logger
	rabbitConnect     *amqp.Connection
	connectionAddress string

	mutex       sync.Mutex
	juryService JureAssignmentsService
}

func New(log *slog.Logger, address string, juryService JureAssignmentsService) *RabbitConsumer {
	clog := log.With(
		slog.String("op", "RabbitConsumer"),
	)

	return &RabbitConsumer{
		log:               clog,
		juryService:       juryService,
		connectionAddress: address,
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
				err := c.ensureConnection()
				if err != nil {
					c.log.Error("failed to connect to RabbitMQ", liblogger.Err(err))
					time.Sleep(5 * time.Second)
					continue
				}

				if err := c.consumeLoop(ctx, queueName); err != nil {
					c.log.Error("consume loop error", liblogger.Err(err))
					time.Sleep(5 * time.Second) // задержка перед переподключением
				}
			}
		}
	}()

}

func (c *RabbitConsumer) ensureConnection() error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if c.rabbitConnect != nil && !c.rabbitConnect.IsClosed() {
		return nil // всё ок
	}

	conn, err := amqp.Dial(c.connectionAddress)
	if err != nil {
		return fmt.Errorf("amqp dial failed: %w", err)
	}
	c.rabbitConnect = conn

	c.log.Info("reconnected to RabbitMQ")
	return nil
}

func (c *RabbitConsumer) consumeLoop(ctx context.Context, queueName string) error {
	closeConnectErrChan := make(chan *amqp.Error)
	c.rabbitConnect.NotifyClose(closeConnectErrChan)

	rabbitChannel, err := c.rabbitConnect.Channel()
	if err != nil {
		return fmt.Errorf("amqp dial failed: %w", err)
	}
	c.log.Info("RabbitMQ connected success")

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

	switch rabbitDTO.Method {
	case CREATE:
		err := c.create(ctx, rabbitDTO.Data.Table, rabbitDTO.Data.Attributes)
		if err != nil {
			c.log.Error("failed create", liblogger.Err(err))
			return fmt.Errorf("failed create data")
		}
		c.log.Debug("success create", rabbitDTO.Data.Table, rabbitDTO.Data.Attributes)
	case UPDATE:
		id, ok := rabbitDTO.Data.SearchAttributes["id"].(string)
		if !ok {
			c.log.Error("failed id does not exist", slog.String("id", id))
			return fmt.Errorf("failed search: ID does not exist")
		}
		c.log.Debug("id", slog.String("id", id))

		err := c.update(ctx, rabbitDTO.Data.Table, rabbitDTO.Data.Attributes, id)
		if err != nil {
			c.log.Error("failed update", liblogger.Err(err))
			return fmt.Errorf("failed update data")
		}
		c.log.Debug("success update", rabbitDTO.Data.Table, rabbitDTO.Data.Attributes)
	case DELETE:
		err := c.delete(ctx, rabbitDTO.Data.Table, rabbitDTO.Data.SearchAttributes)
		if err != nil {
			c.log.Error("failed delete", liblogger.Err(err))
			return fmt.Errorf("failed delete data")
		}
		c.log.Debug("success delete", rabbitDTO.Data.Table, rabbitDTO.Data.Attributes)
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
	case JURY_ASSIGNMENTS:
		juryAssigmentDTO := juryAssignmentsDto.CreateJuryAssignmentsDTO{}
		err := MapToStructViaJSON(data, &juryAssigmentDTO)
		if err != nil {
			return fmt.Errorf("failed decode jury-assignment: %w", err)
		}

		_, err = c.juryService.Create(ctx, juryAssigmentDTO)
		if err != nil {
			return fmt.Errorf("failed create jury-assignment: %w", err)
		}
	default:
		return fmt.Errorf("unexpected table: %s", tableName)
	}

	return nil
}

func (c *RabbitConsumer) update(ctx context.Context, tableName string, data map[string]any, id string) error {
	switch tableName {
	case JURY_ASSIGNMENTS:
		juryAssignmentsDto := juryAssignmentsDto.UpdateJuryAssignmentsDTO{}
		err := MapToStructViaJSON(data, &juryAssignmentsDto)
		if err != nil {
			return fmt.Errorf("failed decode jury-assignment: %w", err)
		}

		err = c.juryService.Update(ctx, id, juryAssignmentsDto)
		if err != nil {
			return fmt.Errorf("failed update jury-assignment: %w", err)
		}
	default:
		return fmt.Errorf("unexpected table: %s", tableName)
	}

	return nil
}

func (c *RabbitConsumer) delete(ctx context.Context, tableName string, data map[string]any) error {
	switch tableName {
	case JURY_ASSIGNMENTS:
		juryAssignmentsDto := juryAssignmentsDto.JuryAssignmentsResponseDTO{}
		err := MapToStructViaJSON(data, &juryAssignmentsDto)
		if err != nil {
			return fmt.Errorf("failed decode jury-assignment: %w", err)
		}

		err = c.juryService.DeleteByFields(ctx, juryAssignmentsDto)
		if err != nil {
			return fmt.Errorf("failed delete jury-assignment: %w", err)
		}
	default:
		return fmt.Errorf("unexpected table: %s", tableName)
	}

	return nil
}
