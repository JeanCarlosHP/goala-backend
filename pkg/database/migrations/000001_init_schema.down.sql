DROP INDEX IF EXISTS idx_users_firebase_uid;
DROP INDEX IF EXISTS idx_food_items_meal;
DROP INDEX IF EXISTS idx_meals_user_date;

DROP TABLE IF EXISTS food_items;
DROP TABLE IF EXISTS meals;
DROP TABLE IF EXISTS user_goals;
DROP TABLE IF EXISTS users;

DROP EXTENSION IF EXISTS "uuid-ossp";
