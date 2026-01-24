package domain

import (
	"time"

	"github.com/google/uuid"
)

type UserStats struct {
	ID                    uuid.UUID  `json:"id" db:"id"`
	UserID                uuid.UUID  `json:"user_id" db:"user_id"`
	CurrentStreak         int32      `json:"currentStreak" db:"current_streak"`
	LongestStreak         int32      `json:"bestStreak" db:"longest_streak"`
	TotalMealsLogged      int32      `json:"totalMealsLogged" db:"total_meals_logged"`
	TotalDaysLogged       int32      `json:"totalDaysLogged" db:"total_days_logged"`
	TotalCaloriesConsumed int32      `json:"totalCaloriesLogged" db:"total_calories_consumed"`
	TotalProteinConsumed  int32      `json:"-" db:"total_protein_consumed"`
	TotalCarbsConsumed    int32      `json:"-" db:"total_carbs_consumed"`
	TotalFatConsumed      int32      `json:"-" db:"total_fat_consumed"`
	LastLogDate           *time.Time `json:"-" db:"last_log_date"`
	CreatedAt             time.Time  `json:"-" db:"created_at"`
	UpdatedAt             time.Time  `json:"-" db:"updated_at"`
}

type UserStatsResponse struct {
	CurrentStreak         int32 `json:"currentStreak" validate:"gte=0"`
	BestStreak            int32 `json:"bestStreak" validate:"gte=0"`
	TotalMealsLogged      int32 `json:"totalMealsLogged" validate:"gte=0"`
	TotalCaloriesLogged   int32 `json:"totalCaloriesLogged" validate:"gte=0"`
	TotalDaysLogged       int32 `json:"totalDaysLogged" validate:"gte=0"`
	AverageCaloriesPerDay int32 `json:"averageCaloriesPerDay" validate:"gte=0"`
}

type DayStats struct {
	Date          time.Time `json:"date" validate:"required"`
	TotalCalories int32     `json:"totalCalories" validate:"gte=0"`
	TotalProtein  int32     `json:"totalProtein" validate:"gte=0"`
	TotalCarbs    int32     `json:"totalCarbs" validate:"gte=0"`
	TotalFat      int32     `json:"totalFat" validate:"gte=0"`
	Meals         []Meal    `json:"meals" validate:"dive"`
	WaterIntake   int32     `json:"waterIntake" validate:"gte=0"`
}

type Pagination struct {
	Page       int  `json:"page" validate:"required,gte=1"`
	Limit      int  `json:"limit" validate:"required,gte=1,lte=100"`
	Total      int  `json:"total" validate:"gte=0"`
	TotalPages int  `json:"totalPages" validate:"gte=0"`
	HasNext    bool `json:"hasNext"`
	HasPrev    bool `json:"hasPrev"`
}

type AggregatedStats struct {
	TotalCalories int32 `json:"totalCalories" validate:"gte=0"`
	TotalProtein  int32 `json:"totalProtein" validate:"gte=0"`
	TotalCarbs    int32 `json:"totalCarbs" validate:"gte=0"`
	TotalFat      int32 `json:"totalFat" validate:"gte=0"`
	AvgCalories   int32 `json:"avgCalories" validate:"gte=0"`
	AvgProtein    int32 `json:"avgProtein" validate:"gte=0"`
	AvgCarbs      int32 `json:"avgCarbs" validate:"gte=0"`
	AvgFat        int32 `json:"avgFat" validate:"gte=0"`
}

type StatsRangeQuery struct {
	StartDate string `query:"startDate" validate:"required"`
	EndDate   string `query:"endDate" validate:"required"`
	Page      int    `query:"page" validate:"omitempty,gte=1"`
	Limit     int    `query:"limit" validate:"omitempty,gte=1,lte=100"`
}

type StatsRangeResponse struct {
	Days       []DayStats      `json:"days" validate:"dive"`
	Pagination Pagination      `json:"pagination"`
	Aggregated AggregatedStats `json:"aggregated"`
}
