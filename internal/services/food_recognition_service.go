package services

import (
	"context"
	"encoding/base64"
	"fmt"
	"io"
	"mime/multipart"
	"time"

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
		aiProvider = NewGeminiProvider(cfg.GeminiAPIKey, logger)
	case "openai":
		// TODO: Implementar OpenAI provider
		logger.Warn("OpenAI provider not implemented yet, falling back to Gemini")
		aiProvider = NewGeminiProvider(cfg.GeminiAPIKey, logger)
	default:
		logger.Warn("Unknown AI provider, using Gemini", "provider", cfg.AIProvider)
		aiProvider = NewGeminiProvider(cfg.GeminiAPIKey, logger)
	}

	return &FoodRecognitionService{
		s3Service:  s3Service,
		aiProvider: aiProvider,
		logger:     logger,
	}
}

func (s *FoodRecognitionService) RecognizeFood(
	ctx context.Context,
	fileHeader *multipart.FileHeader,
	req *domain.FoodRecognitionRequest,
) (*domain.FoodRecognitionResponse, error) {
	tr := otel.Tracer("services/food_recognition_service.go")
	ctx, span := tr.Start(ctx, "RecognizeFood")
	defer span.End()

	startTime := time.Now()

	file, err := fileHeader.Open()
	if err != nil {
		s.logger.Error("failed to open uploaded file", "error", err)
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	// Ler o conteúdo do arquivo
	imageBytes, err := io.ReadAll(file)
	if err != nil {
		s.logger.Error("failed to read file", "error", err)
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	// Upload para S3
	file.Seek(0, 0) // Reset file pointer
	imageURL, err := s.s3Service.UploadImage(ctx, file, fileHeader.Header.Get("Content-Type"))
	if err != nil {
		s.logger.Error("failed to upload image to S3", "error", err)
		return nil, fmt.Errorf("failed to upload image: %w", err)
	}

	s.logger.Info("image uploaded successfully", "url", imageURL)

	// Converter para base64 para enviar para a IA
	imageBase64 := base64.StdEncoding.EncodeToString(imageBytes)

	// Chamar IA para reconhecimento (sem progresso)
	foodItems, err := s.aiProvider.RecognizeFood(ctx, imageBase64)
	if err != nil {
		s.logger.Error("failed to recognize food", "error", err)
		return nil, fmt.Errorf("failed to recognize food: %w", err)
	}

	processingTime := int32(time.Since(startTime).Milliseconds())

	return &domain.FoodRecognitionResponse{
		FoodItems:      foodItems,
		ProcessingTime: processingTime,
	}, nil
}

func (s *FoodRecognitionService) EstimateQuantity(
	ctx context.Context,
	fileHeader *multipart.FileHeader,
	req *domain.EstimateQuantityRequest,
) (*domain.EstimateQuantityResponse, error) {
	tr := otel.Tracer("services/food_recognition_service.go")
	ctx, span := tr.Start(ctx, "EstimateQuantity")
	defer span.End()

	file, err := fileHeader.Open()
	if err != nil {
		s.logger.Error("failed to open uploaded file", "error", err)
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	// Ler o conteúdo do arquivo
	imageBytes, err := io.ReadAll(file)
	if err != nil {
		s.logger.Error("failed to read file", "error", err)
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	// Upload para S3
	file.Seek(0, 0) // Reset file pointer
	imageURL, err := s.s3Service.UploadImage(ctx, file, fileHeader.Header.Get("Content-Type"))
	if err != nil {
		s.logger.Error("failed to upload image to S3", "error", err)
		return nil, fmt.Errorf("failed to upload image: %w", err)
	}

	s.logger.Info("image uploaded successfully", "url", imageURL)

	// Converter para base64
	imageBase64 := base64.StdEncoding.EncodeToString(imageBytes)

	// Chamar IA para estimativa (sem progresso)
	result, err := s.aiProvider.EstimateQuantity(ctx, imageBase64, req)
	if err != nil {
		s.logger.Error("failed to estimate quantity", "error", err)
		return nil, fmt.Errorf("failed to estimate quantity: %w", err)
	}

	return result, nil
}

func (s *FoodRecognitionService) RecognizeFoodByPath(
	ctx context.Context,
	imagePath string,
	req *domain.FoodRecognitionRequest,
) (*domain.FoodRecognitionResponse, error) {
	tr := otel.Tracer("services/food_recognition_service.go")
	ctx, span := tr.Start(ctx, "RecognizeFoodByPath")
	defer span.End()

	startTime := time.Now()

	// Download image from S3
	imageBytes, err := s.s3Service.DownloadImage(ctx, imagePath)
	if err != nil {
		s.logger.Error("failed to download image from S3", "error", err, "imagePath", imagePath)
		return nil, fmt.Errorf("failed to download image: %w", err)
	}

	// Converter para base64 para enviar para a IA
	imageBase64 := base64.StdEncoding.EncodeToString(imageBytes)

	// Chamar IA para reconhecimento
	foodItems, err := s.aiProvider.RecognizeFood(ctx, imageBase64)
	if err != nil {
		s.logger.Error("failed to recognize food", "error", err)
		return nil, fmt.Errorf("failed to recognize food: %w", err)
	}

	processingTime := int32(time.Since(startTime).Milliseconds())

	return &domain.FoodRecognitionResponse{
		FoodItems:      foodItems,
		ProcessingTime: processingTime,
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

	// Converter para base64
	imageBase64 := base64.StdEncoding.EncodeToString(imageBytes)

	// Chamar IA para estimativa
	result, err := s.aiProvider.EstimateQuantity(ctx, imageBase64, req)
	if err != nil {
		s.logger.Error("failed to estimate quantity", "error", err)
		return nil, fmt.Errorf("failed to estimate quantity: %w", err)
	}

	return result, nil
}
