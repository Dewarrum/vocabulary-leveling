package server

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
)

func (s *Server) SubtitlesUpload(router fiber.Router) {
	router.Post("/upload", func(c *fiber.Ctx) error {

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

		videoId := c.FormValue("videoId")

		err = s.Subtitles.FileStorage.Upload(videoId, file, fileHeader.Header.Get("Content-Type"), c.Context())
		if err != nil {
			c.Status(http.StatusInternalServerError).JSON(map[string]string{"error": err.Error()})
			return nil
		}

		return c.Status(http.StatusOK).JSON(map[string]string{"message": "Subtitles uploaded successfully"})
	})
}
