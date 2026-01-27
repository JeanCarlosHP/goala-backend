package handlers

import (
	"io"

	"github.com/gofiber/fiber/v2"

	"github.com/jeancarloshp/calorieai/internal/domain"
	"github.com/jeancarloshp/calorieai/internal/services"
)

type SubscriptionHandler struct {
	subscriptionService *services.SubscriptionService
	revenueCatService   *services.RevenueCatService
	logger              domain.Logger
}

func NewSubscriptionHandler(
	subscriptionService *services.SubscriptionService,
	revenueCatService *services.RevenueCatService,
	log domain.Logger,
) *SubscriptionHandler {
	return &SubscriptionHandler{
		subscriptionService: subscriptionService,
		revenueCatService:   revenueCatService,
		logger:              log,
	}
}

func (h *SubscriptionHandler) GetStatus(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(string)
	ctx := c.Context()

	subscription, err := h.subscriptionService.GetOrCreateSubscription(ctx, userID)
	if err != nil {
		h.logger.Error("Failed to get subscription", map[string]interface{}{
			"error":   err.Error(),
			"user_id": userID,
		})
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "failed to get subscription status",
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data": fiber.Map{
			"is_active":          subscription.IsActive,
			"plan":               subscription.Plan,
			"is_trial":           subscription.IsTrial,
			"current_period_end": subscription.CurrentPeriodEnd,
			"has_access":         subscription.HasAccess(),
		},
	})
}

func (h *SubscriptionHandler) HandleWebhook(c *fiber.Ctx) error {
	ctx := c.Context()

	signature := c.Get("X-Revenuecat-Signature")
	if signature == "" {
		signature = c.Get("Authorization")
	}

	payload, err := io.ReadAll(c.Request().BodyStream())
	if err != nil {
		h.logger.Error("Failed to read webhook payload", map[string]interface{}{
			"error": err.Error(),
		})
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "invalid payload",
		})
	}

	if err := h.revenueCatService.VerifyWebhookSignature(payload, signature); err != nil {
		h.logger.Warn("Invalid webhook signature", map[string]interface{}{
			"error": err.Error(),
		})
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"message": "invalid signature",
		})
	}

	webhook, err := h.revenueCatService.ParseWebhook(payload)
	if err != nil {
		h.logger.Error("Failed to parse webhook", map[string]interface{}{
			"error": err.Error(),
		})
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "invalid webhook format",
		})
	}

	if err := h.subscriptionService.ProcessWebhookEvent(ctx, &webhook.Event); err != nil {
		h.logger.Error("Failed to process webhook event", map[string]interface{}{
			"error":      err.Error(),
			"event_id":   webhook.Event.ID,
			"event_type": webhook.Event.Type,
		})
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "failed to process event",
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "webhook processed",
	})
}
