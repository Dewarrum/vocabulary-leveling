package subtitles

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type DtoSubtitle struct {
	Id       string    `json:"id"`
	VideoId  uuid.UUID `json:"videoId"`
	Sequence int       `json:"sequence"`
	StartMs  int64     `json:"startMs"`
	EndMs    int64     `json:"endMs"`
	Text     string    `json:"text"`
}

func newDtoSubtitle(subtitle *DbSubtitle) *DtoSubtitle {
	return &DtoSubtitle{
		Id:       subtitle.Id,
		VideoId:  subtitle.VideoId,
		Sequence: subtitle.Sequence,
		StartMs:  subtitle.StartMs,
		EndMs:    subtitle.EndMs,
		Text:     subtitle.Text,
	}
}

func search(router fiber.Router, fullTextSearch *FullTextSearch, subtitleRepository *SubtitlesRepository) {
	router.Get("/search", func(c *fiber.Ctx) error {
		query := c.Query("query")
		if query == "" {
			return c.Status(400).JSON(map[string]string{"error": "query is required"})
		}

		ftsSubtitles, err := fullTextSearch.Search(query, c.Context())
		if err != nil {
			return c.Status(500).JSON(map[string]string{"error": err.Error()})
		}

		if len(ftsSubtitles) == 0 {
			return c.Status(200).JSON([]*DtoSubtitle{})
		}

		subtitleIds := make([]string, len(ftsSubtitles))
		for i, ftsSubtitle := range ftsSubtitles {
			subtitleIds[i] = ftsSubtitle.Id
		}

		dbSubtitles, err := subtitleRepository.GetManyByIds(subtitleIds, c.Context())
		if err != nil {
			return c.Status(500).JSON(map[string]string{"error": err.Error()})
		}

		dtoSubtitles := make([]*DtoSubtitle, len(dbSubtitles))
		for i, dbSubtitle := range dbSubtitles {
			dtoSubtitles[i] = newDtoSubtitle(dbSubtitle)
		}

		return c.Status(200).JSON(dtoSubtitles)
	})
}
