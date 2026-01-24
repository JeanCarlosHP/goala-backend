ALTER TABLE food_database ADD COLUMN IF NOT EXISTS barcode VARCHAR(50);
ALTER TABLE food_database ADD COLUMN IF NOT EXISTS calories INT;
ALTER TABLE food_database ADD COLUMN IF NOT EXISTS protein_g DECIMAL(10,2);
ALTER TABLE food_database ADD COLUMN IF NOT EXISTS carbs_g DECIMAL(10,2);
ALTER TABLE food_database ADD COLUMN IF NOT EXISTS fat_g DECIMAL(10,2);
ALTER TABLE food_database ADD COLUMN IF NOT EXISTS serving_size INT;
ALTER TABLE food_database ADD COLUMN IF NOT EXISTS serving_unit VARCHAR(50);
ALTER TABLE food_database ADD COLUMN IF NOT EXISTS verified BOOLEAN DEFAULT false;
ALTER TABLE food_database ADD COLUMN IF NOT EXISTS updated_at TIMESTAMP DEFAULT NOW();

UPDATE food_database SET 
    calories = calories_per_100g,
    protein_g = protein_per_100g,
    carbs_g = carbs_per_100g,
    fat_g = fat_per_100g
WHERE calories IS NULL;

CREATE UNIQUE INDEX IF NOT EXISTS idx_food_database_barcode_unique ON food_database(barcode) WHERE barcode IS NOT NULL;
CREATE INDEX IF NOT EXISTS idx_food_database_barcode ON food_database(barcode);
CREATE INDEX IF NOT EXISTS idx_food_database_name ON food_database(name);
CREATE INDEX IF NOT EXISTS idx_food_database_verified ON food_database(verified);
