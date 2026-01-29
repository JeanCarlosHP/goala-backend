package services

import (
	"context"
	"fmt"

	"github.com/jeancarloshp/calorieai/internal/domain"
	"go.opentelemetry.io/otel"
)

type FoodRecognitionService struct {
	s3Service  *S3Service
	aiProvider AIProvider
	logger     domain.Logger
}

func NewFoodRecognitionService(
	s3Service *S3Service,
	cfg *domain.Config,
	logger domain.Logger,
) *FoodRecognitionService {
	var aiProvider AIProvider

	// Escolhe o provider baseado na config
	switch cfg.AIProvider {
	case "gemini", "":
		aiProvider = NewGeminiProvider(cfg.GeminiAPIKey, cfg.GeminiModel, logger)
	case "openai":
		aiProvider = NewOpenAIProvider(cfg.OpenAIAPIKey, cfg.OpenAIModel, logger)
	default:
		logger.Warn("Unknown AI provider, using Gemini", "provider", cfg.AIProvider)
		aiProvider = NewGeminiProvider(cfg.GeminiAPIKey, cfg.GeminiModel, logger)
	}

	return &FoodRecognitionService{
		s3Service:  s3Service,
		aiProvider: aiProvider,
		logger:     logger,
	}
}

func (s *FoodRecognitionService) RecognizeFoodByPath(
	ctx context.Context,
	imagePath string,
	req *domain.FoodRecognitionRequest,
) (*domain.FoodRecognitionResponse, error) {
	tr := otel.Tracer("services/food_recognition_service.go")
	ctx, span := tr.Start(ctx, "RecognizeFoodByPath")
	defer span.End()

	// Download image from S3
	imageBytes, err := s.s3Service.DownloadImage(ctx, imagePath)
	if err != nil {
		s.logger.Error("failed to download image from S3", "error", err, "imagePath", imagePath)
		return nil, fmt.Errorf("failed to download image: %w", err)
	}

	// Chamar IA para reconhecimento
	foodItems, err := s.aiProvider.RecognizeFood(ctx, imageBytes)
	if err != nil {
		s.logger.Error("failed to recognize food", "error", err)
		return nil, fmt.Errorf("failed to recognize food: %w", err)
	}

	return &domain.FoodRecognitionResponse{
		FoodItems: foodItems,
	}, nil
}

func (s *FoodRecognitionService) EstimateQuantityByPath(
	ctx context.Context,
	imagePath string,
	req *domain.EstimateQuantityRequest,
) (*domain.EstimateQuantityResponse, error) {
	tr := otel.Tracer("services/food_recognition_service.go")
	ctx, span := tr.Start(ctx, "EstimateQuantityByPath")
	defer span.End()

	// Download image from S3
	imageBytes, err := s.s3Service.DownloadImage(ctx, imagePath)
	if err != nil {
		s.logger.Error("failed to download image from S3", "error", err, "imagePath", imagePath)
		return nil, fmt.Errorf("failed to download image: %w", err)
	}

	// Chamar IA para estimativa
	result, err := s.aiProvider.EstimateQuantity(ctx, imageBytes, req)
	if err != nil {
		s.logger.Error("failed to estimate quantity", "error", err)
		return nil, fmt.Errorf("failed to estimate quantity: %w", err)
	}

	return result, nil
}
