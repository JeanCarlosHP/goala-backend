ALTER TABLE food_database
    ADD COLUMN IF NOT EXISTS external_id VARCHAR(100);

ALTER TABLE food_items
    ADD COLUMN IF NOT EXISTS food_database_id UUID REFERENCES food_database(id) ON DELETE SET NULL;

CREATE TABLE IF NOT EXISTS food_portions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    food_id UUID NOT NULL REFERENCES food_database(id) ON DELETE CASCADE,
    name VARCHAR(100) NOT NULL,
    grams DECIMAL(10,2) NOT NULL CHECK (grams > 0),
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    UNIQUE (food_id, name)
);

CREATE TABLE IF NOT EXISTS favorite_foods (
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    food_id UUID NOT NULL REFERENCES food_database(id) ON DELETE CASCADE,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    PRIMARY KEY (user_id, food_id)
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_food_database_external_id_unique
    ON food_database(external_id)
    WHERE external_id IS NOT NULL;

CREATE INDEX IF NOT EXISTS idx_food_items_food_database_id
    ON food_items(food_database_id);

CREATE INDEX IF NOT EXISTS idx_food_portions_food_id
    ON food_portions(food_id);

CREATE INDEX IF NOT EXISTS idx_favorite_foods_user_id
    ON favorite_foods(user_id);
