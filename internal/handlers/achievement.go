package handlers

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/jeancarloshp/calorieai/internal/domain"
	"github.com/jeancarloshp/calorieai/internal/services"
)

type AchievementHandler struct {
	achievementService *services.AchievementService
	validator          *validator.Validate
	logger             domain.Logger
}

func NewAchievementHandler(achievementService *services.AchievementService, logger domain.Logger) *AchievementHandler {
	return &AchievementHandler{
		achievementService: achievementService,
		validator:          validator.New(),
		logger:             logger,
	}
}

func (h *AchievementHandler) GetAchievements(c *fiber.Ctx) error {
	firebaseUID := c.Locals("firebase_uid").(string)

	ctx := c.UserContext()
	userID, ok := c.Locals("user_id").(uuid.UUID)
	if !ok {
		h.logger.Error("Invalid user ID", "firebase_uid", firebaseUID)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "invalid user id",
		})
	}

	achievements, err := h.achievementService.GetUserAchievements(ctx, userID)
	if err != nil {
		h.logger.Error("Failed to get achievements", "user_id", userID.String(), "error", err)
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

	ctx := c.UserContext()
	userID, ok := c.Locals("user_id").(uuid.UUID)
	if !ok {
		h.logger.Error("Invalid user ID", "firebase_uid", firebaseUID)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "invalid user id",
		})
	}

	achievements, err := h.achievementService.SyncAchievements(ctx, userID)
	if err != nil {
		h.logger.Error("Failed to sync achievements", "user_id", userID.String(), "error", err)
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
