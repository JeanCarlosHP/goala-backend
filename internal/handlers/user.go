package handlers

import (
	"context"
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/jeancarloshp/calorieai/internal/domain"
	"github.com/jeancarloshp/calorieai/internal/services"
	"github.com/rs/zerolog/log"
)

type UserHandler struct {
	userService *services.UserService
	s3Service   *services.S3Service
	validator   *validator.Validate
}

func NewUserHandler(userService *services.UserService, s3Service *services.S3Service) *UserHandler {
	return &UserHandler{
		userService: userService,
		s3Service:   s3Service,
		validator:   validator.New(),
	}
}

func (h *UserHandler) GetProfile(c *fiber.Ctx) error {
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

	profile, err := h.userService.GetUserProfile(ctx, user.ID)
	if err != nil {
		log.Error().Err(err).Str("user_id", user.ID.String()).Msg("Failed to get user profile")
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

	reqUserID, err := uuid.Parse(req.ID)
	if err != nil {
		log.Error().Err(err).Str("id", req.ID).Msg("Invalid user ID")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "invalid user id",
		})
	}

	if user.ID != reqUserID {
		log.Warn().
			Str("authenticated_user_id", user.ID.String()).
			Str("requested_user_id", reqUserID.String()).
			Msg("Unauthorized profile update attempt")
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"success": false,
			"message": "unauthorized",
		})
	}

	profile, err := h.userService.UpdateUserProfile(ctx, user.ID, req)
	if err != nil {
		log.Error().Err(err).Str("user_id", user.ID.String()).Msg("Failed to update user profile")
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

	uploadURL, photoURL, err := h.s3Service.GenerateUploadPresignedURL(ctx, firebaseUID, req.ContentType, req.FileSize)
	if err != nil {
		log.Error().Err(err).Str("user_id", user.ID.String()).Msg("Failed to generate presigned URL")
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

	fmt.Print(req)

	ctx := context.Background()
	user, err := h.userService.GetUserByFirebaseUID(ctx, firebaseUID)
	if err != nil {
		log.Error().Err(err).Str("firebase_uid", firebaseUID).Msg("User not found")
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"message": "user not found",
		})
	}

	profile, err := h.userService.PatchUserPreferences(ctx, user.ID, req)
	if err != nil {
		log.Error().Err(err).Str("user_id", user.ID.String()).Msg("Failed to update user preferences")
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
		log.Error().Err(err).Msg("Failed to update goal")
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
