ALTER TABLE user_goals ADD COLUMN IF NOT EXISTS carbs_g INT CHECK (carbs_g >= 0 AND carbs_g <= 2000);
ALTER TABLE user_goals ADD COLUMN IF NOT EXISTS fat_g INT CHECK (fat_g >= 0 AND fat_g <= 1000);
