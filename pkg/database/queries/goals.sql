-- name: UpsertUserGoal :one
INSERT INTO user_goals (user_id, daily_calories, protein, carbs, fat)
VALUES ($1, $2, $3, $4, $5)
ON CONFLICT (user_id)
DO UPDATE SET
    daily_calories = EXCLUDED.daily_calories,
    protein = EXCLUDED.protein,
    carbs = EXCLUDED.carbs,
    fat = EXCLUDED.fat,
    updated_at = NOW()
RETURNING user_id, daily_calories, protein, carbs, fat, updated_at;

-- name: GetUserGoalByUserID :one
SELECT user_id, daily_calories, protein, carbs, fat, updated_at
FROM user_goals
WHERE user_id = $1;
