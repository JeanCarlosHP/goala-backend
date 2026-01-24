package handlers

import (
	"context"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/jeancarloshp/calorieai/internal/domain"
	"github.com/jeancarloshp/calorieai/internal/services"
	"github.com/rs/zerolog/log"
)

type FeedbackHandler struct {
	feedbackService *services.FeedbackService
	userService     *services.UserService
	validator       *validator.Validate
}

func NewFeedbackHandler(feedbackService *services.FeedbackService, userService *services.UserService) *FeedbackHandler {
	return &FeedbackHandler{
		feedbackService: feedbackService,
		userService:     userService,
		validator:       validator.New(),
	}
}

func (h *FeedbackHandler) CreateFeedback(c *fiber.Ctx) error {
	firebaseUID := c.Locals("firebase_uid").(string)

	var req domain.CreateFeedbackRequest
	if err := c.BodyParser(&req); err != nil {
		log.Error().Err(err).Msg("Invalid request body")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "invalid request body",
		})
	}

	if err := h.validator.Struct(&req); err != nil {
		log.Error().Err(err).Msg("Validation error")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "validation error",
			"errors":  err.Error(),
		})
	}

	ctx := context.Background()
	user, err := h.userService.GetUserByFirebaseUID(ctx, firebaseUID)
	if err != nil {
		log.Error().Err(err).Str("firebase_uid", firebaseUID).Msg("User not found")
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"message": "user not found",
		})
	}

	if err := h.feedbackService.CreateFeedback(ctx, user.ID, &req); err != nil {
		log.Error().Err(err).Str("user_id", user.ID.String()).Msg("Failed to create feedback")
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

func (h *FeedbackHandler) GetUserFeedback(c *fiber.Ctx) error {
	firebaseUID := c.Locals("firebase_uid").(string)

	ctx := context.Background()
	user, err := h.userService.GetUserByFirebaseUID(ctx, firebaseUID)
	if err != nil {
		log.Error().Err(err).Str("firebase_uid", firebaseUID).Msg("User not found")
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"message": "user not found",
		})
	}

	feedbacks, err := h.feedbackService.GetUserFeedback(ctx, user.ID)
	if err != nil {
		log.Error().Err(err).Str("user_id", user.ID.String()).Msg("Failed to get user feedback")
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

func (h *FeedbackHandler) GetFeedback(c *fiber.Ctx) error {
	feedbackID := c.Params("id")

	id, err := uuid.Parse(feedbackID)
	if err != nil {
		log.Error().Err(err).Str("feedback_id", feedbackID).Msg("Invalid feedback ID")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "invalid feedback id",
		})
	}

	ctx := context.Background()
	feedback, err := h.feedbackService.GetFeedback(ctx, id)
	if err != nil {
		log.Error().Err(err).Str("feedback_id", feedbackID).Msg("Failed to get feedback")
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
