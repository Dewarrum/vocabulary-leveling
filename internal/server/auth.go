package server

import (
	"fmt"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/valyala/fasthttp/fasthttpadaptor"
)

func (s *Server) SignIn(router fiber.Router) {
	router.Get("/sign-in", func(c *fiber.Ctx) error {
		auth, err := s.NewAuthenticationService(c)
		if err != nil {
			return c.Status(http.StatusInternalServerError).JSON(map[string]string{"error": err.Error()})
		}
		defer auth.Close()

		signInUrl, err := auth.LogtoClient.SignIn(fmt.Sprintf("%s://%s/api/auth/callback", c.Protocol(), c.Hostname()))
		if err != nil {
			return c.Status(http.StatusInternalServerError).JSON(map[string]string{"error": err.Error()})
		}

		return c.Status(http.StatusTemporaryRedirect).Redirect(signInUrl)
	})
}

func (s *Server) SignInCallback(router fiber.Router) {
	router.Get("/callback", func(c *fiber.Ctx) error {
		auth, err := s.NewAuthenticationService(c)
		if err != nil {
			return c.Status(http.StatusInternalServerError).JSON(map[string]string{"error": err.Error()})
		}
		defer auth.Close()

		var r http.Request
		err = fasthttpadaptor.ConvertRequest(c.Context(), &r, true)
		if err != nil {
			s.Logger.Error().Err(err).Msg("Failed to convert request")
			return c.Status(http.StatusInternalServerError).JSON(map[string]string{"error": err.Error()})
		}

		err = auth.LogtoClient.HandleSignInCallback(&r)
		if err != nil {
			s.Logger.Error().Err(err).Msg("Failed to handle sign in callback")
			return c.Status(http.StatusInternalServerError).JSON(map[string]string{"error": err.Error()})
		}

		return c.Status(http.StatusOK).JSON(map[string]string{"message": "Successfully signed in"})
	})
}

func (s *Server) Profile(router fiber.Router) {
	router.Get("/profile", func(c *fiber.Ctx) error {
		auth, err := s.NewAuthenticationService(c)
		if err != nil {
			return c.Status(http.StatusInternalServerError).JSON(map[string]string{"error": err.Error()})
		}
		defer auth.Close()

		idTokenClaims, err := auth.LogtoClient.GetIdTokenClaims()
		if err != nil {
			return c.Status(http.StatusUnauthorized).JSON(map[string]string{"error": err.Error()})
		}

		return c.Status(http.StatusOK).JSON(idTokenClaims)
	})
}
