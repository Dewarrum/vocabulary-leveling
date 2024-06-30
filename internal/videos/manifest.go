package videos

import (
	"bytes"
	"dewarrum/vocabulary-leveling/internal/chunks"
	"dewarrum/vocabulary-leveling/internal/inits"
	"dewarrum/vocabulary-leveling/internal/manifests"
	"dewarrum/vocabulary-leveling/internal/mpd"
	"dewarrum/vocabulary-leveling/internal/subtitles"
	"dewarrum/vocabulary-leveling/internal/utils"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
)

func insertPresignedChunkStreams(segmentList *mpd.SegmentList, presignedUrls []string) error {
	for i, presignedUrl := range presignedUrls {
		segment := &mpd.Segment{
			Media: presignedUrl,
		}
		segmentList.Segments[i] = segment
	}

	return nil
}

func insertPresignedInitStream(segmentList *mpd.SegmentList, presignedUrl string) {
	segmentList.Initialization.SourceURL = presignedUrl
}

func manifest(router fiber.Router, fileStorage *FileStorage, subtitleRepository *subtitles.SubtitlesRepository, manifestsRepository *manifests.ManifestsRepository, initsRepository *inits.InitsRepository, chunksRepository *chunks.ChunksRepository, logger zerolog.Logger) {
	router.Get("/manifest.mpd", func(c *fiber.Ctx) error {
		subtitleIdParts := strings.Split(c.Query("subtitleId"), "/")
		if len(subtitleIdParts) != 2 {
			return c.Status(http.StatusBadRequest).JSON(map[string]string{"error": "subtitleId is invalid"})
		}
		videoId, err := uuid.Parse(subtitleIdParts[0])
		if err != nil {
			return c.Status(http.StatusBadRequest).JSON(map[string]string{"error": err.Error()})
		}
		sequence, err := strconv.Atoi(subtitleIdParts[1])
		if err != nil {
			return c.Status(http.StatusBadRequest).JSON(map[string]string{"error": err.Error()})
		}

		subtitles, err := subtitleRepository.GetRange(videoId, sequence-1, sequence+1, c.Context())
		if err != nil {
			return c.Status(http.StatusInternalServerError).JSON(map[string]string{"error": err.Error()})
		}

		dbManifest, err := manifestsRepository.GetByVideoId(videoId, c.Context())
		if err != nil {
			c.Status(http.StatusInternalServerError).JSON(map[string]string{"error": err.Error()})
			return err
		}

		var manifestMeta mpd.MPD
		err = json.Unmarshal(dbManifest.Meta, &manifestMeta)
		if err != nil {
			c.Status(http.StatusInternalServerError).JSON(map[string]string{"error": err.Error()})
			return err
		}

		dbInits, err := initsRepository.GetByVideoId(videoId, c.Context())
		if err != nil {
			c.Status(http.StatusInternalServerError).JSON(map[string]string{"error": err.Error()})
			return err
		}

		presignedInitStreams := make([]string, len(dbInits))
		for i, init := range dbInits {
			presignedInit, err := fileStorage.PresignObject(init.ContentLocation, c.Context())
			if err != nil {
				c.Status(http.StatusInternalServerError).JSON(map[string]string{"error": err.Error()})
				return err
			}
			presignedInitStreams[i] = presignedInit
		}

		dbChunks, err := chunksRepository.GetMany(videoId, subtitles[0].StartMs-2000, subtitles[len(subtitles)-1].EndMs+2000, c.Context())
		if err != nil {
			c.Status(http.StatusInternalServerError).JSON(map[string]string{"error": err.Error()})
			return err
		}
		logger.Debug().Any("chunks", dbChunks).Msg("Found chunks")

		dbVideoChunks := utils.Filter(dbChunks, func(chunk *chunks.DbChunk) bool {
			return chunk.RepresentationId == "0"
		})
		presignedVideoChunks := make([]string, len(dbVideoChunks))
		for i, chunk := range dbVideoChunks {
			presignedChunk, err := fileStorage.PresignObject(chunk.ContentLocation, c.Context())
			if err != nil {
				c.Status(http.StatusInternalServerError).JSON(map[string]string{"error": err.Error()})
				return err
			}
			presignedVideoChunks[i] = presignedChunk
		}
		videoRepresentation := manifestMeta.GetRepresentation("0")
		var videoRepresentationDuration int64
		for _, chunk := range dbVideoChunks {
			videoRepresentationDuration += chunk.EndMs - chunk.StartMs
		}
		videoRepresentation.SegmentList = &mpd.SegmentList{
			Timescale:      videoRepresentation.SegmentTemplate.Timescale,
			Duration:       fmt.Sprintf("%d", videoRepresentationDuration),
			Initialization: &mpd.Initialization{},
			Segments:       make([]*mpd.Segment, len(dbVideoChunks)),
		}
		insertPresignedInitStream(videoRepresentation.SegmentList, presignedInitStreams[0])
		err = insertPresignedChunkStreams(videoRepresentation.SegmentList, presignedVideoChunks)
		if err != nil {
			c.Status(http.StatusInternalServerError).JSON(map[string]string{"error": err.Error()})
			return nil
		}

		dbAudioChunks := utils.Filter(dbChunks, func(chunk *chunks.DbChunk) bool {
			return chunk.RepresentationId == "1"
		})
		presignedAudioChunks := make([]string, len(dbAudioChunks))
		for i, chunk := range dbAudioChunks {
			presignedChunk, err := fileStorage.PresignObject(chunk.ContentLocation, c.Context())
			if err != nil {
				c.Status(http.StatusInternalServerError).JSON(map[string]string{"error": err.Error()})
				return err
			}
			presignedAudioChunks[i] = presignedChunk
		}
		audioRepresentation := manifestMeta.GetRepresentation("1")
		var audioRepresentationDuration int64
		for _, chunk := range dbAudioChunks {
			audioRepresentationDuration += chunk.EndMs - chunk.StartMs
		}
		audioRepresentation.SegmentList = &mpd.SegmentList{
			Timescale:      audioRepresentation.SegmentTemplate.Timescale,
			Duration:       fmt.Sprintf("%d", audioRepresentationDuration),
			Initialization: &mpd.Initialization{},
			Segments:       make([]*mpd.Segment, len(dbAudioChunks)),
		}
		insertPresignedInitStream(audioRepresentation.SegmentList, presignedInitStreams[1])
		err = insertPresignedChunkStreams(audioRepresentation.SegmentList, presignedAudioChunks)
		if err != nil {
			c.Status(http.StatusInternalServerError).JSON(map[string]string{"error": err.Error()})
			return nil
		}

		videoRepresentation.SegmentTemplate = nil
		audioRepresentation.SegmentTemplate = nil
		c.Set("Content-Type", "application/dash+xml")
		serialized, err := manifestMeta.Serialize()
		if err != nil {
			c.Status(http.StatusInternalServerError).JSON(map[string]string{"error": err.Error()})
			return nil
		}
		serialized = bytes.Replace(serialized, []byte("&amp;"), []byte("&"), -1)
		return c.Status(http.StatusOK).Send(serialized)
	})
}
