package videos

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
)

func export(router fiber.Router, messageQueue *MessageQueue) {
	router.Post("/export", func(c *fiber.Ctx) error {
		videoId := c.Query("videoId")
		if videoId == "" {
			return c.Status(http.StatusBadRequest).JSON(map[string]string{"error": "videoId is required"})
		}

		message := NewExportVideoMessage(videoId)

		err := messageQueue.Send(message, c.Context())
		if err != nil {
			return c.Status(http.StatusInternalServerError).JSON(map[string]string{"error": err.Error()})
		}

		return c.Status(http.StatusOK).JSON(map[string]string{"message": "export started"})
	})
}
