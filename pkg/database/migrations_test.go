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
		"COALESCE(daily_reminder_enabled, notifications_enabled, false)",
		"COALESCE(streak_risk_enabled, notifications_enabled, false)",
		"COALESCE(achievement_unlocked_enabled, notifications_enabled, false)",
		"COALESCE(daily_reminder_time, '20:00')",
		"ALTER COLUMN daily_reminder_time SET DEFAULT '20:00'",
		"users_daily_reminder_time_format_chk",
	}

	for _, snippet := range expectedSnippets {
		if !strings.Contains(sql, snippet) {
			t.Fatalf("expected migration to contain %q", snippet)
		}
	}
}
