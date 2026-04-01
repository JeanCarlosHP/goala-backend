package services

import (
	"testing"

	"github.com/jeancarloshp/calorieai/internal/domain"
)

func TestNotificationPreferencesUpdateFromPatch(t *testing.T) {
	t.Parallel()

	globalEnabled := true
	streakEnabled := false
	reminderTime := "07:45"

	update, ok := notificationPreferencesUpdateFromPatch(domain.PatchUserPreferencesRequest{
		NotificationsEnabled: &globalEnabled,
		NotificationPrefs: &domain.PatchNotificationPreferencesRequest{
			DailyReminder: &domain.PatchDailyReminderPreference{
				Time: &reminderTime,
			},
			StreakRisk: &domain.PatchNotificationCategoryRequest{
				Enabled: &streakEnabled,
			},
		},
	})
	if !ok {
		t.Fatal("expected notification preference changes to be detected")
	}
	if update.NotificationsEnabled == nil || !*update.NotificationsEnabled {
		t.Fatal("expected global notifications update to be preserved")
	}
	if update.DailyReminderTime == nil || *update.DailyReminderTime != "07:45" {
		t.Fatal("expected reminder time update to be preserved")
	}
	if update.StreakRiskEnabled == nil || *update.StreakRiskEnabled {
		t.Fatal("expected streak risk update to be preserved")
	}
	if update.DailyReminderEnabled != nil {
		t.Fatal("expected omitted daily reminder enabled flag to remain nil")
	}
	if update.AchievementUnlockedEnabled != nil {
		t.Fatal("expected omitted achievement flag to remain nil")
	}
}

func TestNotificationPreferencesUpdateFromPatchNoChanges(t *testing.T) {
	t.Parallel()

	update, ok := notificationPreferencesUpdateFromPatch(domain.PatchUserPreferencesRequest{
		NotificationPrefs: &domain.PatchNotificationPreferencesRequest{},
	})
	if ok {
		t.Fatal("expected empty nested preferences to produce no update")
	}
	if update.NotificationsEnabled != nil || update.DailyReminderEnabled != nil || update.DailyReminderTime != nil || update.StreakRiskEnabled != nil || update.AchievementUnlockedEnabled != nil {
		t.Fatal("expected all update fields to remain nil")
	}
}
