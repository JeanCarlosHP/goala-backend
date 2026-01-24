package handlers

import (
	"context"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/jeancarloshp/calorieai/internal/domain"
	"github.com/jeancarloshp/calorieai/internal/services"
	"github.com/rs/zerolog/log"
)

type StatsHandler struct {
	statsService *services.StatsService
	validator    *validator.Validate
}

func NewStatsHandler(statsService *services.StatsService) *StatsHandler {
	return &StatsHandler{
		statsService: statsService,
		validator:    validator.New(),
	}
}

func (h *StatsHandler) GetStats(c *fiber.Ctx) error {
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

	stats, err := h.statsService.GetUserStats(ctx, userID)
	if err != nil {
		log.Error().Err(err).Str("user_id", userID.String()).Msg("Failed to get stats")
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
		log.Error().Err(err).Msg("Failed to parse query")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "invalid query parameters",
		})
	}

	if err := h.validator.Struct(query); err != nil {
		log.Error().Err(err).Msg("Validation failed")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "validation failed",
			"errors":  err.Error(),
		})
	}

	startDate, err := time.Parse("2006-01-02", query.StartDate)
	if err != nil {
		log.Error().Err(err).Str("start_date", query.StartDate).Msg("Invalid start date format")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "invalid start date format, expected yyyy-MM-dd",
		})
	}

	endDate, err := time.Parse("2006-01-02", query.EndDate)
	if err != nil {
		log.Error().Err(err).Str("end_date", query.EndDate).Msg("Invalid end date format")
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

	ctx := context.Background()
	userID, ok := c.Locals("user_id").(uuid.UUID)
	if !ok {
		log.Error().Str("firebase_uid", firebaseUID).Msg("Invalid user ID")
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
		log.Error().Err(err).Str("user_id", userID.String()).Msg("Failed to get stats range")
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
