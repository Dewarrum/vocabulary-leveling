package main

import (
	"context"
	"dewarrum/vocabulary-leveling/internal/app"
	"dewarrum/vocabulary-leveling/internal/subtitles"
	"dewarrum/vocabulary-leveling/internal/videos"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()

	dependencies, err := app.NewDependencies()
	if err != nil {
		panic(err)
	}

	videoExporter, err := videos.NewExporter(dependencies)
	if err != nil {
		panic(err)
	}
	go videoExporter.Run(context.Background())

	subtitlesExporter, err := subtitles.NewExporter(dependencies)
	if err != nil {
		panic(err)
	}
	go subtitlesExporter.Run(context.Background())

	app := fiber.New(fiber.Config{
		BodyLimit: 500 * 1024 * 1024,
	})
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello World")
	})

	api := app.Group("/api")
	api.Use(cors.New(cors.Config{
		AllowOrigins: "*",
	}))

	videos.MapEndpoints(api, dependencies)
	subtitles.MapEndpoints(api, dependencies)

	app.Listen(":3000")
}
