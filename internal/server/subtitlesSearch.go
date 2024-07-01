package server

import (
	"dewarrum/vocabulary-leveling/internal/subtitles"
	"dewarrum/vocabulary-leveling/internal/videos"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func (s *Server) SubtitlesSearch(router fiber.Router) {
	router.Get("/subtitles/search", func(c *fiber.Ctx) error {
		query := c.Query("query")
		if query == "" {
			return c.Status(400).JSON(map[string]string{"error": "query is required"})
		}

		ftsSubtitles, err := s.Subtitles.FullTextSearch.Search(query, c.Context())
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

		dbSubtitles, err := s.Subtitles.Repository.GetManyByIds(subtitleIds, c.Context())
		if err != nil {
			return c.Status(500).JSON(map[string]string{"error": err.Error()})
		}

		videoIds := getVideosIds(dbSubtitles)

		videos, err := s.Videos.Repository.GetManyByIds(videoIds, c.Context())
		if err != nil {
			return c.Status(500).JSON(map[string]string{"error": err.Error()})
		}

		dtoSubtitles := mapToDto(dbSubtitles, getVideoMap(videos))

		return c.Status(200).JSON(dtoSubtitles)
	})
}

func getVideosIds(subtitles []*subtitles.DbSubtitle) []uuid.UUID {
	set := make(map[uuid.UUID]bool)
	videoIds := make([]uuid.UUID, len(subtitles))
	for _, subtitle := range subtitles {
		if _, ok := set[subtitle.VideoId]; !ok {
			set[subtitle.VideoId] = true
			videoIds = append(videoIds, subtitle.VideoId)
		}
	}

	return videoIds
}

func getVideoMap(dbVideos []*videos.DbVideo) map[uuid.UUID]*videos.DbVideo {
	videoMap := make(map[uuid.UUID]*videos.DbVideo)
	for _, video := range dbVideos {
		videoMap[video.Id] = video
	}
	return videoMap
}

func mapToDto(subtitles []*subtitles.DbSubtitle, dbVideoMap map[uuid.UUID]*videos.DbVideo) []*DtoSubtitle {
	dtoSubtitles := make([]*DtoSubtitle, len(subtitles))
	for i, subtitle := range subtitles {
		dbVideo, ok := dbVideoMap[subtitle.VideoId]
		if !ok {
			continue
		}

		dtoSubtitles[i] = &DtoSubtitle{
			Id:        subtitle.Id,
			VideoName: dbVideo.Name,
			StartMs:   subtitle.StartMs,
			EndMs:     subtitle.EndMs,
			Text:      subtitle.Text,
		}
	}
	return dtoSubtitles
}
