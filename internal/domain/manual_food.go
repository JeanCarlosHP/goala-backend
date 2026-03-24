package domain

import "github.com/google/uuid"

type FoodPortion struct {
	Name  string  `json:"name"`
	Grams float64 `json:"grams"`
}

type SearchFood struct {
	ID         *uuid.UUID    `json:"id,omitempty"`
	ExternalID *string       `json:"external_id,omitempty"`
	Name       string        `json:"name" validate:"required"`
	Brand      *string       `json:"brand,omitempty"`
	Calories   float64       `json:"calories" validate:"gte=0"`
	Protein    float64       `json:"protein" validate:"gte=0"`
	Carbs      float64       `json:"carbs" validate:"gte=0"`
	Fat        float64       `json:"fat" validate:"gte=0"`
	Source     string        `json:"source" validate:"required,oneof=internal openfoodfacts"`
	Portions   []FoodPortion `json:"portions,omitempty"`
	IsFavorite bool          `json:"is_favorite"`
}

type FoodSearchRequest struct {
	Query string `json:"query" validate:"max=120"`
	Limit int    `json:"limit" validate:"omitempty,min=1,max=50"`
}

type FoodSearchResponse struct {
	Query     string       `json:"query"`
	Foods     []SearchFood `json:"foods"`
	Recent    []SearchFood `json:"recent"`
	Favorites []SearchFood `json:"favorites"`
}

type LogMealFoodRequest struct {
	Date         string     `json:"date" validate:"required"`
	MealType     string     `json:"meal_type" validate:"required,oneof=breakfast lunch dinner snack"`
	Food         SearchFood `json:"food" validate:"required"`
	Quantity     float64    `json:"quantity" validate:"required,gt=0"`
	PortionName  string     `json:"portion_name" validate:"required"`
	PortionGrams *float64   `json:"portion_grams,omitempty" validate:"omitempty,gt=0"`
}
