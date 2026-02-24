ALTER TABLE food_database RENAME COLUMN fat TO fat_g;
ALTER TABLE food_database RENAME COLUMN carbs TO carbs_g;
ALTER TABLE food_database RENAME COLUMN protein TO protein_g;

ALTER TABLE food_items RENAME COLUMN fat TO fat_g;
ALTER TABLE food_items RENAME COLUMN carbs TO carbs_g;
ALTER TABLE food_items RENAME COLUMN protein TO protein_g;

ALTER TABLE user_goals RENAME COLUMN fat TO fat_g;
ALTER TABLE user_goals RENAME COLUMN carbs TO carbs_g;
ALTER TABLE user_goals RENAME COLUMN protein TO protein_g;
