package services

import (
	"context"
	"encoding/base64"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"strings"
	"time"

	"github.com/jeancarloshp/calorieai/internal/domain"
	"go.opentelemetry.io/otel"
)

type FoodRecognitionService struct {
	s3Service  *S3Service
	aiProvider AIProvider
	config     *domain.Config
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
		config:     cfg,
		logger:     logger,
	}
}

func (s *FoodRecognitionService) RecognizeFoodWithProgress(
	ctx context.Context,
	fileHeader *multipart.FileHeader,
	req *domain.FoodRecognitionRequest,
	progressChan chan<- domain.ProgressUpdate,
) (*domain.FoodRecognitionResponse, error) {
	tr := otel.Tracer("services/food_recognition_service.go")
	ctx, span := tr.Start(ctx, "RecognizeFoodWithProgress")
	defer span.End()

	startTime := time.Now()
	defer close(progressChan)

	progressChan <- domain.ProgressUpdate{
		Stage:      "upload",
		Percentage: 0,
		Message:    "Starting image upload...",
	}

	file, err := fileHeader.Open()
	if err != nil {
		s.logger.Error("failed to open uploaded file", "error", err)
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	progressChan <- domain.ProgressUpdate{
		Stage:      "upload",
		Percentage: 10,
		Message:    "Reading image data...",
	}

	// Ler o conteúdo do arquivo
	imageBytes, err := io.ReadAll(file)
	if err != nil {
		s.logger.Error("failed to read file", "error", err)
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	progressChan <- domain.ProgressUpdate{
		Stage:      "upload",
		Percentage: 20,
		Message:    "Uploading to S3...",
	}

	// Upload para S3
	file.Seek(0, 0) // Reset file pointer
	imageURL, err := s.s3Service.UploadImage(ctx, file, fileHeader.Header.Get("Content-Type"))
	if err != nil {
		s.logger.Error("failed to upload image to S3", "error", err)
		return nil, fmt.Errorf("failed to upload image: %w", err)
	}

	s.logger.Info("image uploaded successfully", "url", imageURL)

	progressChan <- domain.ProgressUpdate{
		Stage:      "upload",
		Percentage: 25,
		Message:    "Image uploaded successfully",
	}

	// Converter para base64 para enviar para a IA
	imageBase64 := base64.StdEncoding.EncodeToString(imageBytes)

	// Chamar IA para reconhecimento
	foodItems, err := s.aiProvider.RecognizeFood(ctx, imageBase64, progressChan)
	if err != nil {
		s.logger.Error("failed to recognize food", "error", err)
		return nil, fmt.Errorf("failed to recognize food: %w", err)
	}

	progressChan <- domain.ProgressUpdate{
		Stage:      "complete",
		Percentage: 100,
		Message:    fmt.Sprintf("Successfully recognized %d food items", len(foodItems)),
	}

	processingTime := int32(time.Since(startTime).Milliseconds())

	return &domain.FoodRecognitionResponse{
		FoodItems:      foodItems,
		ProcessingTime: processingTime,
	}, nil
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
	foodItems, err := s.aiProvider.RecognizeFood(ctx, imageBase64, nil)
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

func (s *FoodRecognitionService) EstimateQuantityWithProgress(
	ctx context.Context,
	fileHeader *multipart.FileHeader,
	req *domain.EstimateQuantityRequest,
	progressChan chan<- domain.ProgressUpdate,
) (*domain.EstimateQuantityResponse, error) {
	tr := otel.Tracer("services/food_recognition_service.go")
	ctx, span := tr.Start(ctx, "EstimateQuantityWithProgress")
	defer span.End()

	defer close(progressChan)

	progressChan <- domain.ProgressUpdate{
		Stage:      "upload",
		Percentage: 0,
		Message:    "Starting image upload...",
	}

	file, err := fileHeader.Open()
	if err != nil {
		s.logger.Error("failed to open uploaded file", "error", err)
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	progressChan <- domain.ProgressUpdate{
		Stage:      "upload",
		Percentage: 10,
		Message:    "Reading image data...",
	}

	// Ler o conteúdo do arquivo
	imageBytes, err := io.ReadAll(file)
	if err != nil {
		s.logger.Error("failed to read file", "error", err)
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	progressChan <- domain.ProgressUpdate{
		Stage:      "upload",
		Percentage: 20,
		Message:    "Uploading to S3...",
	}

	// Upload para S3
	file.Seek(0, 0) // Reset file pointer
	imageURL, err := s.s3Service.UploadImage(ctx, file, fileHeader.Header.Get("Content-Type"))
	if err != nil {
		s.logger.Error("failed to upload image to S3", "error", err)
		return nil, fmt.Errorf("failed to upload image: %w", err)
	}

	s.logger.Info("image uploaded successfully", "url", imageURL)

	progressChan <- domain.ProgressUpdate{
		Stage:      "upload",
		Percentage: 25,
		Message:    "Image uploaded successfully",
	}

	// Converter para base64
	imageBase64 := base64.StdEncoding.EncodeToString(imageBytes)

	// Chamar IA para estimativa
	result, err := s.aiProvider.EstimateQuantity(ctx, imageBase64, req, progressChan)
	if err != nil {
		s.logger.Error("failed to estimate quantity", "error", err)
		return nil, fmt.Errorf("failed to estimate quantity: %w", err)
	}

	progressChan <- domain.ProgressUpdate{
		Stage:      "complete",
		Percentage: 100,
		Message:    "Quantity estimation completed",
	}

	return result, nil
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
	result, err := s.aiProvider.EstimateQuantity(ctx, imageBase64, req, nil)
	if err != nil {
		s.logger.Error("failed to estimate quantity", "error", err)
		return nil, fmt.Errorf("failed to estimate quantity: %w", err)
	}

	return result, nil
}

func (s *FoodRecognitionService) downloadImageFromCDN(ctx context.Context, imagePath string) ([]byte, error) {
	// Remove leading slash if present
	imagePath = strings.TrimPrefix(imagePath, "/")

	// Build CDN URL
	cdnURL := fmt.Sprintf("https://%s/%s", s.config.CDNDomain, imagePath)

	// Create HTTP request
	req, err := http.NewRequestWithContext(ctx, "GET", cdnURL, nil)
	if err != nil {
		s.logger.Error("failed to create HTTP request for CDN", "error", err, "url", cdnURL)
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Make the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		s.logger.Error("failed to download image from CDN", "error", err, "url", cdnURL)
		return nil, fmt.Errorf("failed to download image: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		s.logger.Error("CDN returned non-200 status", "status", resp.StatusCode, "url", cdnURL)
		return nil, fmt.Errorf("failed to download image: status %d", resp.StatusCode)
	}

	// Read the image data
	imageBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		s.logger.Error("failed to read image data from CDN", "error", err, "url", cdnURL)
		return nil, fmt.Errorf("failed to read image data: %w", err)
	}

	return imageBytes, nil
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

	// Download image from CDN
	imageBytes, err := s.downloadImageFromCDN(ctx, imagePath)
	if err != nil {
		s.logger.Error("failed to download image from CDN", "error", err, "imagePath", imagePath)
		return nil, fmt.Errorf("failed to download image: %w", err)
	}

	// Converter para base64 para enviar para a IA
	imageBase64 := base64.StdEncoding.EncodeToString(imageBytes)

	// Chamar IA para reconhecimento
	foodItems, err := s.aiProvider.RecognizeFood(ctx, imageBase64, nil)
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

	// Download image from CDN
	imageBytes, err := s.downloadImageFromCDN(ctx, imagePath)
	if err != nil {
		s.logger.Error("failed to download image from CDN", "error", err, "imagePath", imagePath)
		return nil, fmt.Errorf("failed to download image: %w", err)
	}

	// Converter para base64
	imageBase64 := base64.StdEncoding.EncodeToString(imageBytes)

	// Chamar IA para estimativa
	result, err := s.aiProvider.EstimateQuantity(ctx, imageBase64, req, nil)
	if err != nil {
		s.logger.Error("failed to estimate quantity", "error", err)
		return nil, fmt.Errorf("failed to estimate quantity: %w", err)
	}

	return result, nil
}
