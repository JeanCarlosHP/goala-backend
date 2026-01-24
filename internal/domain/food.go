package domain

import "github.com/google/uuid"

type CreateFoodItemRequest struct {
	MealID      uuid.UUID `json:"meal_id" validate:"required,uuid"`
	Name        string    `json:"name" validate:"required,min=1,max=200"`
	PortionSize float64   `json:"portion_size" validate:"required,min=0,max=10000"`
	PortionUnit string    `json:"portion_unit" validate:"required,oneof=g ml serving cup tbsp tsp oz lb kg piece slice unit"`
	Calories    int       `json:"calories" validate:"required,min=0,max=5000"`
	ProteinG    float64   `json:"protein_g" validate:"required,min=0,max=500"`
	CarbsG      float64   `json:"carbs_g" validate:"required,min=0,max=500"`
	FatG        float64   `json:"fat_g" validate:"required,min=0,max=500"`
	Source      string    `json:"source" validate:"required,oneof=ai_photo ai_text manual barcode"`
}

type UpdateFoodItemRequest struct {
	Name        string  `json:"name" validate:"required,min=1,max=200"`
	PortionSize float64 `json:"portion_size" validate:"required,min=0,max=10000"`
	PortionUnit string  `json:"portion_unit" validate:"required,oneof=g ml serving cup tbsp tsp oz lb kg piece slice unit"`
	Calories    int     `json:"calories" validate:"required,min=0,max=5000"`
	ProteinG    float64 `json:"protein_g" validate:"required,min=0,max=500"`
	CarbsG      float64 `json:"carbs_g" validate:"required,min=0,max=500"`
	FatG        float64 `json:"fat_g" validate:"required,min=0,max=500"`
	Source      string  `json:"source" validate:"required,oneof=ai_photo ai_text manual barcode"`
}

type FoodItemResponse struct {
	ID          uuid.UUID `json:"id"`
	MealID      uuid.UUID `json:"meal_id"`
	Name        string    `json:"name"`
	PortionSize float64   `json:"portion_size"`
	PortionUnit string    `json:"portion_unit"`
	Calories    int       `json:"calories"`
	ProteinG    float64   `json:"protein_g"`
	CarbsG      float64   `json:"carbs_g"`
	FatG        float64   `json:"fat_g"`
	Source      string    `json:"source"`
}
