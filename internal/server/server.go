package server

import (
	"context"
	"dewarrum/vocabulary-leveling/internal/app"
	"dewarrum/vocabulary-leveling/internal/chunks"
	"dewarrum/vocabulary-leveling/internal/inits"
	"dewarrum/vocabulary-leveling/internal/manifests"
	"dewarrum/vocabulary-leveling/internal/subtitles"
	"dewarrum/vocabulary-leveling/internal/videos"

	"github.com/rs/zerolog"
	"go.opentelemetry.io/otel/trace"
)

type Server struct {
	Videos    *VideoContext
	Subtitles *SubtitleContext

	ChunksRepository    *chunks.ChunksRepository
	InitsRepository     *inits.InitsRepository
	ManifestsRepository *manifests.ManifestsRepository

	Logger zerolog.Logger
	Tracer trace.Tracer
}

type SubtitleContext struct {
	Repository     *subtitles.SubtitlesRepository
	Messages       *subtitles.MessageQueue
	FullTextSearch *subtitles.FullTextSearch
	FileStorage    *subtitles.FileStorage
}

type VideoContext struct {
	Repository  *videos.VideosRepository
	Messages    *videos.MessageQueue
	FileStorage *videos.FileStorage
}

func NewServer(dependencies *app.Dependencies, ctx context.Context) (*Server, error) {
	subtitleContext, err := newSubtitleContext(dependencies, ctx)
	if err != nil {
		return nil, err
	}

	videoContext, err := newVideoContext(dependencies)
	if err != nil {
		return nil, err
	}

	return &Server{
		Videos:              videoContext,
		Subtitles:           subtitleContext,
		ChunksRepository:    chunks.NewChunksRepository(dependencies),
		InitsRepository:     inits.NewInitsRepository(dependencies),
		ManifestsRepository: manifests.NewManifestsRepository(dependencies),
		Logger:              dependencies.Logger,
		Tracer:              dependencies.Tracer,
	}, nil
}

func newSubtitleContext(dependencies *app.Dependencies, ctx context.Context) (*SubtitleContext, error) {
	subtitlesMessages, err := subtitles.NewMessageQueue(dependencies)
	if err != nil {
		return nil, err
	}

	subtitlesFullTextSearch, err := subtitles.NewFullTextSearch(dependencies.ElasticsearchClient, ctx)
	if err != nil {
		return nil, err
	}

	return &SubtitleContext{
		Repository:     subtitles.NewSubtitlesRepository(dependencies),
		Messages:       subtitlesMessages,
		FullTextSearch: subtitlesFullTextSearch,
		FileStorage:    subtitles.NewFileStorage(dependencies.S3Client, dependencies.S3PresignClient),
	}, nil
}

func newVideoContext(dependencies *app.Dependencies) (*VideoContext, error) {
	videosMessages, err := videos.NewMessageQueue(dependencies)
	if err != nil {
		return nil, err
	}

	return &VideoContext{
		Repository:  videos.NewVideosRepository(dependencies),
		Messages:    videosMessages,
		FileStorage: videos.NewFileStorage(dependencies.S3Client, dependencies.S3PresignClient),
	}, nil
}
