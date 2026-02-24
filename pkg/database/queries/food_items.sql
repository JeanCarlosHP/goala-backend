-- name: CreateFoodItem :exec
INSERT INTO food_items (id, meal_id, name, portion_size, portion_unit, calories, protein, carbs, fat, source)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10);

-- name: GetFoodItemsByMealID :many
SELECT id, meal_id, name, portion_size, portion_unit, calories, protein, carbs, fat, source
FROM food_items
WHERE meal_id = $1;

-- name: GetFoodItemsByMealIDs :many
SELECT id, meal_id, name, portion_size, portion_unit, calories, protein, carbs, fat, source
FROM food_items
WHERE meal_id = ANY($1::uuid[]);

-- name: SearchFoodDatabase :many
SELECT id, name, brand, calories_per_100g, protein_per_100g, carbs_per_100g, fat_per_100g, source, created_at
FROM food_database
WHERE to_tsvector('portuguese', name) @@ plainto_tsquery('portuguese', $1)
   OR name ILIKE '%' || $2 || '%'
ORDER BY ts_rank(to_tsvector('portuguese', name), plainto_tsquery('portuguese', $1)) DESC
LIMIT $3;

-- name: GetRecentFoods :many
SELECT DISTINCT ON (fi.name) 
    fi.name,
    fi.portion_size,
    fi.portion_unit,
    fi.calories,
    fi.protein,
    fi.carbs,
    fi.fat,
    m.created_at as last_used
FROM food_items fi
JOIN meals m ON fi.meal_id = m.id
WHERE m.user_id = $1
ORDER BY fi.name, m.created_at DESC
LIMIT $2;

-- name: UpdateFoodItem :exec
UPDATE food_items
SET name = $2,
    portion_size = $3,
    portion_unit = $4,
    calories = $5,
    protein = $6,
    carbs = $7,
    fat = $8
WHERE id = $1;

-- name: DeleteFoodItem :exec
DELETE FROM food_items
WHERE id = $1;

-- name: GetFoodItemByID :one
SELECT id, meal_id, name, portion_size, portion_unit, calories, protein, carbs, fat, source
FROM food_items
WHERE id = $1;

-- name: CreateStandaloneFoodItem :one
INSERT INTO food_items (id, meal_id, name, portion_size, portion_unit, calories, protein, carbs, fat, source)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
RETURNING id, meal_id, name, portion_size, portion_unit, calories, protein, carbs, fat, source;

-- name: UpdateFoodItemComplete :one
UPDATE food_items
SET name = $2,
    portion_size = $3,
    portion_unit = $4,
    calories = $5,
    protein = $6,
    carbs = $7,
    fat = $8,
    source = $9
WHERE id = $1
RETURNING id, meal_id, name, portion_size, portion_unit, calories, protein, carbs, fat, source;
