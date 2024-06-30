package subtitles

import (
	"context"
	"dewarrum/vocabulary-leveling/internal/app"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog"
	"go.opentelemetry.io/otel/trace"
)

var (
	ErrFailedToInsertSubtitle  = errors.New("failed to insert subtitle")
	ErrFailedToGetAffectedRows = errors.New("failed to get affected rows")
	ErrFailedToGetSubtitle     = errors.New("failed to get subtitle")
)

type DbSubtitle struct {
	Id        string    `db:"id"`
	VideoId   uuid.UUID `db:"video_id"`
	Sequence  int       `db:"sequence"`
	StartMs   int64     `db:"start_ms"`
	EndMs     int64     `db:"end_ms"`
	Text      string    `db:"text"`
	CreatedAt time.Time `db:"created_at"`
}

func newDbSubtitle(videoId uuid.UUID, text string, sequence int, startMs int64, endMs int64) *DbSubtitle {
	return &DbSubtitle{
		Id:        fmt.Sprintf("%s/%d", videoId, sequence),
		VideoId:   videoId,
		Sequence:  sequence,
		StartMs:   startMs,
		EndMs:     endMs,
		Text:      text,
		CreatedAt: time.Now().In(time.UTC),
	}
}

type SubtitlesRepository struct {
	db     *sqlx.DB
	logger zerolog.Logger
	tracer trace.Tracer
}

func (r *SubtitlesRepository) Insert(subtitle *DbSubtitle, context context.Context) (*DbSubtitle, error) {
	r.logger.Debug().Str("videoId", subtitle.VideoId.String()).Int32("sequence", int32(subtitle.Sequence)).Msg("Inserting subtitle")

	result, err := r.db.NamedExecContext(context, "INSERT INTO subtitles (id, video_id, sequence, start_ms, end_ms, text, created_at) VALUES (:id,:video_id, :sequence, :start_ms, :end_ms, :text, :created_at) ON CONFLICT (video_id, sequence) DO NOTHING", subtitle)
	if err != nil {
		return nil, errors.Join(err, ErrFailedToInsertSubtitle)
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return nil, errors.Join(err, ErrFailedToGetAffectedRows)
	}

	if affected == 1 {
		return subtitle, nil
	}

	rows, err := r.db.QueryxContext(context, "SELECT * FROM subtitles WHERE video_id = $1 AND sequence = $2", subtitle.VideoId, subtitle.Sequence)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var inserted DbSubtitle
	rows.Next()
	err = rows.StructScan(&inserted)
	if err != nil {
		return nil, errors.Join(err, ErrFailedToGetSubtitle)
	}
	return &inserted, nil
}

func (r *SubtitlesRepository) GetManyByIds(ids []string, context context.Context) ([]*DbSubtitle, error) {
	r.logger.Debug().Msg("Searching subtitles by ids")

	query, args, err := sqlx.In("SELECT id, video_id, sequence, start_ms, end_ms, text, created_at FROM subtitles WHERE id IN (?)", ids)
	if err != nil {
		return nil, err
	}

	query = r.db.Rebind(query)
	rows, err := r.db.QueryxContext(context, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var subtitles []*DbSubtitle
	for rows.Next() {
		var subtitle DbSubtitle
		err := rows.StructScan(&subtitle)
		if err != nil {
			return nil, err
		}
		subtitles = append(subtitles, &subtitle)
	}

	return subtitles, nil
}

func (r *SubtitlesRepository) GetById(id string, ctx context.Context) (*DbSubtitle, error) {
	ctx, span := r.tracer.Start(ctx, "subtitles.repository.getById")
	defer span.End()
	r.logger.Debug().Str("id", id).Msg("Searching subtitle by id")

	var subtitle DbSubtitle
	err := sqlx.GetContext(ctx, r.db, &subtitle, "SELECT * FROM subtitles WHERE id = $1 LIMIT 1", id)
	if err != nil {
		return nil, errors.Join(err, ErrFailedToGetSubtitle)
	}

	return &subtitle, nil
}

func NewSubtitlesRepository(dependencies *app.Dependencies) *SubtitlesRepository {
	return &SubtitlesRepository{
		db:     dependencies.Postgres,
		logger: dependencies.Logger,
		tracer: dependencies.Tracer,
	}
}
