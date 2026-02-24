package domain

import (
	"time"

	"github.com/google/uuid"
)

type Meal struct {
	ID        uuid.UUID  `json:"id" db:"id"`
	UserID    uuid.UUID  `json:"user_id" db:"user_id"`
	MealType  string     `json:"meal_type" db:"meal_type"`
	MealDate  time.Time  `json:"meal_date" db:"meal_date"`
	MealTime  *time.Time `json:"meal_time,omitempty" db:"meal_time"`
	PhotoURL  *string    `json:"photo_url,omitempty" db:"photo_url"`
	CreatedAt time.Time  `json:"created_at" db:"created_at"`
	Foods     []FoodItem `json:"foods,omitempty"`
}

type FoodItem struct {
	ID          uuid.UUID `json:"id" db:"id"`
	MealID      uuid.UUID `json:"meal_id" db:"meal_id"`
	Name        string    `json:"name" db:"name"`
	PortionSize float64   `json:"portion_size" db:"portion_size"`
	PortionUnit string    `json:"portion_unit" db:"portion_unit"`
	Calories    int       `json:"calories" db:"calories"`
	Protein     float64   `json:"protein" db:"protein"`
	Carbs       float64   `json:"carbs" db:"carbs"`
	Fat         float64   `json:"fat" db:"fat"`
	Source      string    `json:"source" db:"source"`
}

type CreateMealRequest struct {
	MealType string              `json:"meal_type" validate:"required,oneof=breakfast lunch dinner snack"`
	MealDate string              `json:"meal_date" validate:"required"`
	MealTime *string             `json:"meal_time,omitempty"`
	PhotoURL *string             `json:"photo_url,omitempty"`
	Foods    []CreateFoodRequest `json:"foods" validate:"required,min=1"`
}

type CreateFoodRequest struct {
	Name        string  `json:"name" validate:"required"`
	PortionSize float64 `json:"portion_size" validate:"required,min=0"`
	PortionUnit string  `json:"portion_unit" validate:"required"`
	Calories    int     `json:"calories" validate:"required,min=0"`
	Protein     float64 `json:"protein" validate:"min=0"`
	Carbs       float64 `json:"carbs" validate:"min=0"`
	Fat         float64 `json:"fat" validate:"min=0"`
	Source      string  `json:"source" validate:"required,oneof=ai_photo ai_text manual"`
}

type FoodDatabase struct {
	ID              uuid.UUID `json:"id" db:"id"`
	Name            string    `json:"name" db:"name"`
	Brand           *string   `json:"brand,omitempty" db:"brand"`
	CaloriesPer100g int       `json:"calories_per_100g" db:"calories_per_100g"`
	ProteinPer100g  float64   `json:"protein_per_100g" db:"protein_per_100g"`
	CarbsPer100g    float64   `json:"carbs_per_100g" db:"carbs_per_100g"`
	FatPer100g      float64   `json:"fat_per_100g" db:"fat_per_100g"`
	Source          string    `json:"source" db:"source"`
	CreatedAt       time.Time `json:"created_at" db:"created_at"`
}

type RecentFood struct {
	Name        string    `json:"name"`
	PortionSize float64   `json:"portion_size"`
	PortionUnit string    `json:"portion_unit"`
	Calories    int       `json:"calories"`
	Protein     float64   `json:"protein"`
	Carbs       float64   `json:"carbs"`
	Fat         float64   `json:"fat"`
	LastUsed    time.Time `json:"last_used"`
}

type DailySummary struct {
	Date          string  `json:"date"`
	TotalCalories int     `json:"total_calories"`
	TotalProtein  float64 `json:"total_protein"`
	TotalCarbs    float64 `json:"total_carbs"`
	TotalFat      float64 `json:"total_fat"`
	Meals         []Meal  `json:"meals"`
}
