package server

import (
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/rs/zerolog"
)

type LogtoSessionStorageAdapter struct {
	FiberSession *session.Session
	Logger       zerolog.Logger
}

func (s *LogtoSessionStorageAdapter) GetItem(key string) string {
	item := s.FiberSession.Get(key)
	if item == nil {
		return ""
	}

	return item.(string)
}

func (s *LogtoSessionStorageAdapter) SetItem(key string, value string) {
	s.FiberSession.Set(key, value)
}
