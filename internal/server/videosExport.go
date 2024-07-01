package server

import (
	"dewarrum/vocabulary-leveling/internal/videos"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func (s *Server) VideosExport(router fiber.Router) {
	router.Post("/videos/export", func(c *fiber.Ctx) error {
		videoId, err := uuid.Parse(c.Query("videoId"))
		if err != nil {
			return c.Status(http.StatusBadRequest).JSON(map[string]string{"error": err.Error()})
		}

		message := videos.NewExportVideoMessage(videoId)

		err = s.Videos.Messages.Send(message, c.Context())
		if err != nil {
			return c.Status(http.StatusInternalServerError).JSON(map[string]string{"error": err.Error()})
		}

		return c.Status(http.StatusOK).JSON(map[string]string{"message": "export started"})
	})
}
