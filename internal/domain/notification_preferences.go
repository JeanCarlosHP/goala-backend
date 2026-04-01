package domain

import "regexp"

const DefaultDailyReminderTime = "09:00"

var reminderTimePattern = regexp.MustCompile(`^([01]\d|2[0-3]):[0-5]\d$`)

type NotificationPreference struct {
	Enabled bool `json:"enabled"`
}

type DailyReminderPreference struct {
	Enabled bool   `json:"enabled"`
	Time    string `json:"time"`
}

type NotificationPreferences struct {
	DailyReminder       DailyReminderPreference `json:"dailyReminder"`
	StreakAtRisk        NotificationPreference  `json:"streakAtRisk"`
	AchievementUnlocked NotificationPreference  `json:"achievementUnlocked"`
}

type PatchNotificationPreference struct {
	Enabled *bool `json:"enabled" validate:"omitempty"`
}

type PatchDailyReminderPreference struct {
	Enabled *bool   `json:"enabled" validate:"omitempty"`
	Time    *string `json:"time" validate:"omitempty,notification_time"`
}

type PatchNotificationPreferencesRequest struct {
	DailyReminder       *PatchDailyReminderPreference `json:"dailyReminder" validate:"omitempty"`
	StreakAtRisk        *PatchNotificationPreference  `json:"streakAtRisk" validate:"omitempty"`
	AchievementUnlocked *PatchNotificationPreference  `json:"achievementUnlocked" validate:"omitempty"`
}

type NotificationPreferencesUpdate struct {
	NotificationsEnabled       *bool
	DailyReminderEnabled       *bool
	DailyReminderTime          *string
	StreakAtRiskEnabled        *bool
	AchievementUnlockedEnabled *bool
}

func DefaultNotificationPreferences(enabled bool) NotificationPreferences {
	return NotificationPreferences{
		DailyReminder: DailyReminderPreference{
			Enabled: enabled,
			Time:    DefaultDailyReminderTime,
		},
		StreakAtRisk: NotificationPreference{
			Enabled: enabled,
		},
		AchievementUnlocked: NotificationPreference{
			Enabled: enabled,
		},
	}
}

func IsValidReminderTime(value string) bool {
	return reminderTimePattern.MatchString(value)
}

func (r PatchUserPreferencesRequest) NotificationPreferencesUpdate() NotificationPreferencesUpdate {
	update := NotificationPreferencesUpdate{
		NotificationsEnabled: r.NotificationsEnabled,
	}

	if r.NotificationPreferences == nil {
		return update
	}

	if r.NotificationPreferences.DailyReminder != nil {
		update.DailyReminderEnabled = r.NotificationPreferences.DailyReminder.Enabled
		update.DailyReminderTime = r.NotificationPreferences.DailyReminder.Time
	}

	if r.NotificationPreferences.StreakAtRisk != nil {
		update.StreakAtRiskEnabled = r.NotificationPreferences.StreakAtRisk.Enabled
	}

	if r.NotificationPreferences.AchievementUnlocked != nil {
		update.AchievementUnlockedEnabled = r.NotificationPreferences.AchievementUnlocked.Enabled
	}

	return update
}
