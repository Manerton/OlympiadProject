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

const (
	CREATE = "create"
	UPDATE = "update"
	DELETE = "delete"
)

const (
	USER_TABLE        = "user"
	PARTICIPANT_TABLE = "participant"
)

const op = "RabbitMQ consumer"

type UserService interface {
	Update(ctx context.Context, userDto *user_dto.UpdateUserRequestDTO) error
	Delete(ctx context.Context, id string) error
}

type RabbitConsumer struct {
	rabbitCannel *amqp.Channel
	userService  UserService
}

func New(channel *amqp.Channel, userService UserService) *RabbitConsumer {
	return &RabbitConsumer{
		rabbitCannel: channel,
		userService:  userService,
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

				err = c.handler(rabbitDTO)
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

func (c *RabbitConsumer) handler(rabbitDTO rabbit_dto.RabbitDTO) error {
	switch rabbitDTO.Method {
	case CREATE:
		c.Create(rabbitDTO.Data)
	case UPDATE:
		c.Update()
	case DELETE:
		c.Delete()
	}

	return nil
}

func (c *RabbitConsumer) Create(rabbitData rabbit_dto.RabbitData) {

	switch rabbitData.Table {
	case USER_TABLE:
		userJSON, err := json.Marshal(rabbitData.Attributes)
		if err != nil {
			log.Printf("%s: failed convert from data to json: %v", op, err)
		}

		userDTO := register_dto.RegusterUserRequestDTO{}
		if err := json.Unmarshal(userJSON, &userDTO); err != nil {
			log.Printf("%s: failed parse from json: %v", op, err)
		}

	case PARTICIPANT_TABLE:
	}

}

func (c *RabbitConsumer) Update() {

}

func (c *RabbitConsumer) Delete() {

}
