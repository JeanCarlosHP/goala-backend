package services

import (
	"context"

	"github.com/jeancarloshp/calorieai/internal/domain"
)

type AIProvider interface {
	RecognizeFood(
		ctx context.Context,
		imageBase64 string,
		progressChan chan<- domain.ProgressUpdate,
	) ([]domain.RecognizedFoodItem, error)

	EstimateQuantity(
		ctx context.Context,
		imageBase64 string,
		req *domain.EstimateQuantityRequest,
		progressChan chan<- domain.ProgressUpdate,
	) (*domain.EstimateQuantityResponse, error)
}
