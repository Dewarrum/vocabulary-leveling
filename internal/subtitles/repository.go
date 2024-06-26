package subtitles

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog"
)

var (
	ErrSubtitleCueInsertFailed = errors.New("failed to insert subtitle cue")
	ErrFailedToGetAffectedRows = errors.New("failed to get affected rows")
	ErrFailedToGetSubtitleCue  = errors.New("failed to get subtitle cue")
)

type DbSubtitleCue struct {
	Id        uuid.UUID `db:"id"`
	VideoId   uuid.UUID `db:"video_id"`
	Sequence  int       `db:"sequence"`
	StartMs   int64     `db:"start_ms"`
	EndMs     int64     `db:"end_ms"`
	Text      string    `db:"text"`
	CreatedAt time.Time `db:"created_at"`
}

func newDbSubtitleCue(videoId uuid.UUID, text string, sequence int, startMs int64, endMs int64) *DbSubtitleCue {
	return &DbSubtitleCue{
		Id:        uuid.New(),
		VideoId:   videoId,
		Sequence:  sequence,
		StartMs:   startMs,
		EndMs:     endMs,
		Text:      text,
		CreatedAt: time.Now().In(time.UTC),
	}
}

type SubtitleCueRepository struct {
	db     *sqlx.DB
	logger zerolog.Logger
}

func (r *SubtitleCueRepository) Insert(subtitleCue *DbSubtitleCue, context context.Context) (*DbSubtitleCue, error) {
	r.logger.Debug().Str("videoId", subtitleCue.VideoId.String()).Int32("sequence", int32(subtitleCue.Sequence)).Msg("Inserting subtitle cue")

	result, err := r.db.NamedExecContext(context, "INSERT INTO subtitle_cues (id, video_id, sequence, start_ms, end_ms, text, created_at) VALUES (:id,:video_id, :sequence, :start_ms, :end_ms, :text, :created_at) ON CONFLICT (video_id, sequence) DO NOTHING", subtitleCue)
	if err != nil {
		return nil, errors.Join(err, ErrSubtitleCueInsertFailed)
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return nil, errors.Join(err, ErrFailedToGetAffectedRows)
	}

	if affected == 1 {
		return subtitleCue, nil
	}

	rows, err := r.db.QueryxContext(context, "SELECT * FROM subtitle_cues WHERE video_id = $1 AND sequence = $2", subtitleCue.VideoId, subtitleCue.Sequence)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var inserted DbSubtitleCue
	rows.Next()
	err = rows.StructScan(&inserted)
	if err != nil {
		return nil, errors.Join(err, ErrFailedToGetSubtitleCue)
	}
	return &inserted, nil
}

func (r *SubtitleCueRepository) GetManyByIds(ids []uuid.UUID, context context.Context) ([]*DbSubtitleCue, error) {
	r.logger.Debug().Msg("Searching subtitle cues by ids")

	query, args, err := sqlx.In("SELECT id, video_id, sequence, start_ms, end_ms, text, created_at FROM subtitle_cues WHERE id IN (?)", ids)
	if err != nil {
		return nil, err
	}

	query = r.db.Rebind(query)
	rows, err := r.db.QueryxContext(context, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var subtitleCues []*DbSubtitleCue
	for rows.Next() {
		var subtitleCue DbSubtitleCue
		err := rows.StructScan(&subtitleCue)
		if err != nil {
			return nil, err
		}
		subtitleCues = append(subtitleCues, &subtitleCue)
	}

	return subtitleCues, nil
}

func newSubtitleCueRepository(db *sqlx.DB, logger zerolog.Logger) *SubtitleCueRepository {
	return &SubtitleCueRepository{
		db:     db,
		logger: logger,
	}
}
