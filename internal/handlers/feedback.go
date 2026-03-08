package handlers

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
	"github.com/jeancarloshp/calorieai/internal/domain"
	"github.com/jeancarloshp/calorieai/internal/services"
)

type FeedbackHandler struct {
	feedbackService *services.FeedbackService
	userService     *services.UserService
	validator       *validator.Validate
	logger          domain.Logger
}

func NewFeedbackHandler(feedbackService *services.FeedbackService, userService *services.UserService, logger domain.Logger) *FeedbackHandler {
	return &FeedbackHandler{
		feedbackService: feedbackService,
		userService:     userService,
		validator:       validator.New(),
		logger:          logger,
	}
}

func (h *FeedbackHandler) CreateFeedback(c fiber.Ctx) error {
	firebaseUID := c.Locals("firebase_uid").(string)

	var req domain.CreateFeedbackRequest
	if err := c.Bind().JSON(&req); err != nil {
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

	ctx := c.Context()
	user, err := h.userService.GetUserByFirebaseUID(ctx, firebaseUID)
	if err != nil {
		h.logger.Error("User not found", "firebase_uid", firebaseUID, "error", err)
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"message": "user not found",
		})
	}

	if err := h.feedbackService.CreateFeedback(ctx, user.ID, &req); err != nil {
		h.logger.Error("Failed to create feedback", "user_id", user.ID.String(), "error", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "failed to create feedback",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "feedback submitted successfully",
	})
}

func (h *FeedbackHandler) GetUserFeedback(c fiber.Ctx) error {
	firebaseUID := c.Locals("firebase_uid").(string)

	ctx := c.Context()
	user, err := h.userService.GetUserByFirebaseUID(ctx, firebaseUID)
	if err != nil {
		h.logger.Error("User not found", "firebase_uid", firebaseUID, "error", err)
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"message": "user not found",
		})
	}

	feedbacks, err := h.feedbackService.GetUserFeedback(ctx, user.ID)
	if err != nil {
		h.logger.Error("Failed to get user feedback", "user_id", user.ID.String(), "error", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "failed to get feedback",
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    feedbacks,
		"message": "feedback retrieved successfully",
	})
}

func (h *FeedbackHandler) GetFeedback(c fiber.Ctx) error {
	feedbackID := c.Params("id")

	id, err := uuid.Parse(feedbackID)
	if err != nil {
		h.logger.Error("Invalid feedback ID", "feedback_id", feedbackID, "error", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "invalid feedback id",
		})
	}

	ctx := c.Context()
	feedback, err := h.feedbackService.GetFeedback(ctx, id)
	if err != nil {
		h.logger.Error("Failed to get feedback", "feedback_id", feedbackID, "error", err)
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"message": "feedback not found",
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    feedback,
		"message": "feedback retrieved successfully",
	})
}
