package videos

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func upload(router fiber.Router, fileStorage *FileStorage) {
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

		videoId := uuid.New()

		err = fileStorage.Upload(videoId, file, fileHeader.Header.Get("Content-Type"), c.Context())
		if err != nil {
			c.Status(http.StatusInternalServerError).JSON(map[string]string{"error": err.Error()})
			return nil
		}

		return c.Status(http.StatusOK).JSON(map[string]string{"videoId": videoId.String()})
	})
}
