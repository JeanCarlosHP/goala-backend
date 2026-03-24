DROP INDEX IF EXISTS idx_favorite_foods_user_id;
DROP INDEX IF EXISTS idx_food_portions_food_id;
DROP INDEX IF EXISTS idx_food_items_food_database_id;
DROP INDEX IF EXISTS idx_food_database_external_id_unique;

DROP TABLE IF EXISTS favorite_foods;
DROP TABLE IF EXISTS food_portions;

ALTER TABLE food_items
    DROP COLUMN IF EXISTS food_database_id;

ALTER TABLE food_database
    DROP COLUMN IF EXISTS external_id;
