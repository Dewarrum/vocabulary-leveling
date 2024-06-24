package subtitles

import (
	"time"

	"github.com/google/uuid"
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
