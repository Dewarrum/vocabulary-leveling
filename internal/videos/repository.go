package videos

import (
	"context"
	"dewarrum/vocabulary-leveling/internal/app"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog"
	"go.opentelemetry.io/otel/trace"
)

type DbVideo struct {
	Id        uuid.UUID `db:"id"`
	Name      string    `db:"name"`
	CreatedAt time.Time `db:"created_at"`
}

func NewDbVideo(name string) *DbVideo {
	return &DbVideo{
		Id:        uuid.New(),
		Name:      name,
		CreatedAt: time.Now().In(time.UTC),
	}
}

type VideosRepository struct {
	db     *sqlx.DB
	logger zerolog.Logger
	tracer trace.Tracer
}

func NewVideosRepository(dependencies *app.Dependencies) *VideosRepository {
	return &VideosRepository{
		db:     dependencies.Postgres,
		logger: dependencies.Logger,
		tracer: dependencies.Tracer,
	}
}

func (r *VideosRepository) Insert(video *DbVideo, ctx context.Context) (*DbVideo, error) {
	r.logger.Debug().Str("videoId", video.Id.String()).Msg("Inserting video")

	_, err := r.db.NamedExecContext(ctx, "INSERT INTO videos (id, name, created_at) VALUES (:id,:name, :created_at)", video)
	if err == nil {
		return video, nil
	}

	return nil, err
}

func (r *VideosRepository) GetManyByIds(ids []uuid.UUID, ctx context.Context) ([]*DbVideo, error) {
	r.logger.Debug().Msg("Searching videos by ids")

	query, args, err := sqlx.In("SELECT id, name, created_at FROM videos WHERE id IN (?)", ids)
	if err != nil {
		return nil, err
	}

	query = r.db.Rebind(query)
	rows, err := r.db.QueryxContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var videos []*DbVideo
	for rows.Next() {
		var video DbVideo
		err := rows.StructScan(&video)
		if err != nil {
			return nil, err
		}
		videos = append(videos, &video)
	}

	return videos, nil
}
