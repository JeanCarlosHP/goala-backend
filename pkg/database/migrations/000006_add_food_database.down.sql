DROP INDEX IF EXISTS idx_food_database_verified;
DROP INDEX IF EXISTS idx_food_database_barcode;
DROP INDEX IF EXISTS idx_food_database_barcode_unique;

ALTER TABLE food_database DROP COLUMN IF EXISTS updated_at;
ALTER TABLE food_database DROP COLUMN IF EXISTS verified;
ALTER TABLE food_database DROP COLUMN IF EXISTS serving_unit;
ALTER TABLE food_database DROP COLUMN IF EXISTS serving_size;
ALTER TABLE food_database DROP COLUMN IF EXISTS fat_g;
ALTER TABLE food_database DROP COLUMN IF EXISTS carbs_g;
ALTER TABLE food_database DROP COLUMN IF EXISTS protein_g;
ALTER TABLE food_database DROP COLUMN IF EXISTS calories;
ALTER TABLE food_database DROP COLUMN IF EXISTS barcode;
