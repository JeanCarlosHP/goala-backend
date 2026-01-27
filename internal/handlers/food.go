package handlers

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/jeancarloshp/calorieai/internal/domain"
	"github.com/jeancarloshp/calorieai/internal/services"
)

type FoodHandler struct {
	foodService *services.FoodService
	validator   *validator.Validate
	logger      domain.Logger
}

func NewFoodHandler(foodService *services.FoodService, validator *validator.Validate, logger domain.Logger) *FoodHandler {
	return &FoodHandler{
		foodService: foodService,
		validator:   validator,
		logger:      logger,
	}
}

func (h *FoodHandler) SearchFoods(c *fiber.Ctx) error {
	query := c.Query("q")
	if query == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "query parameter 'q' is required",
		})
	}

	foods, err := h.foodService.SearchFoods(c.Context(), query)
	if err != nil {
		h.logger.Error("failed to search foods", "error", err, "query", query)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to search foods",
		})
	}

	return c.JSON(fiber.Map{
		"foods": foods,
	})
}

func (h *FoodHandler) GetRecentFoods(c *fiber.Ctx) error {
	userID := c.Locals("userID").(string)

	parsedUserID, err := uuid.Parse(userID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid user ID",
		})
	}

	foods, err := h.foodService.GetRecentFoods(c.Context(), parsedUserID)
	if err != nil {
		h.logger.Error("failed to get recent foods", "error", err, "user_id", userID)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to get recent foods",
		})
	}

	return c.JSON(fiber.Map{
		"foods": foods,
	})
}

type AutocompleteRequest struct {
	FoodName string `json:"food_name" validate:"required"`
}

type AutocompleteResponse struct {
	Name            string  `json:"name"`
	CaloriesPer100g int     `json:"calories_per_100g"`
	ProteinPer100g  float64 `json:"protein_per_100g"`
	CarbsPer100g    float64 `json:"carbs_per_100g"`
	FatPer100g      float64 `json:"fat_per_100g"`
}

func (h *FoodHandler) AutocompleteFoodMacros(c *fiber.Ctx) error {
	var req AutocompleteRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid request body",
		})
	}

	if req.FoodName == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "food_name is required",
		})
	}

	// TODO: Implement AI-based autocomplete logic here.
	result := struct {
		AutocompleteResponse string `json:"autocomplete_response"`
	}{
		AutocompleteResponse: req.FoodName,
	}

	return c.JSON(result)
}

func (h *FoodHandler) CreateFoodItem(c *fiber.Ctx) error {
	var req domain.CreateFoodItemRequest
	if err := c.BodyParser(&req); err != nil {
		h.logger.Error("failed to parse request body", "error", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "invalid request body",
		})
	}

	if err := h.validator.Struct(req); err != nil {
		h.logger.Error("validation failed", "error", err, "request", req)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "validation failed",
			"errors":  err.Error(),
		})
	}

	foodItem, err := h.foodService.CreateFoodItem(c.Context(), &req)
	if err != nil {
		h.logger.Error("failed to create food item", "error", err, "request", req)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "failed to create food item",
		})
	}

	response := domain.FoodItemResponse{
		ID:          foodItem.ID,
		MealID:      foodItem.MealID,
		Name:        foodItem.Name,
		PortionSize: foodItem.PortionSize,
		PortionUnit: foodItem.PortionUnit,
		Calories:    foodItem.Calories,
		ProteinG:    foodItem.ProteinG,
		CarbsG:      foodItem.CarbsG,
		FatG:        foodItem.FatG,
		Source:      foodItem.Source,
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"success": true,
		"data":    response,
		"message": "food item created successfully",
	})
}

func (h *FoodHandler) GetFoodItem(c *fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "invalid food item ID",
		})
	}

	foodItem, err := h.foodService.GetFoodItem(c.Context(), id)
	if err != nil {
		h.logger.Error("failed to get food item", "error", err, "id", idParam)
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"message": "food item not found",
		})
	}

	response := domain.FoodItemResponse{
		ID:          foodItem.ID,
		MealID:      foodItem.MealID,
		Name:        foodItem.Name,
		PortionSize: foodItem.PortionSize,
		PortionUnit: foodItem.PortionUnit,
		Calories:    foodItem.Calories,
		ProteinG:    foodItem.ProteinG,
		CarbsG:      foodItem.CarbsG,
		FatG:        foodItem.FatG,
		Source:      foodItem.Source,
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    response,
		"message": "food item retrieved successfully",
	})
}

func (h *FoodHandler) UpdateFoodItem(c *fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "invalid food item ID",
		})
	}

	var req domain.UpdateFoodItemRequest
	if err := c.BodyParser(&req); err != nil {
		h.logger.Error("failed to parse request body", "error", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "invalid request body",
		})
	}

	if err := h.validator.Struct(req); err != nil {
		h.logger.Error("validation failed", "error", err, "request", req)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "validation failed",
			"errors":  err.Error(),
		})
	}

	foodItem, err := h.foodService.UpdateFoodItem(c.Context(), id, &req)
	if err != nil {
		h.logger.Error("failed to update food item", "error", err, "id", idParam, "request", req)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "failed to update food item",
		})
	}

	response := domain.FoodItemResponse{
		ID:          foodItem.ID,
		MealID:      foodItem.MealID,
		Name:        foodItem.Name,
		PortionSize: foodItem.PortionSize,
		PortionUnit: foodItem.PortionUnit,
		Calories:    foodItem.Calories,
		ProteinG:    foodItem.ProteinG,
		CarbsG:      foodItem.CarbsG,
		FatG:        foodItem.FatG,
		Source:      foodItem.Source,
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    response,
		"message": "food item updated successfully",
	})
}

func (h *FoodHandler) DeleteFoodItem(c *fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "invalid food item ID",
		})
	}

	if err := h.foodService.DeleteFoodItem(c.Context(), id); err != nil {
		h.logger.Error("failed to delete food item", "error", err, "id", idParam)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "failed to delete food item",
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "food item deleted successfully",
	})
}
