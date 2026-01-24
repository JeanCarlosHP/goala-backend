-- name: CreateMeal :one
INSERT INTO meals (id, user_id, meal_type, meal_date, meal_time, photo_url)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING id, user_id, meal_type, meal_date, meal_time, photo_url, created_at;

-- name: GetMealsByUserAndDate :many
SELECT id, user_id, meal_type, meal_date, meal_time, photo_url, created_at
FROM meals
WHERE user_id = $1 AND meal_date = $2
ORDER BY meal_time ASC, created_at ASC;

-- name: GetMealByID :one
SELECT id, user_id, meal_type, meal_date, meal_time, photo_url, created_at
FROM meals
WHERE id = $1;
