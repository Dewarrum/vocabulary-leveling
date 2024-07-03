package main

import (
	"context"
	"dewarrum/vocabulary-leveling/internal/app"
	"dewarrum/vocabulary-leveling/internal/server"
	"dewarrum/vocabulary-leveling/internal/subtitles"
	"dewarrum/vocabulary-leveling/internal/videos"
	"fmt"
	"os"
	"os/signal"

	"github.com/gofiber/contrib/fiberzerolog"
	"github.com/gofiber/contrib/otelfiber"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	godotenv.Load(".env")
	godotenv.Load(".env.secret")

	dependencies, err := app.NewDependencies(ctx)
	if err != nil {
		dependencies.Logger.Fatal().Err(err).Msg("Failed to create dependencies")
		panic(err)
	}
	defer dependencies.Close(ctx)

	err = runMigrations(dependencies.Postgres)
	if err != nil {
		dependencies.Logger.Fatal().Err(err).Msg("Failed to run migrations")
		panic(err)
	}

	srv, err := server.NewServer(dependencies, ctx)
	if err != nil {
		dependencies.Logger.Fatal().Err(err).Msg("Failed to create server")
		panic(err)
	}

	videoExporter, err := videos.NewExporter(dependencies)
	if err != nil {
		dependencies.Logger.Fatal().Err(err).Msg("Failed to create video exporter")
		panic(err)
	}
	err = videoExporter.Run(ctx)
	if err != nil {
		dependencies.Logger.Fatal().Err(err).Msg("Failed to run video exporter")
		panic(err)
	}

	subtitlesExporter, err := subtitles.NewExporter(dependencies, ctx)
	if err != nil {
		dependencies.Logger.Fatal().Err(err).Msg("Failed to create subtitles exporter")
		panic(err)
	}
	err = subtitlesExporter.Run(ctx)
	if err != nil {
		dependencies.Logger.Fatal().Err(err).Msg("Failed to run subtitles exporter")
		panic(err)
	}

	app := fiber.New(fiber.Config{
		BodyLimit: 500 * 1024 * 1024,
	})
	app.Use(fiberzerolog.New(fiberzerolog.Config{
		Logger: &dependencies.Logger,
	}))
	app.Use(otelfiber.Middleware())

	api := app.Group("/api", srv.RequireAuthenticationMiddleware())
	api.Use(cors.New(cors.Config{
		AllowOrigins: "*",
	}))

	srv.VideosManifest(api)
	srv.SubtitlesSearch(api)

	adminApi := api.Group("/admin", srv.RequireAuthorizationMiddleware("Admin"))
	srv.VideosUpload(adminApi)

	authApi := app.Group("/auth")

	srv.Profile(authApi)
	srv.SignIn(authApi)
	srv.SignInCallback(authApi)
	srv.SignOut(authApi)

	if os.Getenv("ENVIRONMENT") != "development" {
		app.Static("/", "./web/build")
		app.Get("/*", func(c *fiber.Ctx) error {
			return c.SendFile("./web/build/index.html")
		})
	}

	port := os.Getenv("PORT")
	if port == "" {
		panic("PORT is required")
	}

	if err := app.Listen(fmt.Sprintf(":%s", port)); err != nil {
		dependencies.Logger.Fatal().Err(err).Msg("Failed to start server")
	}
}

func runMigrations(db *sqlx.DB) error {
	driver, err := postgres.WithInstance(db.DB, &postgres.Config{})
	if err != nil {
		return err
	}

	m, err := migrate.NewWithDatabaseInstance("file://db/migrations", "vocabulary-leveling", driver)
	if err != nil {
		return err
	}

	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
		return err
	}

	return nil
}
