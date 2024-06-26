package videos

import (
	"dewarrum/vocabulary-leveling/internal/app"
	"dewarrum/vocabulary-leveling/internal/subtitles"

	"github.com/gofiber/fiber/v2"
)

func MapEndpoints(app fiber.Router, dependencies *app.Dependencies) error {
	subtitleCueRepository := subtitles.NewSubtitleCueRepository(dependencies.Postgres, dependencies.Logger)
	fileStorage := NewFileStorage(dependencies.S3Client, dependencies.S3PresignClient)
	messageQueue, err := NewMessageQueue(dependencies)
	if err != nil {
		return err
	}

	group := app.Group("/videos")
	export(group, messageQueue)
	upload(group, fileStorage)
	manifest(group, fileStorage, subtitleCueRepository)

	return nil
}
