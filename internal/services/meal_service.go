package services

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jeancarloshp/calorieai/internal/domain"
	"github.com/jeancarloshp/calorieai/internal/repositories"
)

type MealService struct {
	mealRepo *repositories.MealRepository
	foodRepo *repositories.FoodRepository
}

func NewMealService(mealRepo *repositories.MealRepository, foodRepo *repositories.FoodRepository) *MealService {
	return &MealService{
		mealRepo: mealRepo,
		foodRepo: foodRepo,
	}
}

func (s *MealService) CreateMeal(ctx context.Context, userID uuid.UUID, req domain.CreateMealRequest) (*domain.Meal, error) {
	mealDate, err := time.Parse("2006-01-02", req.MealDate)
	if err != nil {
		return nil, fmt.Errorf("invalid meal date format: %w", err)
	}

	var mealTime *time.Time
	if req.MealTime != nil {
		t, err := time.Parse("15:04", *req.MealTime)
		if err != nil {
			return nil, fmt.Errorf("invalid meal time format: %w", err)
		}
		mealTime = &t
	}

	meal := &domain.Meal{
		ID:       uuid.New(),
		UserID:   userID,
		MealType: req.MealType,
		MealDate: mealDate,
		MealTime: mealTime,
		PhotoURL: req.PhotoURL,
	}

	if err := s.mealRepo.Create(ctx, meal); err != nil {
		return nil, fmt.Errorf("failed to create meal: %w", err)
	}

	for _, foodReq := range req.Foods {
		food := &domain.FoodItem{
			ID:          uuid.New(),
			MealID:      meal.ID,
			Name:        foodReq.Name,
			PortionSize: foodReq.PortionSize,
			PortionUnit: foodReq.PortionUnit,
			Calories:    foodReq.Calories,
			ProteinG:    foodReq.ProteinG,
			CarbsG:      foodReq.CarbsG,
			FatG:        foodReq.FatG,
			Source:      foodReq.Source,
		}

		if err := s.foodRepo.Create(ctx, food); err != nil {
			return nil, fmt.Errorf("failed to create food item: %w", err)
		}

		meal.Foods = append(meal.Foods, *food)
	}

	return meal, nil
}

func (s *MealService) GetMealsByDate(ctx context.Context, userID uuid.UUID, date time.Time) ([]domain.Meal, error) {
	meals, err := s.mealRepo.GetByUserAndDate(ctx, userID, date)
	if err != nil {
		return nil, fmt.Errorf("failed to get meals: %w", err)
	}

	if len(meals) == 0 {
		return []domain.Meal{}, nil
	}

	mealIDs := make([]uuid.UUID, len(meals))
	for i, meal := range meals {
		mealIDs[i] = meal.ID
	}

	foodsByMeal, err := s.foodRepo.GetByMealIDs(ctx, mealIDs)
	if err != nil {
		return nil, fmt.Errorf("failed to get food items: %w", err)
	}

	for i := range meals {
		if foods, ok := foodsByMeal[meals[i].ID]; ok {
			meals[i].Foods = foods
		}
	}

	return meals, nil
}

func (s *MealService) GetDailySummary(ctx context.Context, userID uuid.UUID, date time.Time, goal *domain.UserGoal) (*domain.DailySummary, error) {
	meals, err := s.GetMealsByDate(ctx, userID, date)
	if err != nil {
		return nil, err
	}

	summary := &domain.DailySummary{
		Date:  date.Format("2006-01-02"),
		Meals: meals,
	}

	if goal != nil {
		summary.GoalCalories = goal.DailyCalorieGoal
		summary.GoalProtein = goal.DailyProteinGoal
		summary.GoalCarbs = goal.DailyCarbsGoal
		summary.GoalFat = goal.DailyFatGoal
	}

	for _, meal := range meals {
		for _, food := range meal.Foods {
			summary.TotalCalories += food.Calories
			summary.TotalProtein += food.ProteinG
			summary.TotalCarbs += food.CarbsG
			summary.TotalFat += food.FatG
		}
	}

	return summary, nil
}
