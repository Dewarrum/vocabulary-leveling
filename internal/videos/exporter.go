package videos

import (
	"context"
	"dewarrum/vocabulary-leveling/internal/app"
	"dewarrum/vocabulary-leveling/internal/chunks"
	"dewarrum/vocabulary-leveling/internal/inits"
	"dewarrum/vocabulary-leveling/internal/manifests"
	"dewarrum/vocabulary-leveling/internal/mpd"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"regexp"
	"strconv"

	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"go.opentelemetry.io/otel/trace"
)

var (
	ErrFailedToRun     = errors.New("failed to run")
	chunkStreamPattern = regexp.MustCompile(`stream-(\d{5})\.m4s`)
)

type Exporter struct {
	manifestsRepository *manifests.ManifestsRepository
	initsRepository     *inits.InitsRepository
	chunksRepository    *chunks.ChunksRepository
	fileStorage         *FileStorage
	messageQueue        *MessageQueue
	logger              zerolog.Logger
	tracer              trace.Tracer
}

func NewExporter(dependencies *app.Dependencies) (*Exporter, error) {
	messageQueue, err := NewMessageQueue(dependencies)
	if err != nil {
		return nil, errors.Join(err, errors.New("failed to create message queue"))
	}

	return &Exporter{
		manifestsRepository: manifests.NewManifestsRepository(dependencies),
		initsRepository:     inits.NewInitsRepository(dependencies),
		chunksRepository:    chunks.NewChunksRepository(dependencies),
		fileStorage:         NewFileStorage(dependencies.S3Client, dependencies.S3PresignClient),
		messageQueue:        messageQueue,
		logger:              dependencies.Logger,
		tracer:              dependencies.Tracer,
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
			e.logger.Info().Str("videoId", message.VideoId.String()).Msg("Exporting video")

			err = e.handleMessage(message, context)
			if err != nil {
				e.logger.Error().Str("videoId", message.VideoId.String()).Err(err).Msg("Failed to export video")
				continue
			}

			e.logger.Info().Str("videoId", message.VideoId.String()).Msg("Video exported successfully")
		}
	}()

	return nil
}

func (e *Exporter) handleMessage(message ExportVideoMessage, context context.Context) error {
	directory := fmt.Sprintf("tmp/%s", message.VideoId)
	err := os.MkdirAll(fmt.Sprintf("%s/chunks/0", directory), 0755)
	if err != nil {
		return errors.Join(err, errors.New("failed to create directory"))
	}
	err = os.MkdirAll(fmt.Sprintf("%s/chunks/1", directory), 0755)
	if err != nil {
		return errors.Join(err, errors.New("failed to create directory"))
	}
	err = os.MkdirAll(fmt.Sprintf("%s/inits/0", directory), 0755)
	if err != nil {
		return errors.Join(err, errors.New("failed to create directory"))
	}
	err = os.MkdirAll(fmt.Sprintf("%s/inits/1", directory), 0755)
	if err != nil {
		return errors.Join(err, errors.New("failed to create directory"))
	}
	// defer os.RemoveAll(directory)

	err = e.downloadVideo(message.VideoId, directory, context)
	if err != nil {
		return errors.Join(err, errors.New("failed to download video"))
	}

	err = e.convertToDash(directory)
	if err != nil {
		return errors.Join(err, errors.New("failed to run ffmpeg"))
	}

	err = e.saveContents(message.VideoId, directory, context)
	if err != nil {
		return errors.Join(err, errors.New("failed to upload video"))
	}

	return nil
}

func (e *Exporter) downloadVideo(videoId uuid.UUID, directory string, context context.Context) error {
	e.logger.Info().Str("videoId", videoId.String()).Msg("Start downloading video")

	response, err := e.fileStorage.Download(videoId, context)
	if err != nil {
		return errors.Join(err, errors.New("failed to download video"))
	}

	e.logger.Info().Str("videoId", videoId.String()).Msg("Finished downloading video")

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

	e.logger.Info().Str("videoId", videoId.String()).Msg("Saved downloaded video to file system")

	return nil
}

func (e *Exporter) convertToDash(directory string) error {
	e.logger.Info().Str("videoId", directory).Msg("Running ffmpeg")

	cmd := exec.Command(
		"ffmpeg",
		"-i", fmt.Sprintf("%s/original", directory),
		"-g", "30",
		"-keyint_min", "30",
		"-sc_threshold", "0",
		"-seg_duration", "5",
		"-init_seg_name", "inits/$RepresentationID$/stream.$ext$",
		"-media_seg_name", "chunks/$RepresentationID$/stream-$Number%05d$.$ext$",
		"-f", "dash",
		fmt.Sprintf("%s/manifest.mpd", directory))

	err := cmd.Run()
	if err != nil {
		return errors.Join(err, errors.New("failed to run ffmpeg"))
	}

	e.logger.Info().Str("videoId", directory).Msg("Successfully exported video via ffmpeg")

	return nil
}

func (e *Exporter) saveInitStream(videoId uuid.UUID, representationId string, directory string, ctx context.Context) error {
	e.logger.Info().Str("videoId", videoId.String()).Msg("Start saving chunks")

	file, err := os.Open(fmt.Sprintf("%s/inits/%s/stream.m4s", directory, representationId))
	if err != nil {
		return errors.Join(err, errors.New("failed to open file"))
	}

	contentLocation, err := e.fileStorage.UploadInitStream(videoId, representationId, "stream.m4s", file, ctx)
	if err != nil {
		return errors.Join(err, errors.New("failed to upload chunk stream"))
	}

	init := inits.NewDbInit(videoId, representationId, contentLocation)
	_, err = e.initsRepository.Insert(init, ctx)
	if err != nil {
		return errors.Join(err, errors.New("failed to save init to database"))
	}

	return nil
}

