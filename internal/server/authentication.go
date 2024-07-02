package server

import (
	"errors"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/logto-io/go/client"
	"github.com/rs/zerolog"
)

func newLogtoConfig() (*client.LogtoConfig, error) {
	endpoint := os.Getenv("LOGTO_ENDPOINT")
	if endpoint == "" {
		return nil, errors.New("LOGTO_ENDPOINT is required")
	}

	appId := os.Getenv("LOGTO_APP_ID")
	if appId == "" {
		return nil, errors.New("LOGTO_APP_ID is required")
	}

	appSecret := os.Getenv("LOGTO_APP_SECRET")
	if appSecret == "" {
		return nil, errors.New("LOGTO_APP_SECRET is required")
	}

	return &client.LogtoConfig{
		Endpoint:  endpoint,
		AppId:     appId,
		AppSecret: appSecret,
	}, nil
}

type AuthenticationService struct {
	Session     *session.Session
	LogtoClient *client.LogtoClient
	Logger      zerolog.Logger
}

func (s *Server) NewAuthenticationService(c *fiber.Ctx) (*AuthenticationService, error) {
	session, err := s.SessionStore.Get(c)
	if err != nil {
		return nil, err
	}

	logtoClient := client.NewLogtoClient(s.LogtoConfig, &LogtoSessionStorageAdapter{
		FiberSession: session,
		Logger:       s.Logger,
	})

	return &AuthenticationService{
		Session:     session,
		LogtoClient: logtoClient,
		Logger:      s.Logger,
	}, nil
}

func (as *AuthenticationService) Close() {
	as.Session.Save()
}
