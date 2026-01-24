-- name: CreateFeedback :one
INSERT INTO feedback (
    user_id, type, title, description, user_email,
    platform, os_version, app_version
) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
RETURNING *;

-- name: ListFeedback :many
SELECT * FROM feedback
ORDER BY created_at DESC
LIMIT $1 OFFSET $2;

-- name: GetFeedback :one
SELECT * FROM feedback WHERE id = $1;

-- name: UpdateFeedbackStatus :exec
UPDATE feedback SET status = $2 WHERE id = $1;

-- name: GetFeedbackByUser :many
SELECT * FROM feedback
WHERE user_id = $1
ORDER BY created_at DESC;

-- name: GetFeedbackByType :many
SELECT * FROM feedback
WHERE type = $1
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;
