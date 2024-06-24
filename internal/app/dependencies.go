package app

import (
	"errors"
	"log"
	"os"

	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/rabbitmq/amqp091-go"
)

const (
	MissingConfigurationError = "missing configuration"
)

type Dependencies struct {
	S3Client        *s3.Client
	S3PresignClient *s3.PresignClient
	RabbitMqChannel *amqp091.Channel
	Postgres        *sqlx.DB
}

func NewDependencies() (*Dependencies, error) {
	s3Client, err := createS3Client()
	if err != nil {
		return nil, err
	}

	s3PresignClient := createS3PresignClient(s3Client)

	rabbitMqChannel, err := createRabbitMqChannel()
	if err != nil {
		return nil, err
	}

	db, err := createPostgresConnection()
	if err != nil {
		return nil, err
	}

	return &Dependencies{
		S3Client:        s3Client,
		S3PresignClient: s3PresignClient,
		RabbitMqChannel: rabbitMqChannel,
		Postgres:        db,
	}, nil
}

func createS3Client() (*s3.Client, error) {
	accessKey := os.Getenv("AWS_ACCESS_KEY_ID")
	secretKey := os.Getenv("AWS_SECRET_ACCESS_KEY")
	region := os.Getenv("AWS_REGION")
	endpointURL := os.Getenv("AWS_ENDPOINT_URL")

	if accessKey == "" || secretKey == "" || region == "" || endpointURL == "" {
		log.Fatal("AWS_ACCESS_KEY_ID, AWS_SECRET_ACCESS_KEY, AWS_REGION and AWS_ENDPOINT_URL must be set")
		return nil, errors.New(MissingConfigurationError)
	}

	client := s3.New(s3.Options{
		Credentials:  credentials.NewStaticCredentialsProvider(accessKey, secretKey, ""),
		Region:       region,
		BaseEndpoint: &endpointURL,
		UsePathStyle: true,
	})

	return client, nil
}

func createS3PresignClient(s3Client *s3.Client) *s3.PresignClient {
	return s3.NewPresignClient(s3Client)
}

func createRabbitMqChannel() (*amqp091.Channel, error) {
	conn, err := amqp091.Dial(os.Getenv("RABBITMQ_URL"))
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %s", err)
		return nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("Failed to open a channel: %s", err)
		return nil, err
	}

	return ch, nil
}

func createRabbitMqVideoExporterQueue(ch *amqp091.Channel) (*amqp091.Queue, error) {
	q, err := ch.QueueDeclare(
		"video-exporter", // name
		true,             // durable
		false,            // delete when unused
		false,            // exclusive
		false,            // no-wait
		nil,              // arguments
	)
	if err != nil {
		log.Fatalf("Failed to declare a queue: %s", err)
		return nil, err
	}

	return &q, nil
}

func createPostgresConnection() (*sqlx.DB, error) {
	db, err := sqlx.Connect("postgres", os.Getenv("POSTGRES_URL"))
	if err != nil {
		log.Fatalf("Failed to connect to Postgres: %s", err)
		return nil, err
	}

	return db, nil
}
