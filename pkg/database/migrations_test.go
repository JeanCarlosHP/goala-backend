package database

import (
	"os"
	"strings"
	"testing"
)

func TestNotificationPreferencesMigrationDefaults(t *testing.T) {
	t.Parallel()

	raw, err := os.ReadFile("migrations/000015_add_notification_preferences_v1.up.sql")
	if err != nil {
		t.Fatalf("read migration: %v", err)
	}

	sql := string(raw)
	expectedSnippets := []string{
		"COALESCE(notification_daily_reminder_enabled, COALESCE(notifications_enabled, false))",
		"COALESCE(notification_streak_at_risk_enabled, COALESCE(notifications_enabled, false))",
		"COALESCE(notification_achievement_unlocked_enabled, COALESCE(notifications_enabled, false))",
		"COALESCE(notification_daily_reminder_time, '09:00')",
		"ALTER COLUMN notification_daily_reminder_time SET DEFAULT '09:00'",
		"users_notification_daily_reminder_time_check",
	}

	for _, snippet := range expectedSnippets {
		if !strings.Contains(sql, snippet) {
			t.Fatalf("expected migration to contain %q", snippet)
		}
	}
}
