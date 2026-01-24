CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    firebase_uid VARCHAR(128) UNIQUE NOT NULL,
    email VARCHAR(255),
    display_name VARCHAR(255),
    photo_url TEXT,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE user_goals (
    user_id UUID PRIMARY KEY REFERENCES users(id) ON DELETE CASCADE,
    daily_calories INT NOT NULL DEFAULT 2000,
    protein_g INT DEFAULT 150,
    carbs_g INT DEFAULT 200,
    fat_g INT DEFAULT 65,
    updated_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE meals (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    meal_type VARCHAR(50),
    meal_date DATE NOT NULL,
    meal_time TIME,
    photo_url TEXT,
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE food_items (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    meal_id UUID NOT NULL REFERENCES meals(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    portion_size DECIMAL(10,2),
    portion_unit VARCHAR(50),
    calories INT NOT NULL,
    protein_g DECIMAL(10,2),
    carbs_g DECIMAL(10,2),
    fat_g DECIMAL(10,2),
    source VARCHAR(50)
);

CREATE INDEX idx_meals_user_date ON meals(user_id, meal_date);
CREATE INDEX idx_food_items_meal ON food_items(meal_id);
CREATE INDEX idx_users_firebase_uid ON users(firebase_uid);
