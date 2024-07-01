package server

import (
	"dewarrum/vocabulary-leveling/internal/videos"
	"net/http"

	"github.com/gofiber/fiber/v2"
)

func (s *Server) VideosUpload(router fiber.Router) {
	router.Post("/videos/upload", func(c *fiber.Ctx) error {
		fileHeader, err := c.FormFile("file")
		if err != nil {
			c.Status(http.StatusBadRequest).JSON(map[string]string{"error": err.Error()})
			return nil
		}

		file, err := fileHeader.Open()
		if err != nil {
			c.Status(http.StatusBadRequest).JSON(map[string]string{"error": err.Error()})
			return nil
		}
		defer file.Close()

		fileName := c.FormValue("name")
		if fileName == "" {
			return c.Status(http.StatusBadRequest).JSON(map[string]string{"error": "name is required"})
		}

		video := videos.NewDbVideo(fileName)
		_, err = s.Videos.Repository.Insert(video, c.Context())
		if err != nil {
			c.Status(http.StatusInternalServerError).JSON(map[string]string{"error": err.Error()})
			return nil
		}

		err = s.Videos.FileStorage.Upload(video.Id, file, fileHeader.Header.Get("Content-Type"), c.Context())
		if err != nil {
			c.Status(http.StatusInternalServerError).JSON(map[string]string{"error": err.Error()})
			return nil
		}

		return c.Status(http.StatusOK).JSON(map[string]string{"videoId": video.Id.String()})
	})
}
