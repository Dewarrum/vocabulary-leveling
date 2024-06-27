package manifests

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
	ErrManifestAlreadyExists = errors.New("manifest already exists")
	ErrFailedToGetManifest   = errors.New("failed to get manifest")
)

type ManifestsRepository struct {
	db     *sqlx.DB
	logger zerolog.Logger
	tracer trace.Tracer
}

func NewManifestsRepository(dependencies *app.Dependencies) *ManifestsRepository {
	return &ManifestsRepository{
		db:     dependencies.Postgres,
		logger: dependencies.Logger,
		tracer: dependencies.Tracer,
	}
}

func (r *ManifestsRepository) Insert(manifest *DbManifest, ctx context.Context) (*DbManifest, error) {
	r.logger.Debug().Str("videoId", manifest.VideoId.String()).Msg("Inserting manifest")

	_, err := r.db.NamedExecContext(ctx, "INSERT INTO manifests (id, video_id, meta) VALUES (:id,:video_id, :meta)", manifest)
	if err == nil {
		return manifest, nil
	}

	pgError := err.(*pq.Error)
	if pgError.Code == "23505" {
		return nil, ErrManifestAlreadyExists
	}

	return nil, err
}

func (r *ManifestsRepository) GetByVideoId(videoId uuid.UUID, ctx context.Context) (*DbManifest, error) {
	ctx, span := r.tracer.Start(ctx, "manifests.repository.getByVideoId")
	defer span.End()
	r.logger.Debug().Str("videoId", videoId.String()).Msg("Searching manifest by video id")

	var manifest DbManifest
	err := r.db.GetContext(ctx, &manifest, "SELECT * FROM manifests WHERE video_id = $1 LIMIT 1", videoId)
	if err != nil {
		return nil, errors.Join(err, ErrFailedToGetManifest)
	}

	return &manifest, nil
}
