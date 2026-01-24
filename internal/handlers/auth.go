package handlers

import (
	"context"

	"github.com/jeancarloshp/calorieai/internal/services"

	firebase "firebase.google.com/go/v4"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/jeancarloshp/calorieai/internal/domain"
	firebaseApp "github.com/jeancarloshp/calorieai/pkg/firebase"
	"github.com/rs/zerolog/log"
)

type AuthHandler struct {
	userService *services.UserService
	firebaseApp *firebase.App
	validator   *validator.Validate
}

func NewAuthHandler(userService *services.UserService, firebaseApp *firebase.App) *AuthHandler {
	return &AuthHandler{
		userService: userService,
		firebaseApp: firebaseApp,
		validator:   validator.New(),
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
		log.Error().Err(err).Msg("Failed to get auth client")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "authentication service unavailable",
		})
	}

	userRecord, err := authClient.GetUser(ctx, req.FirebaseUID)
	if err != nil {
		log.Error().Err(err).Str("firebase_uid", req.FirebaseUID).Msg("Invalid Firebase UID")
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
		log.Error().Err(err).Msg("Failed to register user")
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
		log.Error().Err(err).Str("firebase_uid", firebaseUID).Msg("User not found")
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "user not found",
		})
	}

	goal, err := h.userService.GetUserGoal(ctx, user.ID)
	if err != nil {
		log.Warn().Err(err).Str("user_id", user.ID.String()).Msg("Goal not found")
	}

	return c.JSON(fiber.Map{
		"user": user,
		"goal": goal,
	})
}
