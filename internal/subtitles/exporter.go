package subtitles

import (
	"context"
	"dewarrum/vocabulary-leveling/internal/app"
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/martinlindhe/subtitles"
	"github.com/rs/zerolog"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

var (
	ErrFailedToRunExporter = errors.New("failed to run exporter")
)

type Exporter struct {
	MessageQueue        *MessageQueue
	SubtitlesRepository *SubtitlesRepository
	FileStorage         *FileStorage
	FullTextSearch      *FullTextSearch
	Logger              zerolog.Logger
	Tracer              trace.Tracer
}

func NewExporter(dependencies *app.Dependencies, context context.Context) (*Exporter, error) {
	messageQueue, err := NewMessageQueue(dependencies)
	if err != nil {
		return nil, err
	}
	fullTextSearch, err := NewFullTextSearch(dependencies.ElasticsearchClient, context)
	if err != nil {
		return nil, err
	}

	return &Exporter{
		MessageQueue:        messageQueue,
		SubtitlesRepository: NewSubtitlesRepository(dependencies),
		FileStorage:         NewFileStorage(dependencies.S3Client, dependencies.S3PresignClient),
		FullTextSearch:      fullTextSearch,
		Logger:              dependencies.Logger,
		Tracer:              dependencies.Tracer,
	}, nil
}

func (e *Exporter) Run(ctx context.Context) error {
	e.Logger.Info().Msg("Starting subtitle exporter")
	ctx, span := e.Tracer.Start(ctx, "subtitles.exporter")
	defer span.End()

	messages, err := e.MessageQueue.Consume(ctx)
	if err != nil {
		e.Logger.Fatal().Err(err).Msg("Failed to register a consumer")
		span.RecordError(err, trace.WithStackTrace(true))

		return errors.Join(err, ErrFailedToRunExporter)
	}

	go func() {
		for message := range messages {
			ctx, span := e.Tracer.Start(ctx, "subtitles.exporter.handleMessage", trace.WithAttributes(attribute.String("videoId", message.VideoId.String())))
			err = e.handleMessage(message, ctx)
			if err != nil {
				e.Logger.Error().Str("videoId", message.VideoId.String()).Err(err).Msg("Failed to handle message")
				span.RecordError(err, trace.WithStackTrace(true))
			}
			span.End()
		}
	}()

	return nil
}

func (e *Exporter) handleMessage(message ExportSubtitlesMessage, ctx context.Context) error {
	subtitle, err := e.FileStorage.Download(message.VideoId.String(), ctx)
	if err != nil {
		return err
	}

	for _, caption := range subtitle.Captions {
		dbSubtitle, err := e.saveToDatabase(message.VideoId, caption, ctx)
		if err != nil {
			return err
		}

		err = e.saveToFullTextSearch(dbSubtitle.Id, dbSubtitle.VideoId, caption, ctx)
		if err != nil {
			return err
		}
	}

	return nil
}

func (e *Exporter) saveToDatabase(videoId uuid.UUID, caption subtitles.Caption, context context.Context) (*DbSubtitle, error) {
	emptyDate := (time.Time{}).AddDate(-1, 0, 0)
	subtitle := newDbSubtitle(videoId, strings.Join(caption.Text, "\n"), caption.Seq, caption.Start.Sub(emptyDate).Milliseconds(), caption.End.Sub(emptyDate).Milliseconds())

	inserted, err := e.SubtitlesRepository.Insert(subtitle, context)
	if err != nil {
		return nil, err
	}

	return inserted, nil
}

func (e *Exporter) saveToFullTextSearch(id string, videoId uuid.UUID, caption subtitles.Caption, context context.Context) error {
	subtitle := NewFtsSubtitle(id, videoId, caption.Seq, strings.Join(caption.Text, "\n"))
	err := e.FullTextSearch.Insert(subtitle, context)
	if err != nil {
		return err
	}

	return nil
}
