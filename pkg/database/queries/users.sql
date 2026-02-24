-- name: CreateUser :one
INSERT INTO users (id, firebase_uid, email, display_name, photo_url)
VALUES ($1, $2, $3, $4, $5)
RETURNING id, firebase_uid, email, display_name, photo_url, created_at, updated_at;

-- name: GetUserByFirebaseUID :one
SELECT id, firebase_uid, email, display_name, photo_url, created_at, updated_at,
       weight, height, age, gender, activity_level, language, notifications_enabled, timezone
FROM users
WHERE firebase_uid = $1;

-- name: GetUserByID :one
SELECT id, firebase_uid, email, display_name, photo_url, created_at, updated_at,
       weight, height, age, gender, activity_level, language, notifications_enabled, timezone
FROM users
WHERE id = $1;

-- name: UpdateUser :one
UPDATE users
SET email = $2, display_name = $3, photo_url = $4, updated_at = NOW()
WHERE id = $1
RETURNING id, firebase_uid, email, display_name, photo_url, created_at, updated_at,
          weight, height, age, gender, activity_level, language, notifications_enabled, timezone;

-- name: UpdateUserProfile :exec
UPDATE users SET
    display_name = $2,
    email = $3,
    photo_url = $4,
    weight = $5,
    height = $6,
    age = $7,
    gender = $8,
    activity_level = $9,
    language = $10,
    notifications_enabled = $11,
    timezone = $12,
    updated_at = NOW()
WHERE id = $1;

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

-- name: UpdateUserNotifications :exec
UPDATE users
SET notifications_enabled = $2, updated_at = NOW()
WHERE id = $1;
