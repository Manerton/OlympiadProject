package consumer

import (
	"context"
	"encoding/json"
	"log"
	"log/slog"
	register_dto "main/internal/dto/auth/register"
	rabbit_dto "main/internal/dto/rabbit"
	user_dto "main/internal/dto/user"
	"main/internal/lib/liblogger"

	amqp "github.com/rabbitmq/amqp091-go"
)

// command
const (
	CREATE = "create"
	UPDATE = "update"
	DELETE = "delete"
)

// table name
const (
	USER_TABLE        = "user"
	PARTICIPANT_TABLE = "participant"
)

type UserService interface {
	Update(ctx context.Context, userDto *user_dto.UpdateUserRequestDTO) error
	Delete(ctx context.Context, id string) error
}

type AuthService interface {
	RegisterUser(ctx context.Context, registerUser *register_dto.RegusterUserRequestDTO) error
	RegisterParticipant(ctx context.Context, registerRequst *register_dto.RegisterParticipantRequestDTO) error
}

type RabbitConsumer struct {
	log          *slog.Logger
	rabbitCannel *amqp.Channel
	userService  UserService
	authService  AuthService
}

// construct
func New(log *slog.Logger, channel *amqp.Channel, userService UserService, authService AuthService) *RabbitConsumer {
	const op = "RabbitMQ consumer"

	clog := log.With(
		slog.String("op", op),
	)

	return &RabbitConsumer{
		log:          clog,
		rabbitCannel: channel,
		userService:  userService,
		authService:  authService,
	}
}

func (c *RabbitConsumer) Start(ctx context.Context, queueName string) {

	msgs, err := c.rabbitCannel.Consume(
		queueName,
		"",
		false,
		false,
		false,
		false,
		nil,
	)

	if err != nil {
		c.log.Error("Failed consume queue: %s, %w", queueName, liblogger.Err(err))
	}

	go func() {
		for {
			select {
			case <-ctx.Done():
				log.Println()
				if err := c.rabbitCannel.Cancel("", false); err != nil {
					c.log.Error("failed to cancel consumer: %w", liblogger.Err(err))
				}
				return
			case msg, ok := <-msgs:
				if !ok {
					c.log.Info("message channel closed")
					return
				}

				rabbitDTO := rabbit_dto.RabbitDTO{}
				if err := json.Unmarshal(msg.Body, &rabbitDTO); err != nil {
					c.log.Error("invalid message format: %w", liblogger.Err(err))
					msg.Nack(false, false)
					continue
				}

				err = c.handler(ctx, rabbitDTO)
				if err != nil {
					c.log.Error("task failed: %w", liblogger.Err(err))
					msg.Nack(false, false)
					continue
				}

				msg.Ack(false)
			}
		}
	}()
}

func (c *RabbitConsumer) handler(ctx context.Context, rabbitDTO rabbit_dto.RabbitDTO) error {
	switch rabbitDTO.Method {
	case CREATE:
		c.create(ctx, rabbitDTO.Data)
	case UPDATE:
		c.update()
	case DELETE:
		c.delete()
	}

	return nil
}

func (c *RabbitConsumer) create(ctx context.Context, rabbitData rabbit_dto.RabbitData) {
	dataJSON, err := json.Marshal(rabbitData.Attributes)
	if err != nil {
		c.log.Error("failed convert from data to json: %w", liblogger.Err(err))
	}

	switch rabbitData.Table {
	case USER_TABLE:
		userDTO := register_dto.RegusterUserRequestDTO{}
		if err := json.Unmarshal(dataJSON, &userDTO); err != nil {
			c.log.Error("failed parse from json: %w", liblogger.Err(err))
		}

		err = c.authService.RegisterUser(ctx, &userDTO)
		if err != nil {
			c.log.Error("failed register user: %w", liblogger.Err(err))
		}
	case PARTICIPANT_TABLE:
		participantDTO := register_dto.RegisterParticipantRequestDTO{}
		if err := json.Unmarshal(dataJSON, &participantDTO); err != nil {
			c.log.Error("failed parse from json: %w", liblogger.Err(err))
		}

		err = c.authService.RegisterParticipant(ctx, &participantDTO)
		if err != nil {
			c.log.Error("failed register user: %w", liblogger.Err(err))
		}
	default:
		c.log.Error("unexpected table: %s", slog.String("table", rabbitData.Table))
	}

}

func (c *RabbitConsumer) update() {

}

func (c *RabbitConsumer) delete() {

}
