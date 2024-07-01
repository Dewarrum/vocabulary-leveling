package server

import (
	"bytes"
	"context"
	"dewarrum/vocabulary-leveling/internal/chunks"
	"dewarrum/vocabulary-leveling/internal/inits"
	"dewarrum/vocabulary-leveling/internal/mpd"
	"dewarrum/vocabulary-leveling/internal/utils"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gofiber/fiber/v2"
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

func (s *Server) presignChunks(chunks []*chunks.DbChunk, ctx context.Context) ([]string, error) {
	var presignedUrls []string
	for _, chunk := range chunks {
		presignedChunk, err := s.Videos.FileStorage.PresignObject(chunk.ContentLocation, ctx)
		if err != nil {
			return nil, errors.Join(err, errors.New("failed to presign object"))
		}
		presignedUrls = append(presignedUrls, presignedChunk)
	}
	return presignedUrls, nil
}

func ExtendRange(startMs int64, endMs int64, desiredDuration int64) (s int64, e int64) {
	actualDuration := endMs - startMs
	if actualDuration >= desiredDuration {
		return startMs, endMs
	}

	requiredAddition := desiredDuration - actualDuration
	halfAddition := requiredAddition / 2

	s = startMs - halfAddition
	e = endMs + halfAddition

	if s < 0 {
		e -= s
		s = 0
	}

	return s, e
}

func (s *Server) insertSegmentList(representation *mpd.Representation, chunks []*chunks.DbChunk, init *inits.DbInit, durationMs int64, ctx context.Context) error {
	presignedVideoChunks, err := s.presignChunks(chunks, ctx)
	if err != nil {
		return err
	}

	timescale, err := strconv.ParseInt(representation.SegmentTemplate.Timescale, 10, 64)
	if err != nil {
		return err
	}

	representation.SegmentList = &mpd.SegmentList{
		Timescale:      representation.SegmentTemplate.Timescale,
		Duration:       fmt.Sprintf("%d", durationMs*timescale/1000),
		Initialization: &mpd.Initialization{},
		Segments:       make([]*mpd.Segment, len(chunks)),
	}

	presignedInit, err := s.Videos.FileStorage.PresignObject(init.ContentLocation, ctx)
	if err != nil {
		return err
	}

	representation.SegmentList.Initialization.SourceURL = presignedInit

	err = insertPresignedChunkStreams(representation.SegmentList, presignedVideoChunks)
	if err != nil {
		return err
	}

	representation.SegmentTemplate = nil

	return nil
}

func (s *Server) VideosManifest(router fiber.Router) {
	router.Get("/videos/manifest.mpd", func(c *fiber.Ctx) error {
		subtitleId := c.Query("subtitleId")
		if subtitleId == "" {
			return c.Status(http.StatusBadRequest).JSON(map[string]string{"error": "subtitleId is required"})
		}

		subtitle, err := s.Subtitles.Repository.GetById(subtitleId, c.Context())
		if err != nil {
			return c.Status(http.StatusInternalServerError).JSON(map[string]string{"error": err.Error()})
		}
		videoId := subtitle.VideoId

		dbManifest, err := s.ManifestsRepository.GetByVideoId(videoId, c.Context())
		if err != nil {
			c.Status(http.StatusInternalServerError).JSON(map[string]string{"error": err.Error()})
			return err
		}

		manifestMeta, err := dbManifest.GetMeta()
		if err != nil {
			c.Status(http.StatusInternalServerError).JSON(map[string]string{"error": err.Error()})
			return err
		}

		dbInits, err := s.InitsRepository.GetByVideoId(videoId, c.Context())
		if err != nil {
			c.Status(http.StatusInternalServerError).JSON(map[string]string{"error": err.Error()})
			return err
		}

		chunkDuration, err := manifestMeta.GetChunkDuration()
		s.Logger.Debug().Int64("chunkDuration", chunkDuration).Msg("Chunk duration")
		if err != nil {
			return c.Status(http.StatusInternalServerError).JSON(map[string]string{"error": err.Error()})
		}

		s.Logger.Debug().Int64("startMs", subtitle.StartMs).Int64("endMs", subtitle.EndMs).Msg("Subtitle range")
		subtitleDuration := subtitle.EndMs - subtitle.StartMs
		s.Logger.Debug().Int64("subtitleDuration", subtitleDuration).Msg("Subtitle duration")

		startMs, endMs := ExtendRange(subtitle.StartMs, subtitle.EndMs, max(subtitleDuration, chunkDuration)*2)

		dbChunks, err := s.ChunksRepository.GetMany(videoId, startMs, endMs, c.Context())
		if err != nil {
			return c.Status(http.StatusInternalServerError).JSON(map[string]string{"error": err.Error()})
		}

		dbVideoChunks := utils.Filter(dbChunks, func(chunk *chunks.DbChunk) bool { return chunk.RepresentationId == "0" })
		dbAudioChunks := utils.Filter(dbChunks, func(chunk *chunks.DbChunk) bool { return chunk.RepresentationId == "1" })

		s.insertSegmentList(manifestMeta.GetRepresentation("0"), dbVideoChunks, dbInits[0], 1, c.Context())
		s.insertSegmentList(manifestMeta.GetRepresentation("1"), dbAudioChunks, dbInits[1], 1, c.Context())

		c.Set("Content-Type", "application/dash+xml")
		serialized, err := manifestMeta.Serialize()
		if err != nil {
			return c.Status(http.StatusInternalServerError).JSON(map[string]string{"error": err.Error()})
		}

		serialized = bytes.Replace(serialized, []byte("&amp;"), []byte("&"), -1)
		return c.Status(http.StatusOK).Send(serialized)
	})
}