func getChunkStreamNumber(filename string) (int64, error) {
	matches := chunkStreamPattern.FindStringSubmatch(filename)
	if len(matches) != 2 {
		return 0, errors.New("failed to find chunk stream number")
	}

	number, err := strconv.ParseInt(matches[1], 10, 64)
	if err != nil {
		return 0, errors.Join(err, errors.New("failed to parse chunk stream number"))
	}
	return number, nil
}

func (e *Exporter) saveChunkStreams(videoId uuid.UUID, representationId string, directory string, segmentTemplate *mpd.SegmentTemplate, ctx context.Context) error {
	e.logger.Info().Str("videoId", videoId.String()).Msg("Start saving chunks")

	entries, err := os.ReadDir(fmt.Sprintf("%s/chunks/%s", directory, representationId))
	if err != nil {
		return errors.Join(err, errors.New("failed to read directory"))
	}

	segmentInfos, err := segmentTemplate.GetSegmentInfos()
	if err != nil {
		return errors.Join(err, errors.New("failed to get segment infos"))
	}

	if len(segmentInfos) != len(entries) {
		e.logger.Error().Str("videoId", videoId.String()).Int32("segmentCount", int32(len(segmentInfos))).Int32("chunkStreamCount", int32(len(entries))).Msg("Segment count and chunk stream count do not match")
		return errors.New("segment count and chunk stream count do not match")
	}

	for _, entry := range entries {
		file, err := os.Open(fmt.Sprintf("%s/chunks/%s/%s", directory, representationId, entry.Name()))
		if err != nil {
			return errors.Join(err, errors.New("failed to open file"))
		}

		chunkStreamNumber, err := getChunkStreamNumber(entry.Name())
		if err != nil {
			return errors.Join(err, errors.New("failed to get chunk stream number"))
		}

		contentLocation, err := e.fileStorage.UploadChunkStream(videoId, representationId, entry.Name(), file, ctx)
		if err != nil {
			return errors.Join(err, errors.New("failed to upload chunk stream"))
		}

		segmentInfo := segmentInfos[chunkStreamNumber-1]

		chunk := chunks.NewDbChunk(videoId, representationId, int(chunkStreamNumber), contentLocation, segmentInfo.TimestampMs, segmentInfo.TimestampMs+segmentInfo.DurationMs)
		_, err = e.chunksRepository.Insert(chunk, ctx)
		if err != nil {
			return errors.Join(err, errors.New("failed to save chunk to database"))
		}
	}

	return nil
}

func (e *Exporter) saveContents(videoId uuid.UUID, directory string, ctx context.Context) error {
	e.logger.Info().Str("videoId", videoId.String()).Msg("Start uploading video")

	e.logger.Info().Str("mifestPath", fmt.Sprintf("%s/manifest.mpd", directory)).Msg("Opening manifest file")
	manifestFile, err := os.Open(fmt.Sprintf("%s/manifest.mpd", directory))
	if err != nil {
		return errors.Join(err, errors.New("failed to open manifest file"))
	}

	// TODO: don't upload manifest to s3
	err = e.fileStorage.UploadManifest(videoId, manifestFile, ctx)
	if err != nil {
		return errors.Join(err, errors.New("failed to upload manifest"))
	}

	// TODO: don't read manifest file twice
	manifestFile, err = os.Open(fmt.Sprintf("%s/manifest.mpd", directory))
	if err != nil {
		return errors.Join(err, errors.New("failed to open manifest file"))
	}
	manifestBody, err := io.ReadAll(manifestFile)
	if err != nil {
		return errors.Join(err, errors.New("failed to read manifest file"))
	}
	manifest, err := mpd.Parse(manifestBody)
	if err != nil {
		return errors.Join(err, errors.New("failed to parse manifest"))
	}

	err = e.saveManifest(videoId, manifest, ctx)
	if err != nil {
		return errors.Join(err, errors.New("failed to save manifest to database"))
	}

	entries, err := os.ReadDir(fmt.Sprintf("%s/inits", directory))
	if err != nil {
		return errors.Join(err, errors.New("failed to read directory"))
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		representationId := entry.Name()
		err = e.saveInitStream(videoId, representationId, directory, ctx)
		if err != nil {
			return errors.Join(err, errors.New("failed to save init stream"))
		}
	}

	entries, err = os.ReadDir(fmt.Sprintf("%s/chunks", directory))
	if err != nil {
		return errors.Join(err, errors.New("failed to read directory"))
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		representationId := entry.Name()
		representation := manifest.GetRepresentation(representationId)
		if representation == nil {
			continue
		}
		err = e.saveChunkStreams(videoId, representationId, directory, representation.SegmentTemplate, ctx)
		if err != nil {
			return errors.Join(err, errors.New("failed to save chunk stream"))
		}
	}

	e.logger.Info().Str("videoId", videoId.String()).Msg("Finished uploading video")
	return nil
}

func (e *Exporter) saveManifest(videoId uuid.UUID, manifest *mpd.MPD, ctx context.Context) error {
	dbManifest, err := manifests.NewDbManifest(videoId, manifest)
	if err != nil {
		return errors.Join(err, errors.New("failed to create db manifest"))
	}

	_, err = e.manifestsRepository.Insert(dbManifest, ctx)
	if errors.Is(err, manifests.ErrManifestAlreadyExists) {
		return nil
	}

	if err != nil {
		return errors.Join(err, errors.New("failed to insert manifest"))
	}

	return nil
}
