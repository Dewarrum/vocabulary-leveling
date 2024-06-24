package subtitles

import (
	"dewarrum/vocabulary-leveling/internal/app"

	"github.com/gofiber/fiber/v2"
)

func MapEndpoints(app fiber.Router, dependencies *app.Dependencies) error {
	fileStorage := NewFileStorage(dependencies.S3Client, dependencies.S3PresignClient)
	messageQueue, err := NewMessageQueue(dependencies.RabbitMqChannel)
	if err != nil {
		return err
	}

	group := app.Group("/subtitles")
	upload(group, fileStorage)
	export(group, messageQueue)

	return nil
}
