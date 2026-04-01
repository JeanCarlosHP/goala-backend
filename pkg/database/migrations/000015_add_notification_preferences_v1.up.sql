ALTER TABLE users
    ADD COLUMN IF NOT EXISTS daily_reminder_enabled BOOLEAN,
    ADD COLUMN IF NOT EXISTS daily_reminder_time VARCHAR(5),
    ADD COLUMN IF NOT EXISTS streak_risk_enabled BOOLEAN,
    ADD COLUMN IF NOT EXISTS achievement_unlocked_enabled BOOLEAN;

UPDATE users
SET
    daily_reminder_enabled = COALESCE(daily_reminder_enabled, notifications_enabled, false),
    daily_reminder_time = COALESCE(daily_reminder_time, '20:00'),
    streak_risk_enabled = COALESCE(streak_risk_enabled, notifications_enabled, false),
    achievement_unlocked_enabled = COALESCE(achievement_unlocked_enabled, notifications_enabled, false);

ALTER TABLE users
    ALTER COLUMN daily_reminder_enabled SET DEFAULT false,
    ALTER COLUMN daily_reminder_enabled SET NOT NULL,
    ALTER COLUMN daily_reminder_time SET DEFAULT '20:00',
    ALTER COLUMN daily_reminder_time SET NOT NULL,
    ALTER COLUMN streak_risk_enabled SET DEFAULT false,
    ALTER COLUMN streak_risk_enabled SET NOT NULL,
    ALTER COLUMN achievement_unlocked_enabled SET DEFAULT false,
    ALTER COLUMN achievement_unlocked_enabled SET NOT NULL;

ALTER TABLE users
    ADD CONSTRAINT users_daily_reminder_time_format_chk
    CHECK (daily_reminder_time ~ '^(?:[01][0-9]|2[0-3]):[0-5][0-9]$');
