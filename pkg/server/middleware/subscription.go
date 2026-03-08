package middleware

import (
	"github.com/gofiber/fiber/v3"
	"go.opentelemetry.io/otel"

	"github.com/jeancarloshp/calorieai/internal/domain"
	"github.com/jeancarloshp/calorieai/internal/domain/enum"
	"github.com/jeancarloshp/calorieai/internal/services"
)

func SubscriptionRequired(subscriptionService *services.SubscriptionService, log domain.Logger) fiber.Handler {
	return func(c fiber.Ctx) error {
		ctx := c.Context()

		tr := otel.Tracer("middleware/subscription.go")
		ctx, span := tr.Start(ctx, "SubscriptionRequired")
		defer span.End()

		userID, ok := c.Locals("user_id").(string)
		if !ok || userID == "" {
			log.Warn("Missing user_id in context", nil)
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"success": false,
				"message": "authentication required",
			})
		}

		hasAccess, err := subscriptionService.ValidateAccess(ctx, userID)
		if err != nil {
			log.Error("Failed to validate subscription", map[string]interface{}{
				"error":   err.Error(),
				"user_id": userID,
			})
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"success": false,
				"message": "failed to validate subscription",
			})
		}

		if !hasAccess {
			return c.Status(fiber.StatusPaymentRequired).JSON(fiber.Map{
				"success": false,
				"message": "active subscription required",
				"code":    "SUBSCRIPTION_REQUIRED",
			})
		}

		return c.Next()
	}
}

func AIQuotaCheck(aiUsageService *services.AIUsageService, feature enum.AIFeature, log domain.Logger) fiber.Handler {
	return func(c fiber.Ctx) error {
		ctx := c.Context()

		tr := otel.Tracer("middleware/subscription.go")
		ctx, span := tr.Start(ctx, "AIQuotaCheck")
		defer span.End()

		userID, ok := c.Locals("user_id").(string)
		if !ok || userID == "" {
			log.Warn("Missing user_id in context", nil)
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"success": false,
				"message": "authentication required",
			})
		}

		if err := aiUsageService.CheckAndIncrementUsage(ctx, userID, feature); err != nil {
			if services.IsQuotaExceededError(err) {
				log.Info("Quota exceeded", map[string]interface{}{
					"user_id": userID,
					"feature": feature,
				})
				return c.Status(fiber.StatusPaymentRequired).JSON(fiber.Map{
					"success": false,
					"message": "quota exceeded for this feature",
					"code":    "QUOTA_EXCEEDED",
					"feature": feature,
				})
			}

			log.Error("Failed to check AI quota", map[string]interface{}{
				"error":   err.Error(),
				"user_id": userID,
				"feature": feature,
			})
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"success": false,
				"message": "failed to validate quota",
			})
		}

		return c.Next()
	}
}
