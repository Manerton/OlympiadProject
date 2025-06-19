package consumer

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	register_dto "main/internal/dto/auth/register"
	participant_dto "main/internal/dto/participant"
	rabbit_dto "main/internal/dto/rabbit"
	school_dto "main/internal/dto/school"
	user_dto "main/internal/dto/user"
	"main/internal/lib/liblogger"
	"main/internal/rabbitmq/producer"

	amqp "github.com/rabbitmq/amqp091-go"
)

// command
const (
	GET    = ""
	CREATE = "create"
	UPDATE = "update"
	DELETE = "delete"
)

// table name
const (
	USER_TABLE        = "user"
	PARTICIPANT_TABLE = "participant"
	SCHOOL_TABLE      = "school"
)

type UserService interface {
	GetByFilter(ctx context.Context, userDTO user_dto.SearchAttributesUserDTO) (user_dto.UserResponseDTO, error)

	Update(ctx context.Context, id string, userDto user_dto.UpdateUserRequestDTO) error
	Delete(ctx context.Context, id string) error
}

type ParticipantService interface {
	Update(ctx context.Context, id string, participantDTO participant_dto.UpdateParticipantRequestDTO) error
}

type AuthService interface {
	RegisterUser(ctx context.Context, registerUser *register_dto.RegisterUserRequestDTO) error
	RegisterParticipant(ctx context.Context, registerRequst *register_dto.RegisterParticipantRequestDTO) error
}

type SchoolService interface {
	Create(ctx context.Context, schoolDTO school_dto.CreateSchoolRequestDTO) error
	Update(ctx context.Context, id string, schoolDTO school_dto.UpdateSchoolRequestDTO) error
	Delete(ctx context.Context, id string) error
}

type RabbitConsumer struct {
	log                *slog.Logger
	rabbitCannel       *amqp.Channel
	userService        UserService
	participantService ParticipantService
	authService        AuthService
	schoolService      SchoolService
}

