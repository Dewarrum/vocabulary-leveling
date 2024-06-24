package subtitles

import "github.com/jmoiron/sqlx"

type SubtitleCueRepository struct {
	db *sqlx.DB
}

func (r *SubtitleCueRepository) Insert(subtitleCue *DbSubtitleCue) error {
	_, err := r.db.NamedExec("INSERT INTO subtitle_cues (id, video_id, sequence, start_ms, end_ms, text, created_at) VALUES (:id, :video_id, :sequence, :start_ms, :end_ms, :text, :created_at)", subtitleCue)
	if err != nil {
		return err
	}

	return nil
}

func newSubtitleCueRepository(db *sqlx.DB) *SubtitleCueRepository {
	return &SubtitleCueRepository{
		db: db,
	}
}
