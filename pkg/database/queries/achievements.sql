-- name: GetAllAchievements :many
SELECT * FROM achievements ORDER BY category, target;

-- name: GetAchievementByID :one
SELECT * FROM achievements WHERE id = $1;

-- name: GetAchievementsByCategory :many
SELECT * FROM achievements WHERE category = $1 ORDER BY target;

-- name: GetUserAchievements :many
SELECT ua.*, a.name_key, a.description_key, a.icon, a.target, a.category
FROM user_achievements ua
JOIN achievements a ON ua.achievement_id = a.id
WHERE ua.user_id = $1
ORDER BY ua.unlocked DESC, a.category, a.target;

-- name: UpsertUserAchievement :exec
INSERT INTO user_achievements (user_id, achievement_id, unlocked, progress, unlocked_at)
VALUES ($1, $2, $3, $4, $5)
ON CONFLICT (user_id, achievement_id) DO UPDATE SET
    unlocked = EXCLUDED.unlocked,
    progress = EXCLUDED.progress,
    unlocked_at = EXCLUDED.unlocked_at,
    updated_at = NOW();

-- name: GetUserAchievement :one
SELECT * FROM user_achievements
WHERE user_id = $1 AND achievement_id = $2;

-- name: GetUserUnlockedAchievements :many
SELECT ua.*, a.name_key, a.description_key, a.icon, a.target, a.category
FROM user_achievements ua
JOIN achievements a ON ua.achievement_id = a.id
WHERE ua.user_id = $1 AND ua.unlocked = true
ORDER BY ua.unlocked_at DESC;

-- name: UpdateAchievementProgress :exec
UPDATE user_achievements SET
    progress = $3,
    updated_at = NOW()
WHERE user_id = $1 AND achievement_id = $2;
