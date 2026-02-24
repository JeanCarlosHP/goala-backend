-- name: GetFoodByBarcode :one
SELECT * FROM food_database WHERE barcode = $1;

-- name: CreateFoodFromBarcode :one
INSERT INTO food_database (
    barcode, name, brand, calories, protein, carbs, fat,
    serving_size, serving_unit, source
) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
RETURNING *;

-- name: SearchFoodByName :many
SELECT * FROM food_database
WHERE LOWER(name) LIKE LOWER($1)
ORDER BY verified DESC, name
LIMIT $2;

-- name: GetFoodByID :one
SELECT * FROM food_database WHERE id = $1;

-- name: UpdateFoodVerified :exec
UPDATE food_database SET verified = $2, updated_at = NOW() WHERE id = $1;

-- name: ListVerifiedFoods :many
SELECT * FROM food_database
WHERE verified = true
ORDER BY name
LIMIT $1 OFFSET $2;
