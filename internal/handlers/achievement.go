package handlers

import (
	"context"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/jeancarloshp/calorieai/internal/services"
	"github.com/rs/zerolog/log"
)

type AchievementHandler struct {
	achievementService *services.AchievementService
	validator          *validator.Validate
}

func NewAchievementHandler(achievementService *services.AchievementService) *AchievementHandler {
	return &AchievementHandler{
		achievementService: achievementService,
		validator:          validator.New(),
	}
}

func (h *AchievementHandler) GetAchievements(c *fiber.Ctx) error {
	firebaseUID := c.Locals("firebase_uid").(string)

	ctx := context.Background()
	userID, ok := c.Locals("user_id").(uuid.UUID)
	if !ok {
		log.Error().Str("firebase_uid", firebaseUID).Msg("Invalid user ID")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "invalid user id",
		})
	}

	achievements, err := h.achievementService.GetUserAchievements(ctx, userID)
	if err != nil {
		log.Error().Err(err).Str("user_id", userID.String()).Msg("Failed to get achievements")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "failed to get achievements",
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    achievements,
		"message": "achievements retrieved successfully",
	})
}

func (h *AchievementHandler) SyncAchievements(c *fiber.Ctx) error {
	firebaseUID := c.Locals("firebase_uid").(string)

	ctx := context.Background()
	userID, ok := c.Locals("user_id").(uuid.UUID)
	if !ok {
		log.Error().Str("firebase_uid", firebaseUID).Msg("Invalid user ID")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "invalid user id",
		})
	}

	achievements, err := h.achievementService.SyncAchievements(ctx, userID)
	if err != nil {
		log.Error().Err(err).Str("user_id", userID.String()).Msg("Failed to sync achievements")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "failed to sync achievements",
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    achievements,
		"message": "achievements synced successfully",
	})
}
