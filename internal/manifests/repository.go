package manifests

import (
	"context"
	"dewarrum/vocabulary-leveling/internal/app"
	"errors"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"github.com/rs/zerolog"
	"go.opentelemetry.io/otel/trace"
)

var (
	ErrManifestAlreadyExists = errors.New("manifest already exists")
)

type ManifestsRepository struct {
	db     *sqlx.DB
	logger zerolog.Logger
	tracer trace.Tracer
}

func NewManifestRepository(dependencies *app.Dependencies) *ManifestsRepository {
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
