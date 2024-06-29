package subtitles

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type DtoSubtitleCue struct {
	Id       uuid.UUID `json:"cueId"`
	VideoId  uuid.UUID `json:"videoId"`
	Sequence int       `json:"sequence"`
	StartMs  int64     `json:"startMs"`
	EndMs    int64     `json:"endMs"`
	Text     string    `json:"text"`
}

func search(router fiber.Router, fullTextSearch *FullTextSearch, subtitleCueRepository *SubtitleCueRepository) {
	router.Get("/search", func(c *fiber.Ctx) error {
		query := c.Query("query")
		if query == "" {
			return c.Status(400).JSON(map[string]string{"error": "query is required"})
		}

		ftsSubtitleCues, err := fullTextSearch.Search(query, c.Context())
		if err != nil {
			return c.Status(500).JSON(map[string]string{"error": err.Error()})
		}

		if len(ftsSubtitleCues) == 0 {
			return c.Status(200).JSON([]*DtoSubtitleCue{})
		}

		subtitleCueIds := make([]uuid.UUID, len(ftsSubtitleCues))
		for i, ftsSubtitleCue := range ftsSubtitleCues {
			subtitleCueIds[i] = ftsSubtitleCue.Id
		}

		dbSubtitleCues, err := subtitleCueRepository.GetManyByIds(subtitleCueIds, c.Context())
		if err != nil {
			return c.Status(500).JSON(map[string]string{"error": err.Error()})
		}

		dtoSubtitleCues := make([]*DtoSubtitleCue, len(dbSubtitleCues))
		for i, dbSubtitleCue := range dbSubtitleCues {
			dtoSubtitleCues[i] = &DtoSubtitleCue{
				Id:       dbSubtitleCue.Id,
				VideoId:  dbSubtitleCue.VideoId,
				Sequence: dbSubtitleCue.Sequence,
				StartMs:  dbSubtitleCue.StartMs,
				EndMs:    dbSubtitleCue.EndMs,
				Text:     dbSubtitleCue.Text,
			}
		}

		return c.Status(200).JSON(dtoSubtitleCues)
	})
}
