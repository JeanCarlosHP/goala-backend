package services

import (
	"context"

	"github.com/jeancarloshp/calorieai/internal/domain"
)

type AIProvider interface {
	RecognizeFood(
		ctx context.Context,
		imageBase64 string,
	) ([]domain.RecognizedFoodItem, error)

	EstimateQuantity(
		ctx context.Context,
		imageBase64 string,
		req *domain.EstimateQuantityRequest,
	) (*domain.EstimateQuantityResponse, error)
}
