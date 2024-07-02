package app

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/rabbitmq/amqp091-go"
	"github.com/rs/zerolog"
	"github.com/uptrace/uptrace-go/uptrace"
	"go.opentelemetry.io/otel/trace"
)

type Dependencies struct {
	S3Client            *s3.Client
	S3PresignClient     *s3.PresignClient
	RabbitMqChannel     *amqp091.Channel
	Postgres            *sqlx.DB
	ElasticsearchClient *elasticsearch.TypedClient
	SessionStore        *session.Store
	Logger              zerolog.Logger
	Tracer              trace.Tracer
}

func NewDependencies(ctx context.Context) (*Dependencies, error) {
	logger := createLogger()
	logger.Info().Msg("Creating dependencies")

	logger.Info().Msg("Initializing OpenTelemetry")
	tracer, err := createTracer(ctx)
	if err != nil {
		logger.Fatal().Err(err).Msg("Failed to initialize OpenTelemetry")
		return nil, err
	}

	logger.Info().Msg("Creating S3 client")
	s3Client, err := createS3Client()
	if err != nil {
		logger.Fatal().Err(err).Msg("Failed to create S3 client")
		return nil, err
	}

	s3PresignClient := createS3PresignClient(s3Client)

	logger.Info().Msg("Creating RabbitMQ channel")
	rabbitMqChannel, err := createRabbitMqChannel()
	if err != nil {
		logger.Fatal().Err(err).Msg("Failed to create RabbitMQ channel")
		return nil, err
	}

	logger.Info().Msg("Creating Postgres connection")
	db, err := createPostgresConnection(logger)
	if err != nil {
		logger.Fatal().Err(err).Msg("Failed to create Postgres connection")
		return nil, err
	}

	logger.Info().Msg("Creating Elasticsearch client")
	elasticsearchClient, err := createElasticSearchClient()
	if err != nil {
		logger.Fatal().Err(err).Msg("Failed to create Elasticsearch client")
		return nil, err
	}

	logger.Info().Msg("Creating Redis client")
	sessionStore, err := createSessionStore()
	if err != nil {
		logger.Fatal().Err(err).Msg("Failed to create Redis client")
		return nil, err
	}

	return &Dependencies{
		S3Client:            s3Client,
		S3PresignClient:     s3PresignClient,
		RabbitMqChannel:     rabbitMqChannel,
		Postgres:            db,
		ElasticsearchClient: elasticsearchClient,
		SessionStore:        sessionStore,
		Logger:              logger,
		Tracer:              tracer,
	}, nil
}

func (d *Dependencies) Close(ctx context.Context) error {
	if d.Postgres != nil {
		d.Postgres.Close()
	}

	if d.RabbitMqChannel != nil {
		d.Logger.Info().Msg("Closing RabbitMQ channel")
		d.RabbitMqChannel.Close()
	}

	uptrace.Shutdown(ctx)

	return nil
}
