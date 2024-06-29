package subtitles

import (
	"context"
	"dewarrum/vocabulary-leveling/internal/app"
	"encoding/json"
	"errors"

	"github.com/google/uuid"
	"github.com/rabbitmq/amqp091-go"
	"github.com/rs/zerolog"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

const (
	FailedToDeclareQueue    = "failed to declare queue"
	FailedToSendMessage     = "failed to send message"
	FailedToConsume         = "failed to consume"
	FailedToDeclareExchange = "failed to declare exchange"
	FailedToBindQueue       = "failed to bind queue"
	queue                   = "queue.subtitles.export"
	exchange                = "exchange.subtitles.export"
)

type ExportSubtitlesMessage struct {
	VideoId uuid.UUID `json:"videoId"`
}

func NewExportSubtitlesMessage(videoId uuid.UUID) *ExportSubtitlesMessage {
	return &ExportSubtitlesMessage{
		VideoId: videoId,
	}
}

type MessageQueue struct {
	channel *amqp091.Channel
	queue   *amqp091.Queue
	logger  zerolog.Logger
	tracer  trace.Tracer
}

func (mq *MessageQueue) Send(message *ExportSubtitlesMessage, context context.Context) error {
	_, span := mq.tracer.Start(context, "mq.send.subtitles.export")
	span.SetAttributes(attribute.String("videoId", message.VideoId.String()))
	defer span.End()

	mq.logger.Info().Str("videoId", message.VideoId.String()).Msg("Sending message")

	body, err := json.Marshal(message)
	if err != nil {
		span.RecordError(err, trace.WithStackTrace(true))
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
		span.RecordError(err, trace.WithStackTrace(true))
		return errors.Join(err, errors.New(FailedToSendMessage))
	}

	mq.logger.Info().Str("videoId", message.VideoId.String()).Msg("Message sent successfully")

	return nil
}

func (mq *MessageQueue) Consume(ctx context.Context) (<-chan ExportSubtitlesMessage, error) {
	ctx, span := mq.tracer.Start(ctx, "mq.send.subtitles.export")
	defer span.End()

	mq.logger.Info().Msg("Starting to consume messages")

	deliveries, err := mq.channel.ConsumeWithContext(
		ctx,
		mq.queue.Name,
		"consumer.subtitles.export",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return nil, errors.Join(err, errors.New(FailedToConsume))
	}

	messages := make(chan ExportSubtitlesMessage)

	go func() {
		for delivery := range deliveries {
			var message ExportSubtitlesMessage
			err := json.Unmarshal(delivery.Body, &message)
			if err != nil {
				mq.logger.Error().Err(err).Msg("Failed to unmarshal message")
				continue
			}

			messages <- message
		}
	}()

	return messages, nil
}

func NewMessageQueue(dependencies *app.Dependencies) (*MessageQueue, error) {
	queue, err := createQueue(dependencies.RabbitMqChannel)
	if err != nil {
		return nil, err
	}

	err = dependencies.RabbitMqChannel.ExchangeDeclare(exchange, "direct", false, false, false, false, nil)
	if err != nil {
		return nil, errors.Join(err, errors.New(FailedToDeclareExchange))
	}

	err = dependencies.RabbitMqChannel.QueueBind(queue.Name, "", exchange, false, nil)
	if err != nil {
		return nil, errors.Join(err, errors.New(FailedToBindQueue))
	}

	return &MessageQueue{
		channel: dependencies.RabbitMqChannel,
		queue:   queue,
		logger:  dependencies.Logger,
		tracer:  dependencies.Tracer,
	}, nil
}

func createQueue(channel *amqp091.Channel) (*amqp091.Queue, error) {
	queue, err := channel.QueueDeclare(
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

	return &queue, nil
}
