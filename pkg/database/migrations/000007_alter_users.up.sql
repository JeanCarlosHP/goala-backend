ALTER TABLE users ADD COLUMN IF NOT EXISTS weight INT CHECK (weight > 0 AND weight <= 1000);
ALTER TABLE users ADD COLUMN IF NOT EXISTS height INT CHECK (height > 0 AND height <= 300);
ALTER TABLE users ADD COLUMN IF NOT EXISTS age INT CHECK (age > 0 AND age <= 150);
ALTER TABLE users ADD COLUMN IF NOT EXISTS gender VARCHAR(20) CHECK (gender IN ('male', 'female', 'other'));
ALTER TABLE users ADD COLUMN IF NOT EXISTS activity_level VARCHAR(20) CHECK (activity_level IN ('sedentary', 'light', 'moderate', 'active', 'very_active'));
ALTER TABLE users ADD COLUMN IF NOT EXISTS language VARCHAR(10) DEFAULT 'en-US' CHECK (language IN ('en-US', 'pt-BR'));
ALTER TABLE users ADD COLUMN IF NOT EXISTS notifications_enabled BOOLEAN DEFAULT false;
