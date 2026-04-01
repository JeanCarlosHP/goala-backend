-- name: CreateUser :one
INSERT INTO users (id, firebase_uid, email, display_name, photo_url)
VALUES ($1, $2, $3, $4, $5)
RETURNING id, firebase_uid, email, display_name, photo_url, created_at, updated_at;

-- name: GetUserByFirebaseUID :one
SELECT id, firebase_uid, email, display_name, photo_url, created_at, updated_at,
       weight, height, age, gender, activity_level, language, notifications_enabled, timezone,
       daily_reminder_enabled, daily_reminder_time, streak_risk_enabled, achievement_unlocked_enabled
FROM users
WHERE firebase_uid = $1;

-- name: GetUserByID :one
SELECT id, firebase_uid, email, display_name, photo_url, created_at, updated_at,
       weight, height, age, gender, activity_level, language, notifications_enabled, timezone,
       daily_reminder_enabled, daily_reminder_time, streak_risk_enabled, achievement_unlocked_enabled
FROM users
WHERE id = $1;

-- name: UpdateUser :one
UPDATE users
SET email = $2, display_name = $3, photo_url = $4, updated_at = NOW()
WHERE id = $1
RETURNING id, firebase_uid, email, display_name, photo_url, created_at, updated_at,
          weight, height, age, gender, activity_level, language, notifications_enabled, timezone,
          daily_reminder_enabled, daily_reminder_time, streak_risk_enabled, achievement_unlocked_enabled;

-- name: UpdateUserProfile :exec
UPDATE users SET
    display_name = sqlc.arg(display_name),
    email = sqlc.arg(email),
    photo_url = sqlc.arg(photo_url),
    weight = sqlc.arg(weight),
    height = sqlc.arg(height),
    age = sqlc.arg(age),
    gender = sqlc.arg(gender),
    activity_level = sqlc.arg(activity_level),
    language = sqlc.arg(language),
    notifications_enabled = sqlc.arg(notifications_enabled),
    timezone = sqlc.arg(timezone),
    updated_at = NOW()
WHERE id = sqlc.arg(id);

-- name: ExistsUserByFirebaseUID :one
SELECT EXISTS(SELECT 1 FROM users WHERE firebase_uid = $1);

-- name: UpdateUserAvatar :exec
UPDATE users
SET photo_url = $2, updated_at = NOW()
WHERE id = $1;

-- name: UpdateUserDisplayName :exec
UPDATE users
SET display_name = $2, updated_at = NOW()
WHERE id = $1;

-- name: UpdateUserNotificationPreferences :exec
UPDATE users
SET
    notifications_enabled = COALESCE(sqlc.narg(notifications_enabled), notifications_enabled),
    daily_reminder_enabled = COALESCE(sqlc.narg(daily_reminder_enabled), daily_reminder_enabled),
    daily_reminder_time = COALESCE(sqlc.narg(daily_reminder_time), daily_reminder_time),
    streak_risk_enabled = COALESCE(sqlc.narg(streak_risk_enabled), streak_risk_enabled),
    achievement_unlocked_enabled = COALESCE(sqlc.narg(achievement_unlocked_enabled), achievement_unlocked_enabled),
    updated_at = NOW()
WHERE id = sqlc.arg(id);
