ALTER TABLE users
    ADD COLUMN IF NOT EXISTS notification_daily_reminder_enabled BOOLEAN,
    ADD COLUMN IF NOT EXISTS notification_daily_reminder_time VARCHAR(5),
    ADD COLUMN IF NOT EXISTS notification_streak_at_risk_enabled BOOLEAN,
    ADD COLUMN IF NOT EXISTS notification_achievement_unlocked_enabled BOOLEAN;

UPDATE users
SET
    notification_daily_reminder_enabled = COALESCE(notification_daily_reminder_enabled, COALESCE(notifications_enabled, false)),
    notification_daily_reminder_time = COALESCE(notification_daily_reminder_time, '09:00'),
    notification_streak_at_risk_enabled = COALESCE(notification_streak_at_risk_enabled, COALESCE(notifications_enabled, false)),
    notification_achievement_unlocked_enabled = COALESCE(notification_achievement_unlocked_enabled, COALESCE(notifications_enabled, false));

ALTER TABLE users
    ALTER COLUMN notification_daily_reminder_enabled SET DEFAULT false,
    ALTER COLUMN notification_daily_reminder_enabled SET NOT NULL,
    ALTER COLUMN notification_daily_reminder_time SET DEFAULT '09:00',
    ALTER COLUMN notification_daily_reminder_time SET NOT NULL,
    ALTER COLUMN notification_streak_at_risk_enabled SET DEFAULT false,
    ALTER COLUMN notification_streak_at_risk_enabled SET NOT NULL,
    ALTER COLUMN notification_achievement_unlocked_enabled SET DEFAULT false,
    ALTER COLUMN notification_achievement_unlocked_enabled SET NOT NULL;

ALTER TABLE users
    ADD CONSTRAINT users_notification_daily_reminder_time_check
    CHECK (notification_daily_reminder_time ~ '^([01][0-9]|2[0-3]):[0-5][0-9]$');
