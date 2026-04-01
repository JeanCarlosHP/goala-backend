package handlers

import (
	"context"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
	"github.com/jeancarloshp/calorieai/internal/domain"
)

type userService interface {
	GetUserByFirebaseUID(ctx context.Context, firebaseUID string) (*domain.User, error)
	GetUserProfile(ctx context.Context, userID uuid.UUID) (*domain.UserProfileResponse, error)
	UpdateUserProfile(ctx context.Context, userID uuid.UUID, req domain.UpdateProfileRequest) (*domain.UserProfileResponse, error)
	PatchUserPreferences(ctx context.Context, userID uuid.UUID, req domain.PatchUserPreferencesRequest) (*domain.UserProfileResponse, error)
}

type avatarUploadService interface {
	GenerateUploadPresignedURL(ctx context.Context, firebaseUID, contentType string, fileSize int64) (string, string, error)
}

type UserHandler struct {
	userService userService
	s3Service   avatarUploadService
	validator   *validator.Validate
	logger      domain.Logger
}

func NewUserHandler(userService userService, s3Service avatarUploadService, logger domain.Logger) *UserHandler {
	validate := validator.New()
	_ = validate.RegisterValidation("notification_time", func(fl validator.FieldLevel) bool {
		value, ok := fl.Field().Interface().(string)
		if !ok {
			return false
		}
		return domain.IsValidReminderTime(value)
	})

	return &UserHandler{
		userService: userService,
		s3Service:   s3Service,
		validator:   validate,
		logger:      logger,
	}
}

func (h *UserHandler) GetProfile(c fiber.Ctx) error {
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

func (h *UserHandler) UpdateProfile(c fiber.Ctx) error {
	firebaseUID := c.Locals("firebase_uid").(string)

	var req domain.UpdateProfileRequest
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

func (h *UserHandler) GenerateAvatarUploadURL(c fiber.Ctx) error {
	firebaseUID := c.Locals("firebase_uid").(string)

	var req domain.AvatarUploadRequest
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

func (h *UserHandler) PatchUserPreferences(c fiber.Ctx) error {
	firebaseUID := c.Locals("firebase_uid").(string)

	var req domain.PatchUserPreferencesRequest
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

func (h *AuthHandler) UpdateGoals(c fiber.Ctx) error {
	firebaseUID := c.Locals("firebase_uid").(string)

	var req domain.UpdateGoalRequest
	if err := c.Bind().JSON(&req); err != nil {
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

	ctx := c.Context()
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
