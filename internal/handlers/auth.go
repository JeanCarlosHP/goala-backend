package handlers

import (
	"context"

	"github.com/jeancarloshp/calorieai/internal/services"

	firebase "firebase.google.com/go/v4"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/jeancarloshp/calorieai/internal/domain"
	firebaseApp "github.com/jeancarloshp/calorieai/pkg/firebase"
)

type AuthHandler struct {
	userService *services.UserService
	firebaseApp *firebase.App
	validator   *validator.Validate
	logger      domain.Logger
}

func NewAuthHandler(userService *services.UserService, firebaseApp *firebase.App, logger domain.Logger) *AuthHandler {
	return &AuthHandler{
		userService: userService,
		firebaseApp: firebaseApp,
		validator:   validator.New(),
		logger:      logger,
	}
}

func (h *AuthHandler) Register(c *fiber.Ctx) error {
	var req domain.RegisterRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid request body",
		})
	}

	ctx := context.Background()
	authClient, err := firebaseApp.GetAuthClient(ctx, h.firebaseApp)
	if err != nil {
		h.logger.Error("Failed to get auth client", "error", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "authentication service unavailable",
		})
	}

	userRecord, err := authClient.GetUser(ctx, req.FirebaseUID)
	if err != nil {
		h.logger.Error("Invalid Firebase UID", "firebase_uid", req.FirebaseUID, "error", err)
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "invalid firebase user",
		})
	}

	if req.Email == "" {
		req.Email = userRecord.Email
	}
	if req.DisplayName == "" {
		req.DisplayName = userRecord.DisplayName
	}
	if req.PhotoURL == nil && userRecord.PhotoURL != "" {
		req.PhotoURL = &userRecord.PhotoURL
	}

	user, err := h.userService.RegisterUser(ctx, req)
	if err != nil {
		if err == domain.ErrUserAlreadyExists {
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{
				"error": "user already exists",
			})
		}
		h.logger.Error("Failed to register user", "error", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to register user",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(user)
}

func (h *AuthHandler) GetMe(c *fiber.Ctx) error {
	firebaseUID := c.Locals("firebase_uid").(string)

	ctx := context.Background()
	user, err := h.userService.GetUserByFirebaseUID(ctx, firebaseUID)
	if err != nil {
		h.logger.Error("User not found", "firebase_uid", firebaseUID, "error", err)
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "user not found",
		})
	}

	goal, err := h.userService.GetUserGoal(ctx, user.ID)
	if err != nil {
		h.logger.Warn("Goal not found", "user_id", user.ID.String(), "error", err)
	}

	return c.JSON(fiber.Map{
		"user": user,
		"goal": goal,
	})
}
