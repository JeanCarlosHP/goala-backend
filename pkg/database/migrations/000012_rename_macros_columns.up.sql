ALTER TABLE user_goals RENAME COLUMN protein_g TO protein;
ALTER TABLE user_goals RENAME COLUMN carbs_g TO carbs;
ALTER TABLE user_goals RENAME COLUMN fat_g TO fat;

ALTER TABLE food_items RENAME COLUMN protein_g TO protein;
ALTER TABLE food_items RENAME COLUMN carbs_g TO carbs;
ALTER TABLE food_items RENAME COLUMN fat_g TO fat;

ALTER TABLE food_database RENAME COLUMN protein_g TO protein;
ALTER TABLE food_database RENAME COLUMN carbs_g TO carbs;
ALTER TABLE food_database RENAME COLUMN fat_g TO fat;
