package handlers

import (
	"context"
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/jeancarloshp/calorieai/internal/domain"
	"github.com/jeancarloshp/calorieai/internal/services"
)

type UserHandler struct {
	userService *services.UserService
	s3Service   *services.S3Service
	validator   *validator.Validate
	logger      domain.Logger
}

func NewUserHandler(userService *services.UserService, s3Service *services.S3Service, logger domain.Logger) *UserHandler {
	return &UserHandler{
		userService: userService,
		s3Service:   s3Service,
		validator:   validator.New(),
		logger:      logger,
	}
}

func (h *UserHandler) GetProfile(c *fiber.Ctx) error {
	firebaseUID := c.Locals("firebase_uid").(string)

	ctx := context.Background()
	user, err := h.userService.GetUserByFirebaseUID(ctx, firebaseUID)
	if err != nil {
		h.logger.Error("User not found", "firebase_uid", firebaseUID, "error", err)
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"message": "user not found",
		})
	}

	profile, err := h.userService.GetUserProfile(ctx, user.ID)
	if err != nil {
		h.logger.Error("Failed to get user profile", "user_id", user.ID.String(), "error", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "failed to get user profile",
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    profile,
		"message": "profile retrieved successfully",
	})
}

func (h *UserHandler) UpdateProfile(c *fiber.Ctx) error {
	firebaseUID := c.Locals("firebase_uid").(string)

	var req domain.UpdateProfileRequest
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

	ctx := context.Background()
	user, err := h.userService.GetUserByFirebaseUID(ctx, firebaseUID)
	if err != nil {
		h.logger.Error("User not found", "firebase_uid", firebaseUID, "error", err)
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"message": "user not found",
		})
	}

	reqUserID, err := uuid.Parse(req.ID)
	if err != nil {
		h.logger.Error("Invalid user ID in request", "user_id", req.ID, "error", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "invalid user id",
		})
	}

	if user.ID != reqUserID {
		h.logger.Error("Unauthorized profile update attempt", "authenticated_user_id", user.ID.String(), "requested_user_id", reqUserID.String())
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"success": false,
			"message": "unauthorized",
		})
	}

	profile, err := h.userService.UpdateUserProfile(ctx, user.ID, req)
	if err != nil {
		h.logger.Error("Failed to update user profile", "user_id", user.ID.String(), "error", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "failed to update user profile",
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    profile,
		"message": "profile updated successfully",
	})
}

func (h *UserHandler) GenerateAvatarUploadURL(c *fiber.Ctx) error {
	firebaseUID := c.Locals("firebase_uid").(string)

	var req domain.AvatarUploadRequest
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

	ctx := context.Background()
	user, err := h.userService.GetUserByFirebaseUID(ctx, firebaseUID)
	if err != nil {
		h.logger.Error("User not found", "firebase_uid", firebaseUID, "error", err)
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"message": "user not found",
		})
	}

	uploadURL, photoURL, err := h.s3Service.GenerateUploadPresignedURL(ctx, firebaseUID, req.ContentType, req.FileSize)
	if err != nil {
		h.logger.Error("Failed to generate presigned URL", "user_id", user.ID.String(), "error", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data": domain.AvatarUploadResponse{
			UploadURL: uploadURL,
			PhotoURL:  photoURL,
			ExpiresIn: 300,
		},
		"message": "presigned URL generated successfully",
	})
}

func (h *UserHandler) PatchUserPreferences(c *fiber.Ctx) error {
	firebaseUID := c.Locals("firebase_uid").(string)

	var req domain.PatchUserPreferencesRequest
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

	fmt.Print(req)

	ctx := context.Background()
	user, err := h.userService.GetUserByFirebaseUID(ctx, firebaseUID)
	if err != nil {
		h.logger.Error("User not found", "firebase_uid", firebaseUID, "error", err)
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"message": "user not found",
		})
	}

	profile, err := h.userService.PatchUserPreferences(ctx, user.ID, req)
	if err != nil {
		h.logger.Error("Failed to update user preferences", "user_id", user.ID.String(), "error", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "failed to update user preferences",
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    profile,
		"message": "preferences updated successfully",
	})
}

func (h *AuthHandler) UpdateGoals(c *fiber.Ctx) error {
	firebaseUID := c.Locals("firebase_uid").(string)

	var req domain.UpdateGoalRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "invalid request body",
		})
	}

	if err := h.validator.Struct(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "validation error",
		})
	}

	ctx := context.Background()
	user, err := h.userService.GetUserByFirebaseUID(ctx, firebaseUID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"message": "user not found",
		})
	}

	profile, err := h.userService.UpdateUserGoal(ctx, user.ID, req)
	if err != nil {
		h.logger.Error("Failed to update goal", "user_id", user.ID.String(), "error", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "failed to update goal",
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    profile,
		"message": "goal updated successfully",
	})
}
