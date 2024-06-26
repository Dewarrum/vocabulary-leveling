package subtitles

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

func export(router fiber.Router, messageQueue *MessageQueue, tracer trace.Tracer) {
	router.Post("/export", func(c *fiber.Ctx) error {
		ctx, span := tracer.Start(c.Context(), "export", trace.WithAttributes(attribute.String("videoId", c.Query("videoId"))))
		defer span.End()
		videoId, err := uuid.Parse(c.Query("videoId"))
		if err != nil {
			return c.Status(http.StatusBadRequest).JSON(map[string]string{
				"error": err.Error(),
			})
		}

		message := NewExportSubtitlesMessage(videoId)
		err = messageQueue.Send(message, ctx)
		if err != nil {
			return c.Status(http.StatusInternalServerError).JSON(map[string]string{
				"error": err.Error(),
			})
		}

		return c.Status(http.StatusOK).JSON(map[string]string{
			"message": "export started",
		})
	})
}
