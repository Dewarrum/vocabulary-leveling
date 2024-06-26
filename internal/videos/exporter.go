package videos

import (
	"context"
	"dewarrum/vocabulary-leveling/internal/app"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"

	"github.com/rs/zerolog"
	"go.opentelemetry.io/otel/trace"
)

var (
	ErrFailedToRun = errors.New("failed to run")
)

type Exporter struct {
	fileStorage  *FileStorage
	messageQueue *MessageQueue
	logger       zerolog.Logger
	tracer       trace.Tracer
}

func NewExporter(dependencies *app.Dependencies) (*Exporter, error) {
	messageQueue, err := NewMessageQueue(dependencies.RabbitMqChannel)
	if err != nil {
		return nil, errors.Join(err, errors.New("failed to create message queue"))
	}

	return &Exporter{
		fileStorage:  NewFileStorage(dependencies.S3Client, dependencies.S3PresignClient),
		messageQueue: messageQueue,
		logger:       dependencies.Logger,
		tracer:       dependencies.Tracer,
	}, nil
}

func (e *Exporter) Run(context context.Context) error {
	e.logger.Info().Msg("Starting video exporter")

	messages, err := e.messageQueue.Consume()
	if err != nil {
		e.logger.Fatal().Err(err).Msg("Failed to register a consumer")
		return errors.Join(err, ErrFailedToRun)
	}

	go func() {
		for message := range messages {
			e.logger.Info().Str("videoId", message.VideoId).Msg("Exporting video")

			err = e.handleMessage(message, context)
			if err != nil {
				e.logger.Error().Str("videoId", message.VideoId).Err(err).Msg("Failed to export video")
				continue
			}

			e.logger.Info().Str("videoId", message.VideoId).Msg("Video exported successfully")
		}
	}()

	return nil
}

func (e *Exporter) handleMessage(message ExportVideoMessage, context context.Context) error {
	directory := fmt.Sprintf("tmp/%s", message.VideoId)
	err := os.Mkdir(directory, 0755)
	if err != nil {
		return errors.Join(err, errors.New("failed to create directory"))
	}

	err = e.downloadVideo(message.VideoId, directory, context)
	if err != nil {
		return errors.Join(err, errors.New("failed to download video"))
	}

	err = e.exportVideo(directory)
	if err != nil {
		return errors.Join(err, errors.New("failed to run ffmpeg"))
	}

	err = e.uploadVideo(message.VideoId, directory, context)
	if err != nil {
		return errors.Join(err, errors.New("failed to upload video"))
	}

	err = clearDirectory(directory)
	if err != nil {
		return errors.Join(err, errors.New("failed to clear directory"))
	}

	return nil
}

func (e *Exporter) downloadVideo(videoId, directory string, context context.Context) error {
	e.logger.Info().Str("videoId", videoId).Msg("Start downloading video")

	response, err := e.fileStorage.Download(videoId, context)
	if err != nil {
		return errors.Join(err, errors.New("failed to download video"))
	}

	e.logger.Info().Str("videoId", videoId).Msg("Finished downloading video")

	path := fmt.Sprintf("%s/original", directory)
	fi, err := os.Create(path)
	if err != nil {
		return errors.Join(err, errors.New("failed to create file"))
	}
	defer fi.Close()

	_, err = io.Copy(fi, response.Body)
	if err != nil {
		return errors.Join(err, errors.New("failed to download video"))
	}

	e.logger.Info().Str("videoId", videoId).Msg("Saved downloaded video to file system")

	return nil
}

func (e *Exporter) exportVideo(directory string) error {
	e.logger.Info().Str("videoId", directory).Msg("Running ffmpeg")

	cmd := exec.Command(
		"ffmpeg",
		"-i", fmt.Sprintf("%s/original", directory),
		"-g", "30",
		"-keyint_min", "30",
		"-sc_threshold", "0",
		"-seg_duration", "5",
		"-use_template", "0",
		"-init_seg_name", "output-init-stream$RepresentationID$.$ext$",
		"-media_seg_name", "output-chunk-stream$RepresentationID$-$Number%05d$.$ext$",
		"-f", "dash",
		fmt.Sprintf("%s/manifest.mpd", directory))

	err := cmd.Run()
	if err != nil {
		return errors.Join(err, errors.New("failed to run ffmpeg"))
	}

	e.logger.Info().Str("videoId", directory).Msg("Successfully exported video via ffmpeg")

	return nil
}

func (e *Exporter) uploadVideo(videoId, directory string, context context.Context) error {
	e.logger.Info().Str("videoId", videoId).Msg("Start uploading video")

	entries, err := os.ReadDir(directory)
	if err != nil {
		return errors.Join(err, errors.New("failed to read directory"))
	}

	for _, entry := range entries {
		if entry.IsDir() {
			e.logger.Info().Str("videoId", videoId).Str("file", entry.Name()).Msg("Skipping file")
			continue
		}

		if isChunkStream(entry.Name()) {
			file, err := os.Open(fmt.Sprintf("%s/%s", directory, entry.Name()))
			if err != nil {
				return errors.Join(err, errors.New("failed to open file"))
			}

			err = e.fileStorage.UploadChunkStream(videoId, entry.Name(), file, context)
			if err != nil {
				return errors.Join(err, errors.New("failed to upload chunk stream"))
			}
			continue
		}

		if isInitStream(entry.Name()) {
			file, err := os.Open(fmt.Sprintf("%s/%s", directory, entry.Name()))
			if err != nil {
				return errors.Join(err, errors.New("failed to open file"))
			}

			err = e.fileStorage.UploadInitStream(videoId, entry.Name(), file, context)
			if err != nil {
				return errors.Join(err, errors.New("failed to upload init stream"))
			}
		}

		if entry.Name() == "manifest.mpd" {
			file, err := os.Open(fmt.Sprintf("%s/%s", directory, entry.Name()))
			if err != nil {
				return errors.Join(err, errors.New("failed to open file"))
			}

			err = e.fileStorage.UploadManifest(videoId, file, context)
			if err != nil {
				return errors.Join(err, errors.New("failed to upload manifest"))
			}
		}

		e.logger.Info().Str("videoId", videoId).Str("file", entry.Name()).Msg("Skipping unknown file")
	}

	e.logger.Info().Str("videoId", videoId).Msg("Finished uploading video")
	return nil
}

func clearDirectory(directory string) error {
	err := os.RemoveAll(directory)
	if err != nil {
		return errors.Join(err, errors.New("failed to clear directory"))
	}

	return nil
}
