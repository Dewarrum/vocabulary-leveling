package inits

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
	ErrInitAlreadyExists = errors.New("init already exists")
	ErrFailedToGetInits  = errors.New("failed to get inits")
)

type DbInit struct {
	Id               uuid.UUID `db:"id"`
	VideoId          uuid.UUID `db:"video_id"`
	RepresentationId string    `db:"representation_id"`
	ContentLocation  string    `db:"content_location"`
}

func NewDbInit(videoId uuid.UUID, representationId string, contentLocation string) *DbInit {
	return &DbInit{
		Id:               uuid.New(),
		VideoId:          videoId,
		RepresentationId: representationId,
		ContentLocation:  contentLocation,
	}
}

type InitsRepository struct {
	db     *sqlx.DB
	logger zerolog.Logger
	tracer trace.Tracer
}

func NewInitsRepository(dependencies *app.Dependencies) *InitsRepository {
	return &InitsRepository{
		db:     dependencies.Postgres,
		logger: dependencies.Logger,
		tracer: dependencies.Tracer,
	}
}

func (r *InitsRepository) Insert(init *DbInit, ctx context.Context) (*DbInit, error) {
	ctx, span := r.tracer.Start(ctx, "inits.repository.insert")
	defer span.End()
	r.logger.Debug().Str("videoId", init.VideoId.String()).Msg("Inserting init")

	_, err := r.db.NamedExecContext(ctx, "INSERT INTO inits (id, video_id, representation_id, content_location) VALUES (:id,:video_id, :representation_id, :content_location)", init)
	if err == nil {
		return init, nil
	}

	pgError := err.(*pq.Error)
	if pgError.Code == "23505" {
		return nil, ErrInitAlreadyExists
	}

	return nil, err
}

func (r *InitsRepository) GetByVideoId(videoId uuid.UUID, ctx context.Context) ([]*DbInit, error) {
	ctx, span := r.tracer.Start(ctx, "inits.repository.getByVideoId")
	defer span.End()
	r.logger.Debug().Str("videoId", videoId.String()).Msg("Searching inits by video id")

	var inits []*DbInit
	err := r.db.SelectContext(ctx, &inits, "SELECT * FROM inits WHERE video_id = $1 ORDER BY representation_id", videoId)
	if err != nil {
		return nil, errors.Join(err, ErrFailedToGetInits)
	}

	return inits, nil
}
