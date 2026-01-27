package middleware

import (
	"context"
	"strings"

	firebase "firebase.google.com/go/v4"
	"github.com/gofiber/fiber/v2"
	fbApp "github.com/jeancarloshp/calorieai/pkg/firebase"

	"github.com/jeancarloshp/calorieai/internal/domain"
	"github.com/jeancarloshp/calorieai/internal/repositories"
)

func AuthRequired(firebaseApp *firebase.App, logger domain.Logger) fiber.Handler {
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
			logger.Error("Failed to get auth client", "error", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "authentication service unavailable",
			})
		}

		decodedToken, err := authClient.VerifyIDToken(ctx, token)
		if err != nil {
			logger.Warn("Invalid or expired token", "error", err)
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "invalid or expired token",
			})
		}

		c.Locals("firebase_uid", decodedToken.UID)

		return c.Next()
	}
}

func UserContext(userRepo *repositories.UserRepository, logger domain.Logger) fiber.Handler {
	return func(c *fiber.Ctx) error {
		firebaseUID, ok := c.Locals("firebase_uid").(string)
		if !ok || firebaseUID == "" {
			logger.Warn("Firebase UID not found in context")
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "authentication required",
			})
		}

		ctx := context.Background()
		user, err := userRepo.GetByFirebaseUID(ctx, firebaseUID)
		if err != nil {
			logger.Error("User not found", "firebase_uid", firebaseUID, "error", err)
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "user not registered",
			})
		}

		c.Locals("user_id", user.ID)
		c.Locals("user", user)

		return c.Next()
	}
}
