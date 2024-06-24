package videos

import (
	"context"
	"encoding/json"
	"errors"
	"log"

	"github.com/rabbitmq/amqp091-go"
)

const (
	FailedToDeclareQueue    = "failed to declare queue"
	FailedToSendMessage     = "failed to send message"
	FailedToConsume         = "failed to consume"
	FailedToDeclareExchange = "failed to declare exchange"
	FailedToBindQueue       = "failed to bind queue"
	queue                   = "queue.videos.export"
	exchange                = "exchange.videos.export"
)

type MessageQueue struct {
	channel *amqp091.Channel
	queue   *amqp091.Queue
}

type ExportVideoMessage struct {
	VideoId string `json:"videoId"`
}

func NewExportVideoMessage(videoId string) *ExportVideoMessage {
	return &ExportVideoMessage{
		VideoId: videoId,
	}
}

func (mq *MessageQueue) Send(message *ExportVideoMessage, context context.Context) error {
	body, err := json.Marshal(message)
	if err != nil {
		return err
	}

	err = mq.channel.PublishWithContext(
		context,
		exchange, // exchange
		"",       // routing key
		false,    // mandatory
		false,    // immediate
		amqp091.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)
	if err != nil {
		return errors.Join(err, errors.New(FailedToSendMessage))
	}

	return nil
}

func (mq *MessageQueue) Consume() (<-chan ExportVideoMessage, error) {
	deliveries, err := mq.channel.Consume(
		mq.queue.Name,
		"consumer.videos.export",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return nil, errors.Join(err, errors.New(FailedToConsume))
	}

	messages := make(chan ExportVideoMessage)

	go func() {
		for delivery := range deliveries {
			var message ExportVideoMessage
			err := json.Unmarshal(delivery.Body, &message)
			if err != nil {
				log.Fatalf("Failed to unmarshal message: %s", err)
				continue
			}

			messages <- message
		}
	}()

	return messages, nil
}

func NewMessageQueue(channel *amqp091.Channel) (*MessageQueue, error) {
	queue, err := createQueue(channel)
	if err != nil {
		return nil, err
	}

	err = channel.ExchangeDeclare(exchange, "direct", true, false, false, false, nil)
	if err != nil {
		return nil, errors.Join(err, errors.New(FailedToDeclareExchange))
	}

	err = channel.QueueBind(queue.Name, "", exchange, false, nil)
	if err != nil {
		return nil, errors.Join(err, errors.New(FailedToBindQueue))
	}

	return &MessageQueue{
		channel: channel,
		queue:   queue,
	}, nil
}

func createQueue(channel *amqp091.Channel) (*amqp091.Queue, error) {
	q, err := channel.QueueDeclare(
		queue, // name
		true,  // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		return nil, errors.Join(err, errors.New(FailedToDeclareQueue))
	}

	return &q, nil
}
