package chunks

import (
	"context"
	"dewarrum/vocabulary-leveling/internal/app"
	"errors"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"github.com/rs/zerolog"
	"go.opentelemetry.io/otel/trace"
)

var (
	ErrChunkAlreadyExists = errors.New("chunk already exists")
	ErrFailedToGetChunks  = errors.New("failed to get chunks")
)

type DbChunk struct {
	Id               uuid.UUID `db:"id"`
	VideoId          uuid.UUID `db:"video_id"`
	RepresentationId string    `db:"representation_id"`
	Sequence         int       `db:"sequence"`
	ContentLocation  string    `db:"content_location"`
	StartMs          int64     `db:"start_ms"`
	EndMs            int64     `db:"end_ms"`
}

func NewDbChunk(videoId uuid.UUID, representationId string, sequence int, contentLocation string, startMs int64, endMs int64) *DbChunk {
	return &DbChunk{
		Id:               uuid.New(),
		VideoId:          videoId,
		RepresentationId: representationId,
		Sequence:         sequence,
		ContentLocation:  contentLocation,
		StartMs:          startMs,
		EndMs:            endMs,
	}
}

type ChunksRepository struct {
	db     *sqlx.DB
	logger zerolog.Logger
	tracer trace.Tracer
}

func NewChunksRepository(dependencies *app.Dependencies) *ChunksRepository {
	return &ChunksRepository{
		db:     dependencies.Postgres,
		logger: dependencies.Logger,
		tracer: dependencies.Tracer,
	}
}

func (r *ChunksRepository) Insert(chunk *DbChunk, ctx context.Context) (*DbChunk, error) {
	r.logger.Debug().Str("videoId", chunk.VideoId.String()).Msg("Inserting chunk")

	_, err := r.db.NamedExecContext(ctx, "INSERT INTO chunks (id, video_id, representation_id, sequence, content_location, start_ms, end_ms) VALUES (:id,:video_id, :representation_id, :sequence, :content_location, :start_ms, :end_ms)", chunk)
	if err == nil {
		return chunk, nil
	}

	pgError := err.(*pq.Error)
	if pgError.Code == "23505" {
		return nil, ErrChunkAlreadyExists
	}

	return nil, err
}

func (r *ChunksRepository) GetMany(videoId uuid.UUID, startMs, endMs int64, ctx context.Context) ([]*DbChunk, error) {
	ctx, span := r.tracer.Start(ctx, "chunks.repository.getMany")
	defer span.End()
	r.logger.Debug().Str("videoId", videoId.String()).Msg("Searching chunks")

	var chunks []*DbChunk
	err := r.db.SelectContext(ctx, &chunks, "SELECT * FROM chunks WHERE video_id = $1 AND start_ms >= $2 AND end_ms <= $3", videoId, startMs, endMs)
	if err != nil {
		return nil, errors.Join(err, ErrFailedToGetChunks)
	}

	return chunks, nil
}
