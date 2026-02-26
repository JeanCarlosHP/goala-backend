-- name: GetMealsWithFoodsInRange :many
SELECT 
    m.id as meal_id,
    m.user_id,
    m.meal_type,
    m.meal_date,
    m.meal_time,
    m.photo_url,
    m.created_at as meal_created_at,
    fi.id as food_id,
    fi.name as food_name,
    fi.portion_size,
    fi.portion_unit,
    fi.calories,
    fi.protein,
    fi.carbs,
    fi.fat,
    fi.source
FROM meals m
LEFT JOIN food_items fi ON fi.meal_id = m.id
WHERE m.user_id = $1 
  AND m.meal_date >= $2 
  AND m.meal_date <= $3
ORDER BY m.meal_date ASC, m.meal_time ASC, m.created_at ASC;
