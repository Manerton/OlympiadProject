package rabbitmq

import (
	"context"
	"fmt"
	"log/slog"
	"main/internal/lib/liblogger"
	"sync"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

type ConnectionManager struct {
	connectionAddress string
	connect           *amqp.Connection
	log               *slog.Logger

	mu sync.RWMutex
}

func New(connection string, log *slog.Logger) *ConnectionManager {
	mlog := log.With(
		slog.String("op", "ConnectionManager"),
	)

	manager := &ConnectionManager{
		connectionAddress: connection,
		log:               mlog,
	}
	return manager
}

func (m *ConnectionManager) Start(ctx context.Context) {
	go func() {
		for {
			select {
			case <-ctx.Done():
				m.log.Info("consumer stopped by context")
				return
			default:
				err := m.ensureConnection()
				if err != nil {
					m.log.Error("failed to connect to RabbitMQ", liblogger.Err(err))
					time.Sleep(5 * time.Second)
					continue
				}
			}
		}
	}()
}

func (m *ConnectionManager) ensureConnection() error {
	connect, err := amqp.Dial(m.connectionAddress)
	if err != nil {
		return fmt.Errorf("amqp dial failed: %w", err)
	}

	m.mu.Lock()

	m.connect = connect
	notifyClose := make(chan *amqp.Error)
	m.connect.NotifyClose(notifyClose)

	m.mu.Unlock()
	m.log.Info("RabbitMQ connected")

	// ожидаем закрытия соединения
	err = <-notifyClose
	m.log.Info("RabbitMQ connection closed", liblogger.Err(err))

	return nil
}

func (m *ConnectionManager) GetConnection() (*amqp.Connection, error) {
	m.mu.RLock()
	connect := m.connect
	m.mu.RUnlock()

	if connect != nil && !connect.IsClosed() {
		return connect, nil
	}

	return nil, fmt.Errorf("failed connect")
}
