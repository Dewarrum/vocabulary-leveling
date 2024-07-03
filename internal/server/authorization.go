package server

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
)

func (s *Server) RequireAuthorizationMiddleware(requiredRole string) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		auth, err := s.NewAuthenticationService(c)
		if err != nil {
			return c.Status(http.StatusInternalServerError).JSON(map[string]string{"error": err.Error()})
		}

		claims, err := auth.LogtoClient.GetIdTokenClaims()
		if err != nil {
			return c.Status(http.StatusUnauthorized).JSON(map[string]string{"error": err.Error()})
		}

		for _, role := range claims.Roles {
			if role == requiredRole {
				return c.Next()
			}
		}

		s.Logger.Info().Any("roles", claims.Roles).Str("requiredRole", requiredRole).Msg("User does not have required role")

		return c.Status(http.StatusForbidden).JSON(map[string]string{"error": "forbidden"})
	}
}
