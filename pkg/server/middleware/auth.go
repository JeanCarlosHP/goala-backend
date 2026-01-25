package middleware

import (
	"context"
	"strings"

	firebase "firebase.google.com/go/v4"
	"github.com/gofiber/fiber/v2"
	fbApp "github.com/jeancarloshp/calorieai/pkg/firebase"
	"github.com/rs/zerolog/log"

	"github.com/jeancarloshp/calorieai/internal/repositories"
)

func AuthRequired(firebaseApp *firebase.App) fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "missing authorization header",
			})
		}

		token := strings.TrimPrefix(authHeader, "Bearer ")
		if token == authHeader {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "invalid authorization header format",
			})
		}

		ctx := context.Background()
		authClient, err := fbApp.GetAuthClient(ctx, firebaseApp)
		if err != nil {
			log.Error().Err(err).Msg("Failed to get auth client")
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "authentication service unavailable",
			})
		}

		decodedToken, err := authClient.VerifyIDToken(ctx, token)
		if err != nil {
			log.Error().Err(err).Msg("Invalid token")
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "invalid or expired token",
			})
		}

		c.Locals("firebase_uid", decodedToken.UID)

		return c.Next()
	}
}

func UserContext(userRepo *repositories.UserRepository) fiber.Handler {
	return func(c *fiber.Ctx) error {
		firebaseUID, ok := c.Locals("firebase_uid").(string)
		if !ok || firebaseUID == "" {
			log.Warn().Msg("Missing firebase_uid in context")
			return c.Next()
		}

		ctx := context.Background()
		user, err := userRepo.GetByFirebaseUID(ctx, firebaseUID)
		if err != nil {
			log.Warn().Err(err).Str("firebase_uid", firebaseUID).Msg("User not found")
			return c.Next()
		}

		c.Locals("user_id", user.ID)
		c.Locals("user", user)

		return c.Next()
	}
}
