-- name: UpsertUserGoal :one
INSERT INTO user_goals (user_id, daily_calories, protein_g, carbs_g, fat_g)
VALUES ($1, $2, $3, $4, $5)
ON CONFLICT (user_id)
DO UPDATE SET
    daily_calories = EXCLUDED.daily_calories,
    protein_g = EXCLUDED.protein_g,
    carbs_g = EXCLUDED.carbs_g,
    fat_g = EXCLUDED.fat_g,
    updated_at = NOW()
RETURNING user_id, daily_calories, protein_g, carbs_g, fat_g, updated_at;

-- name: GetUserGoalByUserID :one
SELECT user_id, daily_calories, protein_g, carbs_g, fat_g, updated_at
FROM user_goals
WHERE user_id = $1;
