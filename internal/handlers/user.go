package handlers

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v3"
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
	userID, ok := c.Locals("user_id").(uuid.UUID)
	if !ok {
		h.logger.Warn("Missing user_id in context")
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"message": "authentication required",
		})
	}

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

	uploadURL, photoURL, err := h.s3Service.GenerateUploadPresignedURL(ctx, userID.String(), req.ContentType, req.FileSize)
	if err != nil {
		h.logger.Error("Failed to generate presigned URL", "user_id", userID.String(), "error", err)
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

func (h *UserHandler) GetAvatar(c fiber.Ctx) error {
	userID := strings.TrimSpace(c.Params("userID"))
	filename := strings.TrimSpace(c.Params("filename"))
	if userID == "" || filename == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "invalid avatar path",
		})
	}

	objectPath := fmt.Sprintf("/users/%s/avatars/%s", userID, filename)
	objectData, err := h.s3Service.GetObject(c.Context(), objectPath)
	if err != nil {
		h.logger.Error("Failed to fetch avatar", "firebase_uid", userID, "path", objectPath, "error", err)
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"message": "image not found",
		})
	}

	c.Set(fiber.HeaderContentType, objectData.ContentType)
	c.Set(fiber.HeaderCacheControl, "public, max-age=300")
	return c.Send(objectData.Body)
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

	fmt.Print(req)

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
