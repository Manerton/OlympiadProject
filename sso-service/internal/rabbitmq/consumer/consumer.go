package consumer

import (
	"context"
	"encoding/json"
	"log"
	register_dto "main/internal/dto/auth/register"
	rabbit_dto "main/internal/dto/rabbit"
	user_dto "main/internal/dto/user"

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

const op = "RabbitMQ consumer"

type UserService interface {
	Update(ctx context.Context, userDto *user_dto.UpdateUserRequestDTO) error
	Delete(ctx context.Context, id string) error
}

type AuthService interface {
	RegisterUser(ctx context.Context, registerUser *register_dto.RegusterUserRequestDTO) error
	RegisterParticipant(ctx context.Context, registerRequst *register_dto.RegisterParticipantRequestDTO) error
}

type RabbitConsumer struct {
	rabbitCannel *amqp.Channel
	userService  UserService
	authService  AuthService
}

// construct
func New(channel *amqp.Channel, userService UserService, authService AuthService) *RabbitConsumer {
	return &RabbitConsumer{
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
		log.Printf("Failed consume queue: %s", queueName)
	}

	go func() {
		for {
			select {
			case <-ctx.Done():
				log.Println()
				if err := c.rabbitCannel.Cancel("", false); err != nil {
					log.Printf("%s: failed to cancel consumer", op)
				}
				return
			case msg, ok := <-msgs:
				if !ok {
					log.Printf("%s: message channel closed", op)
					return
				}

				rabbitDTO := rabbit_dto.RabbitDTO{}
				if err := json.Unmarshal(msg.Body, &rabbitDTO); err != nil {
					log.Printf("%s: invalid message format: %v", op, err)
					msg.Nack(false, false)
					continue
				}

				err = c.handler(ctx, rabbitDTO)
				if err != nil {
					log.Printf("%s: task failed: %v", op, err)
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
		log.Printf("%s: failed convert from data to json: %v", op, err)
	}

	switch rabbitData.Table {
	case USER_TABLE:
		userDTO := register_dto.RegusterUserRequestDTO{}
		if err := json.Unmarshal(dataJSON, &userDTO); err != nil {
			log.Printf("%s: failed parse from json: %v", op, err)
		}

		err = c.authService.RegisterUser(ctx, &userDTO)
		if err != nil {
			log.Printf("%s: failed register user: %v", op, err)
		}
	case PARTICIPANT_TABLE:
		participantDTO := register_dto.RegisterParticipantRequestDTO{}
		if err := json.Unmarshal(dataJSON, &participantDTO); err != nil {
			log.Printf("%s: failed parse from json: %v", op, err)
		}

		err = c.authService.RegisterParticipant(ctx, &participantDTO)
		if err != nil {
			log.Printf("%s: failed register user: %v", op, err)
		}
	}

}

func (c *RabbitConsumer) update() {

}

func (c *RabbitConsumer) delete() {

}
