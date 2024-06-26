package videos

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func export(router fiber.Router, messageQueue *MessageQueue) {
	router.Post("/export", func(c *fiber.Ctx) error {
		videoId, err := uuid.Parse(c.Query("videoId"))
		if err != nil {
			return c.Status(http.StatusBadRequest).JSON(map[string]string{"error": err.Error()})
		}

		message := NewExportVideoMessage(videoId)

		err = messageQueue.Send(message, c.Context())
		if err != nil {
			return c.Status(http.StatusInternalServerError).JSON(map[string]string{"error": err.Error()})
		}

		return c.Status(http.StatusOK).JSON(map[string]string{"message": "export started"})
	})
}