// construct
func New(log *slog.Logger, channel *amqp.Channel,
	userService UserService, participantService ParticipantService, authService AuthService, schoolService SchoolService) *RabbitConsumer {
	const op = "RabbitMQ consumer"

	clog := log.With(
		slog.String("op", op),
	)

	return &RabbitConsumer{
		log:                clog,
		rabbitCannel:       channel,
		userService:        userService,
		participantService: participantService,
		authService:        authService,
		schoolService:      schoolService,
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
		c.log.Error("Failed consume queue", slog.String("queue name", queueName), liblogger.Err(err))
		return
	}

	go func() {
		for {
			select {
			case <-ctx.Done():
				c.log.Info("init cancel consumer")
				if err := c.rabbitCannel.Cancel("", false); err != nil {
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
	attributesDataJSON, err := json.Marshal(rabbitDTO.Data.Attributes)
	if err != nil {
		c.log.Error("failed convert Attributes from data to json", liblogger.Err(err))
	}

	searchDataJSON, err := json.Marshal(rabbitDTO.Data.SearchAttributes)
	if err != nil {
		c.log.Error("failed convert SearchAttributes from data to json", liblogger.Err(err))
	}

	id := ""
	if rabbitDTO.Method == UPDATE || rabbitDTO.Method == DELETE {
		id, ok := rabbitDTO.Data.SearchAttributes["id"].(string)
		if !ok {
			c.log.Error("failed id does not exist", slog.String("id", id))
			return fmt.Errorf("failed update: ID does not exist")
		}
	}

	switch rabbitDTO.Method {
	case GET:
		result, err := c.get(ctx, rabbitDTO.Data.Table, searchDataJSON)
		if err != nil {
			c.log.Error("failed get", liblogger.Err(err))
			return fmt.Errorf("failed get data")
		}
		producer.SendToQueue(c.rabbitCannel, rabbitDTO.AppName, result)
	case CREATE:
		err := c.create(ctx, rabbitDTO.Data.Table, attributesDataJSON)
		if err != nil {
			c.log.Error("failed create", liblogger.Err(err))
			c.log.Debug("dataJSON", slog.Any("data", rabbitDTO.Data.Attributes))
			return fmt.Errorf("failed create data")
		}
	case UPDATE:
		err := c.update(ctx, rabbitDTO.Data.Table, attributesDataJSON, id)
		if err != nil {
			c.log.Error("failed update", liblogger.Err(err))
			return fmt.Errorf("failed update data")
		}
	case DELETE:
		err := c.delete(ctx, rabbitDTO.Data.Table, id)
		if err != nil {
			c.log.Error("failed delete", liblogger.Err(err))
			return fmt.Errorf("failed delete data")
		}
	}

	return nil
}

func (c *RabbitConsumer) get(ctx context.Context, tableName string, dataJSON []byte) ([]byte, error) {
	jsonResult := []byte{}
	switch tableName {
	case USER_TABLE:
		userDTO := user_dto.SearchAttributesUserDTO{}
		if err := json.Unmarshal(dataJSON, &userDTO); err != nil {
			return nil, fmt.Errorf("failed parse from json: %w", err)
		}

		userResult, err := c.userService.GetByFilter(ctx, userDTO)
		if err != nil {
			return nil, fmt.Errorf("%w", err)
		}

		jsonResult, err = json.Marshal(userResult)
		if err != nil {
			return nil, fmt.Errorf("failed convert data to json: %w", err)
		}
	case PARTICIPANT_TABLE:
		// TODO!! realise another tables
	default:
		return nil, fmt.Errorf("unexpected table: %s", tableName)
	}
	return jsonResult, nil
}

func (c *RabbitConsumer) create(ctx context.Context, tableName string, dataJSON []byte) error {
	switch tableName {
	case USER_TABLE:
		userDTO := register_dto.RegisterUserRequestDTO{}
		if err := userDTO.UnmarshalJSON(dataJSON); err != nil {
			return fmt.Errorf("failed parse json: %w", err)
		}

		err := c.authService.RegisterUser(ctx, &userDTO)
		if err != nil {
			return fmt.Errorf("%w", err)
		}
	case PARTICIPANT_TABLE:
		participantDTO := register_dto.RegisterParticipantRequestDTO{}
		if err := json.Unmarshal(dataJSON, &participantDTO); err != nil {
			return fmt.Errorf("failed parse json: %w", err)
		}

		err := c.authService.RegisterParticipant(ctx, &participantDTO)
		if err != nil {
			return fmt.Errorf("%w", err)
		}
	case SCHOOL_TABLE:
		schoolDTO := school_dto.CreateSchoolRequestDTO{}
		if err := json.Unmarshal(dataJSON, &schoolDTO); err != nil {
			return fmt.Errorf("failed parse json: %w", err)
		}

		err := c.schoolService.Create(ctx, schoolDTO)
		if err != nil {
			return fmt.Errorf("%w", err)
		}
	default:
		return fmt.Errorf("unexpected table: %s", tableName)
	}
	return nil
}

func (c *RabbitConsumer) update(ctx context.Context, tableName string, dataJSON []byte, id string) error {
	switch tableName {
	case USER_TABLE:
		userDTO := user_dto.UpdateUserRequestDTO{}
		err := json.Unmarshal(dataJSON, &userDTO)
		if err != nil {
			return fmt.Errorf("failed parse json: %w", err)
		}

		err = c.userService.Update(ctx, id, userDTO)
		if err != nil {
			return fmt.Errorf("%w", err)
		}
	case PARTICIPANT_TABLE:
		participantDTO := participant_dto.UpdateParticipantRequestDTO{}
		err := json.Unmarshal(dataJSON, &participantDTO)
		if err != nil {
			return fmt.Errorf("failed parse json: %w", err)
		}

		err = c.participantService.Update(ctx, id, participantDTO)
		if err != nil {
			return fmt.Errorf("%w", err)
		}
	case SCHOOL_TABLE:
		schoolDTO := school_dto.UpdateSchoolRequestDTO{}
		err := json.Unmarshal(dataJSON, &schoolDTO)
		if err != nil {
			return fmt.Errorf("failed parse json: %w", err)
		}

		err = c.schoolService.Update(ctx, id, schoolDTO)
		if err != nil {
			return fmt.Errorf("%w", err)
		}
	default:
		return fmt.Errorf("unexpected table: %s", tableName)
	}
	return nil
}

func (c *RabbitConsumer) delete(ctx context.Context, tableName string, id string) error {
	switch tableName {
	case USER_TABLE:
		err := c.userService.Delete(ctx, id)
		if err != nil {
			return fmt.Errorf("%w", err)
		}
	case SCHOOL_TABLE:
		err := c.schoolService.Delete(ctx, id)
		if err != nil {
			return fmt.Errorf("%w", err)
		}
	default:
		return fmt.Errorf("unexpected table: %s", tableName)
	}
	return nil
}
