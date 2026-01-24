CREATE TABLE IF NOT EXISTS achievements (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name_key VARCHAR(100) UNIQUE NOT NULL,
    description_key VARCHAR(255) NOT NULL,
    icon VARCHAR(100) NOT NULL,
    target INT NOT NULL,
    category VARCHAR(50) NOT NULL CHECK (category IN ('streak', 'meals', 'calories', 'protein', 'general')),
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_achievements_category ON achievements(category);

CREATE TABLE IF NOT EXISTS user_achievements (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    achievement_id UUID NOT NULL REFERENCES achievements(id) ON DELETE CASCADE,
    unlocked BOOLEAN DEFAULT false NOT NULL,
    progress INT DEFAULT 0 NOT NULL,
    unlocked_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    UNIQUE(user_id, achievement_id)
);

CREATE INDEX IF NOT EXISTS idx_user_achievements_user_id ON user_achievements(user_id);
CREATE INDEX IF NOT EXISTS idx_user_achievements_unlocked ON user_achievements(unlocked);
CREATE INDEX IF NOT EXISTS idx_user_achievements_achievement_id ON user_achievements(achievement_id);

INSERT INTO achievements (name_key, description_key, icon, target, category) VALUES
('first_meal', 'achievement_first_meal_desc', '🍽️', 1, 'meals'),
('10_meals', 'achievement_10_meals_desc', '📊', 10, 'meals'),
('50_meals', 'achievement_50_meals_desc', '🏆', 50, 'meals'),
('100_meals', 'achievement_100_meals_desc', '🎯', 100, 'meals'),
('first_streak', 'achievement_first_streak_desc', '🔥', 1, 'streak'),
('7_day_streak', 'achievement_7_day_streak_desc', '⚡', 7, 'streak'),
('30_day_streak', 'achievement_30_day_streak_desc', '💪', 30, 'streak'),
('100_day_streak', 'achievement_100_day_streak_desc', '👑', 100, 'streak'),
('protein_goal_met', 'achievement_protein_goal_met_desc', '💪', 1, 'protein'),
('calorie_goal_met', 'achievement_calorie_goal_met_desc', '🎯', 1, 'calories');
