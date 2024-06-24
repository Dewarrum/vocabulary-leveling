package videos

import (
	"dewarrum/vocabulary-leveling/internal/mpd"
	"errors"
	"net/http"
	"slices"
	"strings"

	"github.com/gofiber/fiber/v2"
)

func insertPresignedChunkStreams(segmentList *mpd.SegmentList, presignedUrls []string) error {
	for i, segment := range segmentList.Segments {
		ix := slices.IndexFunc(presignedUrls, func(s string) bool {
			return strings.Contains(s, segment.Media)
		})
		if ix == -1 {
			return errors.New("failed to find segment")
		}
		segmentList.Segments[i].Media = presignedUrls[ix]
	}

	return nil
}

func insertPresignedInitStream(segmentList *mpd.SegmentList, presignedUrl string) {
	segmentList.Initialization.SourceURL = presignedUrl
}

func manifest(router fiber.Router, fileStorage *FileStorage) {
	router.Get("/manifest.mpd", func(c *fiber.Ctx) error {
		videoId := c.Query("videoId")

		presignedChunkStreams, err := fileStorage.PresignChunkStreams(videoId, c.Context())
		if err != nil {
			c.Status(http.StatusInternalServerError).JSON(map[string]string{"error": err.Error()})
			return err
		}

		presignedInitStreams, err := fileStorage.PresignInitStreams(videoId, c.Context())
		if err != nil {
			c.Status(http.StatusInternalServerError).JSON(map[string]string{"error": err.Error()})
			return err
		}

		manifest, err := fileStorage.DownloadManifest(videoId, c.Context())
		if err != nil {
			c.Status(http.StatusInternalServerError).JSON(map[string]string{"error": err.Error()})
			return nil
		}

		videoSegmentList := manifest.Periods[0].AdaptationSets[0].Representations[0].SegmentList
		insertPresignedInitStream(videoSegmentList, presignedInitStreams[0])
		err = insertPresignedChunkStreams(videoSegmentList, presignedChunkStreams)
		if err != nil {
			c.Status(http.StatusInternalServerError).JSON(map[string]string{"error": err.Error()})
			return nil
		}

		audioSegmentList := manifest.Periods[0].AdaptationSets[1].Representations[0].SegmentList
		insertPresignedInitStream(audioSegmentList, presignedInitStreams[1])
		err = insertPresignedChunkStreams(audioSegmentList, presignedChunkStreams)
		if err != nil {
			c.Status(http.StatusInternalServerError).JSON(map[string]string{"error": err.Error()})
			return nil
		}

		c.Set("Content-Type", "application/dash+xml")
		serialized, err := manifest.Serialize()
		if err != nil {
			c.Status(http.StatusInternalServerError).JSON(map[string]string{"error": err.Error()})
			return nil
		}
		return c.Status(http.StatusOK).Send(serialized)
	})
}
