package app

import (
	"errors"
	"os"

	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/gofiber/storage/redis/v3"
)

var (
	ErrRedisUrlIsRequired = errors.New("REDIS_URL is required")
)

func createSessionStore() (*session.Store, error) {
	url := os.Getenv("REDIS_URL")
	if url == "" {
		return nil, ErrRedisUrlIsRequired
	}

	redisStorage := redis.New(redis.Config{
		URL:      url,
		Database: 0,
	})

	sessionStore := session.New(session.Config{
		Storage: redisStorage,
	})

	return sessionStore, nil
}
