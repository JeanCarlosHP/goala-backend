package services

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jeancarloshp/calorieai/internal/domain"
)

type stubUserRepository struct {
	user        *domain.User
	lastUpdate  domain.NotificationPreferencesUpdate
	updateCalls int
}

func (s *stubUserRepository) Create(_ context.Context, user *domain.User) error {
	s.user = user
	return nil
}

func (s *stubUserRepository) ExistsByFirebaseUID(_ context.Context, _ string) (bool, error) {
	return false, nil
}

func (s *stubUserRepository) GetByFirebaseUID(_ context.Context, _ string) (*domain.User, error) {
	return s.user, nil
}

func (s *stubUserRepository) GetByID(_ context.Context, _ uuid.UUID) (*domain.User, error) {
	return s.user, nil
}

func (s *stubUserRepository) UpdateProfile(_ context.Context, user *domain.User) error {
	s.user = user
	return nil
}

func (s *stubUserRepository) UpdateAvatar(_ context.Context, _ uuid.UUID, photoURL *string) error {
	s.user.PhotoURL = photoURL
	return nil
}

func (s *stubUserRepository) UpdateDisplayName(_ context.Context, _ uuid.UUID, displayName *string) error {
	if displayName != nil {
		s.user.DisplayName = *displayName
	}
	return nil
}

func (s *stubUserRepository) UpdateNotificationPreferences(_ context.Context, _ uuid.UUID, update domain.NotificationPreferencesUpdate) error {
	s.lastUpdate = update
	s.updateCalls++

	if update.NotificationsEnabled != nil {
		s.user.NotificationsEnabled = *update.NotificationsEnabled
	}
	if update.DailyReminderEnabled != nil {
		s.user.NotificationPreferences.DailyReminder.Enabled = *update.DailyReminderEnabled
	}
	if update.DailyReminderTime != nil {
		s.user.NotificationPreferences.DailyReminder.Time = *update.DailyReminderTime
	}
	if update.StreakAtRiskEnabled != nil {
		s.user.NotificationPreferences.StreakAtRisk.Enabled = *update.StreakAtRiskEnabled
	}
	if update.AchievementUnlockedEnabled != nil {
		s.user.NotificationPreferences.AchievementUnlocked.Enabled = *update.AchievementUnlockedEnabled
	}

	return nil
}

type stubGoalRepository struct {
	goal *domain.UserGoal
}

func (s *stubGoalRepository) GetByUserID(_ context.Context, _ uuid.UUID) (*domain.UserGoal, error) {
	return s.goal, nil
}

func (s *stubGoalRepository) Upsert(_ context.Context, goal *domain.UserGoal) error {
	s.goal = goal
	return nil
}

func TestGetUserProfileIncludesNotificationPreferences(t *testing.T) {
	t.Parallel()

	userID := uuid.New()
	createdAt := time.Date(2026, 4, 1, 12, 0, 0, 0, time.UTC)
	updatedAt := createdAt.Add(10 * time.Minute)
	userRepo := &stubUserRepository{
		user: &domain.User{
			ID:                   userID,
			Email:                "user@example.com",
			DisplayName:          "Jean",
			PhotoURL:             stringRef("/avatars/user.jpg"),
			Language:             "en-US",
			Timezone:             "America/Sao_Paulo",
			NotificationsEnabled: true,
			NotificationPreferences: domain.NotificationPreferences{
				DailyReminder: domain.DailyReminderPreference{
					Enabled: true,
					Time:    "08:45",
				},
				StreakAtRisk: domain.NotificationPreference{
					Enabled: false,
				},
				AchievementUnlocked: domain.NotificationPreference{
					Enabled: true,
				},
			},
			CreatedAt: createdAt,
			UpdatedAt: updatedAt,
		},
	}
	goalRepo := &stubGoalRepository{
		goal: &domain.UserGoal{
			UserID:           userID,
			DailyCalorieGoal: 2100,
			DailyProteinGoal: 160,
			DailyCarbsGoal:   220,
			DailyFatGoal:     70,
		},
	}

	service := NewUserService(userRepo, goalRepo, "https://cdn.example.com")

	profile, err := service.GetUserProfile(context.Background(), userID)
	if err != nil {
		t.Fatalf("GetUserProfile returned error: %v", err)
	}

	if profile.NotificationPreferences.DailyReminder.Time != "08:45" {
		t.Fatalf("expected daily reminder time to be preserved, got %q", profile.NotificationPreferences.DailyReminder.Time)
	}

	if profile.NotificationPreferences.StreakAtRisk.Enabled {
		t.Fatalf("expected streak-at-risk preference to remain false")
	}

	if profile.Photo == nil || *profile.Photo != "https://cdn.example.com/avatars/user.jpg" {
		t.Fatalf("expected CDN photo URL, got %#v", profile.Photo)
	}
}

