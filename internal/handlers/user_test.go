package handlers

import (
	"context"
	"encoding/json"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
	"github.com/jeancarloshp/calorieai/internal/domain"
)

type stubHandlerUserService struct {
	user        *domain.User
	profile     *domain.UserProfileResponse
	patchCalled bool
}

func (s *stubHandlerUserService) GetUserByFirebaseUID(_ context.Context, _ string) (*domain.User, error) {
	return s.user, nil
}

func (s *stubHandlerUserService) GetUserProfile(_ context.Context, _ uuid.UUID) (*domain.UserProfileResponse, error) {
	return s.profile, nil
}

func (s *stubHandlerUserService) UpdateUserProfile(_ context.Context, _ uuid.UUID, _ domain.UpdateProfileRequest) (*domain.UserProfileResponse, error) {
	return s.profile, nil
}

func (s *stubHandlerUserService) PatchUserPreferences(_ context.Context, _ uuid.UUID, _ domain.PatchUserPreferencesRequest) (*domain.UserProfileResponse, error) {
	s.patchCalled = true
	return s.profile, nil
}

type stubAvatarUploadService struct{}

func (s *stubAvatarUploadService) GenerateUploadPresignedURL(_ context.Context, _, _ string, _ int64) (string, string, error) {
	return "", "", nil
}

type stubLogger struct{}

func (l stubLogger) Panic(v any)           {}
func (l stubLogger) Panicf(string, ...any) {}
func (l stubLogger) Fatal(...any)          {}
func (l stubLogger) Fatalf(string, ...any) {}
func (l stubLogger) Error(string, ...any)  {}
func (l stubLogger) Warn(string, ...any)   {}
func (l stubLogger) Info(string, ...any)   {}
func (l stubLogger) Infof(string, ...any)  {}
func (l stubLogger) Debug(string, ...any)  {}

func TestPatchUserPreferencesRejectsInvalidReminderTime(t *testing.T) {
	t.Parallel()

	service := &stubHandlerUserService{
		user: &domain.User{ID: uuid.New()},
	}
	handler := NewUserHandler(service, &stubAvatarUploadService{}, stubLogger{})

	app := fiber.New()
	app.Patch("/api/v1/user/profile", func(c fiber.Ctx) error {
		c.Locals("firebase_uid", "firebase-123")
		return handler.PatchUserPreferences(c)
	})

	req := httptest.NewRequest("PATCH", "/api/v1/user/profile", strings.NewReader(`{"notificationPreferences":{"dailyReminder":{"time":"25:99"}}}`))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("app.Test returned error: %v", err)
	}

	if resp.StatusCode != fiber.StatusBadRequest {
		t.Fatalf("expected 400 status, got %d", resp.StatusCode)
	}

	if service.patchCalled {
		t.Fatalf("expected patch service not to be called for invalid payload")
	}
}

func TestGetProfileReturnsNotificationPreferencesShape(t *testing.T) {
	t.Parallel()

	userID := uuid.New()
	service := &stubHandlerUserService{
		user: &domain.User{ID: userID},
		profile: &domain.UserProfileResponse{
			ID:                   userID.String(),
			DisplayName:          "Jean",
			Email:                "user@example.com",
			Language:             "en-US",
			Timezone:             "America/Sao_Paulo",
			NotificationsEnabled: true,
			NotificationPreferences: domain.NotificationPreferences{
				DailyReminder: domain.DailyReminderPreference{
					Enabled: true,
					Time:    "07:30",
				},
				StreakAtRisk:        domain.NotificationPreference{Enabled: true},
				AchievementUnlocked: domain.NotificationPreference{Enabled: false},
			},
		},
	}
	handler := NewUserHandler(service, &stubAvatarUploadService{}, stubLogger{})

	app := fiber.New()
	app.Get("/api/v1/user/profile", func(c fiber.Ctx) error {
		c.Locals("firebase_uid", "firebase-123")
		return handler.GetProfile(c)
	})

	req := httptest.NewRequest("GET", "/api/v1/user/profile", nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("app.Test returned error: %v", err)
	}

	var body struct {
		Success bool `json:"success"`
		Data    struct {
			NotificationPreferences domain.NotificationPreferences `json:"notificationPreferences"`
		} `json:"data"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		t.Fatalf("failed to decode response body: %v", err)
	}

	if !body.Success {
		t.Fatalf("expected success response")
	}

	if body.Data.NotificationPreferences.DailyReminder.Time != "07:30" {
		t.Fatalf("expected notification preferences shape in response, got %+v", body.Data.NotificationPreferences)
	}

	if body.Data.NotificationPreferences.AchievementUnlocked.Enabled {
		t.Fatalf("expected achievement-unlocked flag to remain false")
	}
}
