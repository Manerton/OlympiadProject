package rabbitmq

import (
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

func MustConnect(connectionStr string) *amqp.Connection {
	connect, err := amqp.Dial(connectionStr)
	if err != nil {
		log.Fatalf("Failed connect rabbitmq %s", connectionStr)
	}

	return connect
}
