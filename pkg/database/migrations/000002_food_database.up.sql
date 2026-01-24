CREATE TABLE food_database (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(255) NOT NULL,
    brand VARCHAR(255),
    calories_per_100g INT,
    protein_per_100g DECIMAL(10,2),
    carbs_per_100g DECIMAL(10,2),
    fat_per_100g DECIMAL(10,2),
    source VARCHAR(50),
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_food_name_search ON food_database USING GIN(to_tsvector('portuguese', name));
CREATE INDEX idx_food_brand ON food_database(brand);
CREATE INDEX idx_food_source ON food_database(source);
