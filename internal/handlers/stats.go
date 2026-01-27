package handlers

import (
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/jeancarloshp/calorieai/internal/domain"
	"github.com/jeancarloshp/calorieai/internal/services"
)

type StatsHandler struct {
	statsService *services.StatsService
	validator    *validator.Validate
	logger       domain.Logger
}

func NewStatsHandler(statsService *services.StatsService, logger domain.Logger) *StatsHandler {
	return &StatsHandler{
		statsService: statsService,
		validator:    validator.New(),
		logger:       logger,
	}
}

func (h *StatsHandler) GetStats(c *fiber.Ctx) error {
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

	stats, err := h.statsService.GetUserStats(ctx, userID)
	if err != nil {
		h.logger.Error("Failed to get stats", "user_id", userID.String(), "error", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "failed to get stats",
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    stats,
		"message": "stats retrieved successfully",
	})
}

func (h *StatsHandler) GetStatsRange(c *fiber.Ctx) error {
	firebaseUID := c.Locals("firebase_uid").(string)

	var query domain.StatsRangeQuery
	if err := c.QueryParser(&query); err != nil {
		h.logger.Error("Invalid query parameters", "firebase_uid", firebaseUID, "error", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "invalid query parameters",
		})
	}

	if err := h.validator.Struct(query); err != nil {
		h.logger.Error("Validation failed for query parameters", "firebase_uid", firebaseUID, "error", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "validation failed",
			"errors":  err.Error(),
		})
	}

	startDate, err := time.Parse("2006-01-02", query.StartDate)
	if err != nil {
		h.logger.Error("Invalid start date format", "start_date", query.StartDate, "error", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "invalid start date format, expected yyyy-MM-dd",
		})
	}

	endDate, err := time.Parse("2006-01-02", query.EndDate)
	if err != nil {
		h.logger.Error("Invalid end date format", "end_date", query.EndDate, "error", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "invalid end date format, expected yyyy-MM-dd",
		})
	}

	if startDate.After(endDate) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "start date must be before end date",
		})
	}

	ctx := c.UserContext()
	userID, ok := c.Locals("user_id").(uuid.UUID)
	if !ok {
		h.logger.Error("Invalid user ID", "firebase_uid", firebaseUID)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "invalid user id",
		})
	}

	page := query.Page
	if page == 0 {
		page = 1
	}

	limit := query.Limit
	if limit == 0 {
		limit = 30
	}

	stats, err := h.statsService.GetStatsRange(ctx, userID, startDate, endDate, page, limit)
	if err != nil {
		h.logger.Error("Failed to get stats range", "user_id", userID.String(), "error", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "failed to get stats",
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    stats,
		"message": "stats retrieved successfully",
	})
}
