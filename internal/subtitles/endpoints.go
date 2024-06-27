package subtitles

import (
	"context"
	"dewarrum/vocabulary-leveling/internal/app"

	"github.com/gofiber/fiber/v2"
)

func MapEndpoints(app fiber.Router, dependencies *app.Dependencies, ctx context.Context) error {
	fileStorage := NewFileStorage(dependencies.S3Client, dependencies.S3PresignClient)
	subtitleCueRepository := NewSubtitleCueRepository(dependencies)
	fullTextSearch, err := NewFullTextSearch(dependencies.ElasticsearchClient, ctx)
	if err != nil {
		return err
	}
	messageQueue, err := NewMessageQueue(dependencies)
	if err != nil {
		return err
	}

	group := app.Group("/subtitles")
	upload(group, fileStorage)
	export(group, messageQueue, dependencies.Tracer)
	search(group, fullTextSearch, subtitleCueRepository)

	return nil
}
