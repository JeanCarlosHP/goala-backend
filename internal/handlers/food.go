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
			"success": false,
			"error":   "missing_query_parameter",
			"message": "query parameter 'q' is required",
		})
	}

	foods, err := h.foodService.SearchFoods(c.UserContext(), query)
	if err != nil {
		h.logger.Error("failed to search foods", "error", err, "query", query)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error":   "search_foods_failed",
			"message": "failed to search foods",
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    foods,
	})
}

func (h *FoodHandler) GetRecentFoods(c *fiber.Ctx) error {
	userID := c.Locals("userID").(string)

	parsedUserID, err := uuid.Parse(userID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "invalid_user_id",
			"message": "invalid user ID",
		})
	}

	foods, err := h.foodService.GetRecentFoods(c.UserContext(), parsedUserID)
	if err != nil {
		h.logger.Error("failed to get recent foods", "error", err, "user_id", userID)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error":   "get_recent_foods_failed",
			"message": "failed to get recent foods",
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    foods,
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
			"success": false,
			"error":   "invalid_request_body",
			"message": "invalid request body",
		})
	}

	if req.FoodName == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "missing_food_name",
			"message": "food_name is required",
		})
	}

	// TODO: Implement AI-based autocomplete logic here.
	result := struct {
		AutocompleteResponse string `json:"autocomplete_response"`
	}{
		AutocompleteResponse: req.FoodName,
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    result,
	})
}

func (h *FoodHandler) CreateFoodItem(c *fiber.Ctx) error {
	var req domain.CreateFoodItemRequest
	if err := c.BodyParser(&req); err != nil {
		h.logger.Error("failed to parse request body", "error", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "invalid_request_body",
			"message": "invalid request body",
		})
	}

	if err := h.validator.Struct(req); err != nil {
		h.logger.Error("validation failed", "error", err, "request", req)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "validation_failed",
			"message": err.Error(),
		})
	}

	foodItem, err := h.foodService.CreateFoodItem(c.UserContext(), &req)
	if err != nil {
		h.logger.Error("failed to create food item", "error", err, "request", req)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error":   "create_food_item_failed",
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
			"error":   "invalid_food_item_id",
			"message": "invalid food item ID",
		})
	}

	foodItem, err := h.foodService.GetFoodItem(c.UserContext(), id)
	if err != nil {
		h.logger.Error("failed to get food item", "error", err, "id", idParam)
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"error":   "food_item_not_found",
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
	})
}

func (h *FoodHandler) UpdateFoodItem(c *fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "invalid_food_item_id",
			"message": "invalid food item ID",
		})
	}

	var req domain.UpdateFoodItemRequest
	if err := c.BodyParser(&req); err != nil {
		h.logger.Error("failed to parse request body", "error", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "invalid_request_body",
			"message": "invalid request body",
		})
	}

	if err := h.validator.Struct(req); err != nil {
		h.logger.Error("validation failed", "error", err, "request", req)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "validation_failed",
			"message": err.Error(),
		})
	}

	foodItem, err := h.foodService.UpdateFoodItem(c.UserContext(), id, &req)
	if err != nil {
		h.logger.Error("failed to update food item", "error", err, "id", idParam, "request", req)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error":   "update_food_item_failed",
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
			"error":   "invalid_food_item_id",
			"message": "invalid food item ID",
		})
	}

	if err := h.foodService.DeleteFoodItem(c.UserContext(), id); err != nil {
		h.logger.Error("failed to delete food item", "error", err, "id", idParam)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error":   "delete_food_item_failed",
			"message": "failed to delete food item",
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "food item deleted successfully",
	})
}
