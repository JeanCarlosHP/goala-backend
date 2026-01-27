package handlers

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/jeancarloshp/calorieai/internal/domain"
	"github.com/jeancarloshp/calorieai/internal/domain/enum"
	"github.com/jeancarloshp/calorieai/internal/services"
)

type FoodRecognitionHandler struct {
	foodRecognitionService *services.FoodRecognitionService
	barcodeService         *services.BarcodeService
	validator              *validator.Validate
	aiUsageService         *services.AIUsageService
	s3Service              *services.S3Service
	logger                 domain.Logger
}

func NewFoodRecognitionHandler(
	foodRecognitionService *services.FoodRecognitionService,
	barcodeService *services.BarcodeService,
	aiUsageService *services.AIUsageService,
	s3Service *services.S3Service,
	logger domain.Logger,
) *FoodRecognitionHandler {
	return &FoodRecognitionHandler{
		foodRecognitionService: foodRecognitionService,
		barcodeService:         barcodeService,
		validator:              validator.New(),
		aiUsageService:         aiUsageService,
		s3Service:              s3Service,
		logger:                 logger,
	}
}

func (h *FoodRecognitionHandler) RecognizeFood(c *fiber.Ctx) error {
	userID, ok := c.Locals("user_id").(uuid.UUID)
	if !ok {
		h.logger.Warn("Missing user_id in context")
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"message": "authentication required",
		})
	}

	userIDStr := userID.String()

	ctx := c.UserContext()
	if err := h.aiUsageService.CheckAndIncrementUsage(ctx, userIDStr, enum.FeatureFoodRecognition); err != nil {
		if services.IsQuotaExceededError(err) {
			h.logger.Warn("Quota exceeded for user %s and feature %s", userID, enum.FeatureFoodRecognition)
			return c.Status(fiber.StatusPaymentRequired).JSON(fiber.Map{
				"success": false,
				"message": "quota exceeded for this feature",
				"code":    "QUOTA_EXCEEDED",
				"feature": enum.FeatureFoodRecognition,
			})
		}

		h.logger.Error("Failed to validate quota: %v, user_id: %s", err, userIDStr)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "failed to validate quota",
		})
	}

	var req domain.FoodRecognitionRequest
	if err := c.BodyParser(&req); err != nil {
		h.logger.Error("Invalid request body", "error", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "invalid request body",
		})
	}

	if err := h.validator.Struct(&req); err != nil {
		h.logger.Error("Validation failed: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "validation error",
			"errors":  err.Error(),
		})
	}

	// Processar reconhecimento de forma síncrona
	result, err := h.foodRecognitionService.RecognizeFoodByPath(ctx, req.ImagePath, &req)
	if err != nil {
		h.logger.Error("Failed to recognize food: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "failed to recognize food",
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    result,
		"message": "food recognized successfully",
	})
}

func (h *FoodRecognitionHandler) GetFoodByBarcode(c *fiber.Ctx) error {
	ctx := c.UserContext()

	barcode := c.Params("barcode")
	if barcode == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "barcode is required",
		})
	}

	result, err := h.barcodeService.GetFoodByBarcode(ctx, barcode)
	if err != nil {
		h.logger.Error("Failed to get food by barcode: %v, barcode: %s", err, barcode)
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"message": "food not found for barcode",
		})
	}

	if err := h.validator.Struct(result); err != nil {
		h.logger.Error("Validation failed for barcode response: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "invalid response from barcode service",
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    result,
		"message": "food retrieved successfully",
	})
}

func (h *FoodRecognitionHandler) EstimateQuantity(c *fiber.Ctx) error {
	ctx := c.UserContext()

	var req domain.EstimateQuantityRequest
	if err := c.BodyParser(&req); err != nil {
		h.logger.Error("Invalid request body", "error", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "invalid request body",
		})
	}

	if err := h.validator.Struct(&req); err != nil {
		h.logger.Error("Validation failed: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "validation error",
			"errors":  err.Error(),
		})
	}

	// Processar estimativa de quantidade de forma síncrona
	result, err := h.foodRecognitionService.EstimateQuantityByPath(ctx, req.ImagePath, &req)
	if err != nil {
		h.logger.Error("Failed to estimate quantity: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "failed to estimate quantity",
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    result,
		"message": "quantity estimated successfully",
	})
}

func (h *FoodRecognitionHandler) GenerateFoodImageUploadURL(c *fiber.Ctx) error {
	userID, ok := c.Locals("user_id").(uuid.UUID)
	if !ok {
		h.logger.Warn("Missing user_id in context")
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"message": "authentication required",
		})
	}

	var req domain.FoodImageUploadRequest
	if err := c.BodyParser(&req); err != nil {
		h.logger.Error("Invalid request body", "error", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "invalid request body",
		})
	}

	if err := h.validator.Struct(&req); err != nil {
		h.logger.Error("Validation error", "error", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "validation error",
			"errors":  err.Error(),
		})
	}

	ctx := c.UserContext()
	uploadURL, imagePath, err := h.s3Service.GenerateFoodImageUploadPresignedURL(ctx, userID.String(), req.ContentType, req.FileSize)
	if err != nil {
		h.logger.Error("Failed to generate presigned URL", "user_id", userID.String(), "error", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data": domain.FoodImageUploadResponse{
			UploadURL: uploadURL,
			ImagePath: imagePath,
			ExpiresIn: 300,
		},
		"message": "presigned URL generated successfully",
	})
}
