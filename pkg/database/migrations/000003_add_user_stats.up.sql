CREATE TABLE IF NOT EXISTS user_stats (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID UNIQUE NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    current_streak INT DEFAULT 0 NOT NULL,
    longest_streak INT DEFAULT 0 NOT NULL,
    total_meals_logged INT DEFAULT 0 NOT NULL,
    total_days_logged INT DEFAULT 0 NOT NULL,
    total_calories_consumed INT DEFAULT 0 NOT NULL,
    total_protein_consumed INT DEFAULT 0 NOT NULL,
    total_carbs_consumed INT DEFAULT 0 NOT NULL,
    total_fat_consumed INT DEFAULT 0 NOT NULL,
    last_log_date DATE,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_user_stats_user_id ON user_stats(user_id);
CREATE INDEX IF NOT EXISTS idx_user_stats_current_streak ON user_stats(current_streak DESC);
CREATE INDEX IF NOT EXISTS idx_user_stats_last_log_date ON user_stats(last_log_date DESC);
