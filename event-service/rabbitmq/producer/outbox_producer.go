package producer

import (
	"context"
	"fmt"
	"log/slog"
	"main/internal/lib/liblogger"
	"main/internal/models/outbox"
	"main/internal/storage/orm"
	"main/rabbitmq"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

type OutboxRepositoryInteface interface {
	GetAllNotProcessed(ctx context.Context, orm orm.ORM, limit *int) ([]outbox.Outbox, error)
	Update(ctx context.Context, orm orm.ORM, outbox outbox.Outbox) error
}

type OutboxProducer struct {
	db                orm.ORM
	log               *slog.Logger
	outboxRepository  OutboxRepositoryInteface
	connectionManager *rabbitmq.ConnectionManager
}

func New(log *slog.Logger, manager *rabbitmq.ConnectionManager, orm orm.ORM, repository OutboxRepositoryInteface) *OutboxProducer {
	return &OutboxProducer{
		log:               log,
		db:                orm,
		connectionManager: manager,
		outboxRepository:  repository,
	}
}

func (p *OutboxProducer) Start(ctx context.Context) {
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			default:
				outboxMessages, err := p.outboxRepository.GetAllNotProcessed(ctx, p.db, nil)
				if err != nil {
					p.log.Error("failed get outbox message", liblogger.Err(err))
					continue
				}

				if len(outboxMessages) == 0 {
					p.log.Info("outbox message is empty")
					time.Sleep(5 * time.Second)
					continue
				}

				connection, err := p.connectionManager.GetConnection()
				if err != nil {
					// p.log.Error("failed get connection", liblogger.Err(err))
					continue
				}

				channel, err := connection.Channel()
				if err != nil {
					p.log.Error("failed get channel", liblogger.Err(err))
					continue
				}

				for _, outboxMessage := range outboxMessages {
					p.log.Debug("test")
					err := p.Send(outboxMessage.QueueName, channel, outboxMessage.Payload)
					if err != nil {
						p.log.Error("failed send outbox message to rabbit", liblogger.Err(err))
					} else {
						p.outboxRepository.Update(ctx, p.db, outbox.Outbox{Processed: true})
					}
				}
			}
		}
	}()
}

func (p *OutboxProducer) Send(queueName string, channel *amqp.Channel, data []byte) error {

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := channel.PublishWithContext(
		ctx, "", queueName, false, false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        data,
		},
	)
	if err != nil {
		return fmt.Errorf("failed publish message: %w", err)
	} else {
		p.log.Info("message was send to: ", slog.String("queue name", queueName))
	}
	return nil
}
