package domain

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/go-playground/validator/v10"
)

func TestPatchUserPreferencesRequestValidation(t *testing.T) {
	t.Parallel()

	validate := validator.New()

	validReq := PatchUserPreferencesRequest{
		NotificationPrefs: &PatchNotificationPreferencesRequest{
			DailyReminder: &PatchDailyReminderPreference{},
		},
	}
	if err := validate.Struct(validReq); err != nil {
		t.Fatalf("expected omitted nested fields to validate, got %v", err)
	}

	invalidTime := "25:99"
	invalidReq := PatchUserPreferencesRequest{
		NotificationPrefs: &PatchNotificationPreferencesRequest{
			DailyReminder: &PatchDailyReminderPreference{
				Time: &invalidTime,
			},
		},
	}
	if err := validate.Struct(invalidReq); err == nil {
		t.Fatal("expected invalid reminder time to fail validation")
	}
}

func TestNotificationPreferencesEffective(t *testing.T) {
	t.Parallel()

	prefs := NotificationPreferences{
		DailyReminder: DailyReminderPreference{
			Enabled: true,
			Time:    "09:30",
		},
		StreakRisk: NotificationCategoryPreference{
			Enabled: true,
		},
		AchievementUnlocked: NotificationCategoryPreference{
			Enabled: true,
		},
	}

	effective := prefs.Effective(false)
	if effective.DailyReminder.Enabled {
		t.Fatal("expected daily reminder to be disabled when global notifications are off")
	}
	if effective.DailyReminder.Time != "09:30" {
		t.Fatalf("expected reminder time to be preserved, got %q", effective.DailyReminder.Time)
	}
	if effective.StreakRisk.Enabled {
		t.Fatal("expected streak risk to be disabled when global notifications are off")
	}
	if effective.AchievementUnlocked.Enabled {
		t.Fatal("expected achievement unlocked to be disabled when global notifications are off")
	}
}

func TestUserProfileResponseJSONShape(t *testing.T) {
	t.Parallel()

	profile := UserProfileResponse{
		ID:                   "user-1",
		DisplayName:          "Jean",
		Email:                "jean@example.com",
		DailyCalorieGoal:     2000,
		DailyProteinGoal:     150,
		Language:             "en-US",
		Timezone:             "America/Sao_Paulo",
		NotificationsEnabled: true,
		NotificationPrefs: NotificationPreferences{
			DailyReminder: DailyReminderPreference{
				Enabled: true,
				Time:    "08:15",
			},
			StreakRisk: NotificationCategoryPreference{
				Enabled: true,
			},
			AchievementUnlocked: NotificationCategoryPreference{
				Enabled: false,
			},
		},
		CreatedAt: time.Unix(0, 0).UTC(),
		UpdatedAt: time.Unix(0, 0).UTC(),
	}

	raw, err := json.Marshal(profile)
	if err != nil {
		t.Fatalf("marshal profile: %v", err)
	}

	var payload map[string]any
	if err := json.Unmarshal(raw, &payload); err != nil {
		t.Fatalf("unmarshal profile json: %v", err)
	}

	if _, ok := payload["notificationsEnabled"]; !ok {
		t.Fatal("expected legacy notificationsEnabled key in profile response")
	}

	prefs, ok := payload["notificationPreferences"].(map[string]any)
	if !ok {
		t.Fatal("expected notificationPreferences object in profile response")
	}
	if _, ok := prefs["dailyReminder"].(map[string]any); !ok {
		t.Fatal("expected dailyReminder object in notificationPreferences")
	}
	if _, ok := prefs["streakRisk"].(map[string]any); !ok {
		t.Fatal("expected streakRisk object in notificationPreferences")
	}
	if _, ok := prefs["achievementUnlocked"].(map[string]any); !ok {
		t.Fatal("expected achievementUnlocked object in notificationPreferences")
	}
}
