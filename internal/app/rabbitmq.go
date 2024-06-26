package app

import (
	"errors"
	"os"

	"github.com/rabbitmq/amqp091-go"
)

var (
	ErrFailedToConnectToRabbitMq     = errors.New("failed to connect to RabbitMQ")
	ErrFailedToCreateRabbitMqChannel = errors.New("failed to create RabbitMQ channel")
)

func createRabbitMqChannel() (*amqp091.Channel, error) {
	conn, err := amqp091.Dial(os.Getenv("RABBITMQ_URL"))
	if err != nil {
		return nil, errors.Join(err, ErrFailedToConnectToRabbitMq)
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, errors.Join(err, ErrFailedToCreateRabbitMqChannel)
	}

	return ch, nil
}
