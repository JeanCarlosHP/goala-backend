package domain

import "github.com/google/uuid"

type CreateFoodItemRequest struct {
	MealID      uuid.UUID `json:"meal_id" validate:"required,uuid"`
	Name        string    `json:"name" validate:"required,min=1,max=200"`
	PortionSize float64   `json:"portion_size" validate:"required,min=0,max=10000"`
	PortionUnit string    `json:"portion_unit" validate:"required,oneof=g ml serving cup tbsp tsp oz lb kg piece slice unit"`
	Calories    int       `json:"calories" validate:"required,min=0,max=5000"`
	Protein     float64   `json:"protein" validate:"required,min=0,max=500"`
	Carbs       float64   `json:"carbs" validate:"required,min=0,max=500"`
	Fat         float64   `json:"fat" validate:"required,min=0,max=500"`
	Source      string    `json:"source" validate:"required,oneof=ai_photo ai_text manual barcode"`
}

type UpdateFoodItemRequest struct {
	Name        string  `json:"name" validate:"required,min=1,max=200"`
	PortionSize float64 `json:"portion_size" validate:"required,min=0,max=10000"`
	PortionUnit string  `json:"portion_unit" validate:"required,oneof=g ml serving cup tbsp tsp oz lb kg piece slice unit"`
	Calories    int     `json:"calories" validate:"required,min=0,max=5000"`
	Protein     float64 `json:"protein" validate:"required,min=0,max=500"`
	Carbs       float64 `json:"carbs" validate:"required,min=0,max=500"`
	Fat         float64 `json:"fat" validate:"required,min=0,max=500"`
	Source      string  `json:"source" validate:"required,oneof=ai_photo ai_text manual barcode"`
}

type FoodItemResponse struct {
	ID          uuid.UUID `json:"id"`
	MealID      uuid.UUID `json:"meal_id"`
	Name        string    `json:"name"`
	PortionSize float64   `json:"portion_size"`
	PortionUnit string    `json:"portion_unit"`
	Calories    int       `json:"calories"`
	Protein     float64   `json:"protein"`
	Carbs       float64   `json:"carbs"`
	Fat         float64   `json:"fat"`
	Source      string    `json:"source"`
}
