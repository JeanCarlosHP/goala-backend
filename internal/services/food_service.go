package services

import (
	"context"

	"github.com/jeancarloshp/calorieai/internal/repositories"

	"github.com/google/uuid"
	"github.com/jeancarloshp/calorieai/internal/domain"
)

type FoodService struct {
	foodRepo *repositories.FoodRepository
}

func NewFoodService(foodRepo *repositories.FoodRepository) *FoodService {
	return &FoodService{
		foodRepo: foodRepo,
	}
}

func (s *FoodService) SearchFoods(ctx context.Context, query string) ([]domain.FoodDatabase, error) {
	if query == "" {
		return []domain.FoodDatabase{}, nil
	}

	return s.foodRepo.SearchFoodDatabase(ctx, query, 20)
}

func (s *FoodService) GetRecentFoods(ctx context.Context, userID uuid.UUID) ([]domain.RecentFood, error) {
	return s.foodRepo.GetRecentFoods(ctx, userID, 20)
}

func (s *FoodService) CreateFoodItem(ctx context.Context, req *domain.CreateFoodItemRequest) (*domain.FoodItem, error) {
	return s.foodRepo.CreateStandalone(ctx, req)
}

func (s *FoodService) GetFoodItem(ctx context.Context, id uuid.UUID) (*domain.FoodItem, error) {
	return s.foodRepo.GetByID(ctx, id)
}

func (s *FoodService) UpdateFoodItem(ctx context.Context, id uuid.UUID, req *domain.UpdateFoodItemRequest) (*domain.FoodItem, error) {
	_, err := s.foodRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return s.foodRepo.Update(ctx, id, req)
}

func (s *FoodService) DeleteFoodItem(ctx context.Context, id uuid.UUID) error {
	_, err := s.foodRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	return s.foodRepo.Delete(ctx, id)
}
