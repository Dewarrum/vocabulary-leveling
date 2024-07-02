package server

import (
	"dewarrum/vocabulary-leveling/internal/subtitles"
	"dewarrum/vocabulary-leveling/internal/videos"
	"net/http"

	"github.com/gofiber/fiber/v2"
)

func (s *Server) VideosUpload(router fiber.Router) {
	router.Post("/videos/upload", func(c *fiber.Ctx) error {
		videoHeader, err := c.FormFile("video")
		if err != nil {
			c.Status(http.StatusBadRequest).JSON(map[string]string{"error": err.Error()})
			return nil
		}

		videoFile, err := videoHeader.Open()
		if err != nil {
			c.Status(http.StatusBadRequest).JSON(map[string]string{"error": err.Error()})
			return nil
		}
		defer videoFile.Close()

		videoName := c.FormValue("videoName")
		if videoName == "" {
			return c.Status(http.StatusBadRequest).JSON(map[string]string{"error": "name is required"})
		}

		subtitlesHeader, err := c.FormFile("subtitles")
		if err != nil {
			c.Status(http.StatusBadRequest).JSON(map[string]string{"error": err.Error()})
			return nil
		}

		subtitlesFile, err := subtitlesHeader.Open()
		if err != nil {
			c.Status(http.StatusBadRequest).JSON(map[string]string{"error": err.Error()})
			return nil
		}
		defer subtitlesFile.Close()

		video := videos.NewDbVideo(videoName)
		_, err = s.Videos.Repository.Insert(video, c.Context())
		if err != nil {
			c.Status(http.StatusInternalServerError).JSON(map[string]string{"error": err.Error()})
			return nil
		}

		err = s.Videos.FileStorage.Upload(video.Id, videoFile, videoHeader.Header.Get("Content-Type"), c.Context())
		if err != nil {
			c.Status(http.StatusInternalServerError).JSON(map[string]string{"error": err.Error()})
			return nil
		}

		exportVideoMessage := videos.NewExportVideoMessage(video.Id)
		err = s.Videos.Messages.Send(exportVideoMessage, c.Context())
		if err != nil {
			c.Status(http.StatusInternalServerError).JSON(map[string]string{"error": err.Error()})
			return nil
		}

		err = s.Subtitles.FileStorage.Upload(video.Id, subtitlesFile, subtitlesHeader.Header.Get("Content-Type"), c.Context())
		if err != nil {
			c.Status(http.StatusInternalServerError).JSON(map[string]string{"error": err.Error()})
			return nil
		}

		exportSubtitlesMessage := subtitles.NewExportSubtitlesMessage(video.Id)
		err = s.Subtitles.Messages.Send(exportSubtitlesMessage, c.Context())
		if err != nil {
			c.Status(http.StatusInternalServerError).JSON(map[string]string{"error": err.Error()})
			return nil
		}

		return c.Status(http.StatusOK).JSON(map[string]string{"videoId": video.Id.String()})
	})
}
