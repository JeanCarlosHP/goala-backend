ALTER TABLE users
    DROP CONSTRAINT IF EXISTS users_notification_daily_reminder_time_check;

ALTER TABLE users
    DROP COLUMN IF EXISTS notification_daily_reminder_enabled,
    DROP COLUMN IF EXISTS notification_daily_reminder_time,
    DROP COLUMN IF EXISTS notification_streak_at_risk_enabled,
    DROP COLUMN IF EXISTS notification_achievement_unlocked_enabled;
