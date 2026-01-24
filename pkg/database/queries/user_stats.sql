-- name: CreateUserStats :one
INSERT INTO user_stats (user_id)
VALUES ($1)
RETURNING *;

-- name: GetUserStats :one
SELECT * FROM user_stats WHERE user_id = $1;

-- name: UpdateUserStats :exec
UPDATE user_stats SET
    current_streak = $2,
    longest_streak = $3,
    total_meals_logged = $4,
    total_days_logged = $5,
    total_calories_consumed = $6,
    total_protein_consumed = $7,
    total_carbs_consumed = $8,
    total_fat_consumed = $9,
    last_log_date = $10,
    updated_at = NOW()
WHERE user_id = $1;

-- name: IncrementMealCount :exec
UPDATE user_stats SET
    total_meals_logged = total_meals_logged + 1,
    updated_at = NOW()
WHERE user_id = $1;

-- name: UpdateStreakAndLastLogDate :exec
UPDATE user_stats SET
    current_streak = $2,
    longest_streak = GREATEST(longest_streak, $2),
    last_log_date = $3,
    updated_at = NOW()
WHERE user_id = $1;

-- name: AddNutritionToStats :exec
UPDATE user_stats SET
    total_calories_consumed = total_calories_consumed + $2,
    total_protein_consumed = total_protein_consumed + $3,
    total_carbs_consumed = total_carbs_consumed + $4,
    total_fat_consumed = total_fat_consumed + $5,
    updated_at = NOW()
WHERE user_id = $1;
