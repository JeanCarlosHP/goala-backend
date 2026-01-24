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

type MealHandler struct {
	mealService *services.MealService
	userService *services.UserService
	validator   *validator.Validate
}

func NewMealHandler(mealService *services.MealService, userService *services.UserService) *MealHandler {
	return &MealHandler{
		mealService: mealService,
		userService: userService,
		validator:   validator.New(),
	}
}

func (h *MealHandler) CreateMeal(c *fiber.Ctx) error {
	firebaseUID := c.Locals("firebase_uid").(string)
	userID := c.Locals("user_id").(uuid.UUID)

	var req domain.CreateMealRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid request body",
		})
	}

	if err := h.validator.Struct(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	ctx := context.Background()
	meal, err := h.mealService.CreateMeal(ctx, userID, req)
	if err != nil {
		log.Error().Err(err).Str("firebase_uid", firebaseUID).Msg("Failed to create meal")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to create meal",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(meal)
}

func (h *MealHandler) GetMeals(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uuid.UUID)

	dateStr := c.Query("date")
	if dateStr == "" {
		dateStr = time.Now().Format("2006-01-02")
	}

	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid date format, use YYYY-MM-DD",
		})
	}

	ctx := context.Background()
	meals, err := h.mealService.GetMealsByDate(ctx, userID, date)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get meals")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to get meals",
		})
	}

	return c.JSON(meals)
}

func (h *MealHandler) GetDailySummary(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uuid.UUID)

	dateStr := c.Query("date")
	if dateStr == "" {
		dateStr = time.Now().Format("2006-01-02")
	}

	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid date format, use YYYY-MM-DD",
		})
	}

	ctx := context.Background()

	goal, err := h.userService.GetUserGoal(ctx, userID)
	if err != nil {
		log.Warn().Err(err).Msg("Failed to get user goals, using defaults")
	}

	summary, err := h.mealService.GetDailySummary(ctx, userID, date, goal)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get daily summary")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to get daily summary",
		})
	}

	return c.JSON(summary)
}
