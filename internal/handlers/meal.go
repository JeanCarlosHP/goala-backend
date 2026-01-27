package handlers

import (
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/jeancarloshp/calorieai/internal/domain"
	"github.com/jeancarloshp/calorieai/internal/services"
)

type MealHandler struct {
	mealService *services.MealService
	userService *services.UserService
	validator   *validator.Validate
	logger      domain.Logger
}

func NewMealHandler(mealService *services.MealService, userService *services.UserService, logger domain.Logger) *MealHandler {
	return &MealHandler{
		mealService: mealService,
		userService: userService,
		validator:   validator.New(),
		logger:      logger,
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

	ctx := c.UserContext()
	meal, err := h.mealService.CreateMeal(ctx, userID, req)
	if err != nil {
		h.logger.Error("Failed to create meal", "firebase_uid", firebaseUID, "error", err)
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

	ctx := c.UserContext()
	meals, err := h.mealService.GetMealsByDate(ctx, userID, date)
	if err != nil {
		h.logger.Error("Failed to get meals", "user_id", userID.String(), "error", err)
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

	ctx := c.UserContext()

	goal, err := h.userService.GetUserGoal(ctx, userID)
	if err != nil {
		h.logger.Warn("Failed to get user goal", "user_id", userID.String(), "error", err)
	}

	summary, err := h.mealService.GetDailySummary(ctx, userID, date, goal)
	if err != nil {
		h.logger.Error("Failed to get daily summary", "user_id", userID.String(), "error", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to get daily summary",
		})
	}

	return c.JSON(summary)
}
