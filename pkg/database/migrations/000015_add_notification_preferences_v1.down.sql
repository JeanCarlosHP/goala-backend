ALTER TABLE users
    DROP CONSTRAINT IF EXISTS users_daily_reminder_time_format_chk;

ALTER TABLE users
    DROP COLUMN IF EXISTS achievement_unlocked_enabled,
    DROP COLUMN IF EXISTS streak_risk_enabled,
    DROP COLUMN IF EXISTS daily_reminder_time,
    DROP COLUMN IF EXISTS daily_reminder_enabled;
