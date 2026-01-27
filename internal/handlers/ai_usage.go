package handlers

import (
	"github.com/gofiber/fiber/v2"

	"github.com/jeancarloshp/calorieai/internal/domain"
	"github.com/jeancarloshp/calorieai/internal/domain/enum"
	"github.com/jeancarloshp/calorieai/internal/services"
)

type AIUsageHandler struct {
	aiUsageService *services.AIUsageService
	logger         domain.Logger
}

func NewAIUsageHandler(aiUsageService *services.AIUsageService, log domain.Logger) *AIUsageHandler {
	return &AIUsageHandler{
		aiUsageService: aiUsageService,
		logger:         log,
	}
}

func (h *AIUsageHandler) GetUsage(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(string)
	ctx := c.UserContext()

	usages, err := h.aiUsageService.ListUserUsage(ctx, userID)
	if err != nil {
		h.logger.Error("Failed to get AI usage", map[string]interface{}{
			"error":   err.Error(),
			"user_id": userID,
		})
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "failed to get usage data",
		})
	}

	usageData := make([]fiber.Map, 0, len(usages))
	for _, usage := range usages {
		usageData = append(usageData, fiber.Map{
			"feature":   usage.Feature,
			"used":      usage.UsageCount,
			"quota":     usage.Quota,
			"remaining": usage.RemainingQuota(),
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    usageData,
	})
}

func (h *AIUsageHandler) CheckFeatureQuota(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(string)
	featureStr := c.Params("feature")

	ctx := c.UserContext()

	feature := enum.AIFeature(featureStr)
	if !feature.IsValid() {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "invalid feature",
		})
	}

	usage, err := h.aiUsageService.GetUsage(ctx, userID, feature)
	if err != nil {
		h.logger.Error("Failed to check quota", map[string]interface{}{
			"error":   err.Error(),
			"user_id": userID,
			"feature": feature,
		})
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "failed to check quota",
		})
	}

	if usage == nil {
		return c.JSON(fiber.Map{
			"success": true,
			"data": fiber.Map{
				"has_quota": true,
				"used":      0,
				"quota":     0,
				"remaining": 0,
			},
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data": fiber.Map{
			"has_quota": usage.HasQuota(),
			"used":      usage.UsageCount,
			"quota":     usage.Quota,
			"remaining": usage.RemainingQuota(),
		},
	})
}