func TestPatchUserPreferencesOnlyUpdatesTargetedNotificationFields(t *testing.T) {
	t.Parallel()

	userID := uuid.New()
	userRepo := &stubUserRepository{
		user: &domain.User{
			ID:                   userID,
			Email:                "user@example.com",
			DisplayName:          "Jean",
			Language:             "en-US",
			Timezone:             "America/Sao_Paulo",
			NotificationsEnabled: true,
			NotificationPreferences: domain.NotificationPreferences{
				DailyReminder: domain.DailyReminderPreference{
					Enabled: true,
					Time:    "09:00",
				},
				StreakAtRisk: domain.NotificationPreference{
					Enabled: true,
				},
				AchievementUnlocked: domain.NotificationPreference{
					Enabled: false,
				},
			},
		},
	}
	goalRepo := &stubGoalRepository{
		goal: &domain.UserGoal{
			UserID:           userID,
			DailyCalorieGoal: 2000,
			DailyProteinGoal: 150,
			DailyCarbsGoal:   200,
			DailyFatGoal:     65,
		},
	}

	service := NewUserService(userRepo, goalRepo, "https://cdn.example.com")
	req := domain.PatchUserPreferencesRequest{
		NotificationPreferences: &domain.PatchNotificationPreferencesRequest{
			DailyReminder: &domain.PatchDailyReminderPreference{
				Time: stringRef("08:30"),
			},
		},
	}

	profile, err := service.PatchUserPreferences(context.Background(), userID, req)
	if err != nil {
		t.Fatalf("PatchUserPreferences returned error: %v", err)
	}

	if userRepo.updateCalls != 1 {
		t.Fatalf("expected one notification update call, got %d", userRepo.updateCalls)
	}

	if userRepo.lastUpdate.DailyReminderTime == nil || *userRepo.lastUpdate.DailyReminderTime != "08:30" {
		t.Fatalf("expected daily reminder time update to be captured, got %#v", userRepo.lastUpdate.DailyReminderTime)
	}

	if userRepo.lastUpdate.NotificationsEnabled != nil {
		t.Fatalf("expected global notifications flag to remain untouched")
	}

	if userRepo.lastUpdate.StreakAtRiskEnabled != nil || userRepo.lastUpdate.AchievementUnlockedEnabled != nil {
		t.Fatalf("expected unrelated notification categories to remain untouched")
	}

	if profile.NotificationPreferences.DailyReminder.Time != "08:30" {
		t.Fatalf("expected returned profile to include updated daily reminder time, got %q", profile.NotificationPreferences.DailyReminder.Time)
	}

	if !profile.NotificationPreferences.StreakAtRisk.Enabled {
		t.Fatalf("expected streak-at-risk preference to remain true")
	}

	if profile.NotificationPreferences.AchievementUnlocked.Enabled {
		t.Fatalf("expected achievement-unlocked preference to remain false")
	}
}

func stringRef(value string) *string {
	return &value
}
