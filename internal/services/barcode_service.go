package services

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/jeancarloshp/calorieai/internal/domain"
	"github.com/jeancarloshp/calorieai/pkg/database/db"
	"go.opentelemetry.io/otel"
)

type BarcodeService struct {
	foodRepo         *db.Queries
	openFoodFactsURL string
	httpClient       *http.Client
	logger           domain.Logger
}

func NewBarcodeService(
	foodRepo *db.Queries,
	cfg *domain.Config,
	logger domain.Logger,
) *BarcodeService {
	openFoodFactsURL := cfg.OpenFoodFactsAPIURL
	if openFoodFactsURL == "" {
		openFoodFactsURL = "https://world.openfoodfacts.org/api/v2"
	}

	return &BarcodeService{
		foodRepo:         foodRepo,
		openFoodFactsURL: openFoodFactsURL,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
		logger: logger,
	}
}

func (s *BarcodeService) GetFoodByBarcode(ctx context.Context, barcode string) (*domain.FoodBarcodeResponse, error) {
	tr := otel.Tracer("services/barcode_service.go")
	ctx, span := tr.Start(ctx, "GetFoodByBarcode")
	defer span.End()

	cached, err := s.foodRepo.GetFoodByBarcode(ctx, &barcode)
	if err == nil {
		s.logger.Info("found food in cache", "barcode", barcode)
		return s.mapDBToResponse(ctx, &cached), nil
	}

	s.logger.Info("food not in cache, fetching from OpenFoodFacts", "barcode", barcode)

	offFood, err := s.fetchFromOpenFoodFacts(ctx, barcode)
	if err != nil {
		return nil, err
	}

	if err := s.cacheFoodInDB(ctx, barcode, offFood); err != nil {
		s.logger.Warn("failed to cache food in database", "error", err)
	}

	return offFood, nil
}

func (s *BarcodeService) fetchFromOpenFoodFacts(ctx context.Context, barcode string) (*domain.FoodBarcodeResponse, error) {
	tr := otel.Tracer("services/barcode_service.go")
	ctx, span := tr.Start(ctx, "fetchFromOpenFoodFacts")
	defer span.End()

	url := fmt.Sprintf("%s/product/%s.json", s.openFoodFactsURL, barcode)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("User-Agent", "CalorieAI/1.0")

	resp, err := s.httpClient.Do(req)
	if err != nil {
		s.logger.Error("failed to call OpenFoodFacts", "error", err)
		return nil, fmt.Errorf("failed to call OpenFoodFacts: %w", err)
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			s.logger.Warn("failed to close response body", "error", err)
		}
	}()

	if resp.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("barcode not found")
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("OpenFoodFacts returned status %d", resp.StatusCode)
	}

	var offResp struct {
		Status  int `json:"status"`
		Product struct {
			ProductName       string `json:"product_name"`
			Brands            string `json:"brands"`
			NutrimentsPer100g struct {
				EnergyKcal100g float64 `json:"energy-kcal_100g"`
				Proteins100g   float64 `json:"proteins_100g"`
				Carbs100g      float64 `json:"carbohydrates_100g"`
				Fat100g        float64 `json:"fat_100g"`
			} `json:"nutriments"`
		} `json:"product"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&offResp); err != nil {
		s.logger.Error("failed to decode OpenFoodFacts response", "error", err)
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	if offResp.Status != 1 {
		return nil, fmt.Errorf("product not found in OpenFoodFacts")
	}

	source := "OpenFoodFacts"
	brand := offResp.Product.Brands
	servingSize := int32(100)
	servingUnit := "g"

	return &domain.FoodBarcodeResponse{
		Barcode:     barcode,
		Name:        offResp.Product.ProductName,
		Brand:       &brand,
		Calories:    int32(offResp.Product.NutrimentsPer100g.EnergyKcal100g),
		Protein:     int32(offResp.Product.NutrimentsPer100g.Proteins100g),
		Carbs:       int32(offResp.Product.NutrimentsPer100g.Carbs100g),
		Fat:         int32(offResp.Product.NutrimentsPer100g.Fat100g),
		ServingSize: &servingSize,
		ServingUnit: &servingUnit,
		Source:      &source,
	}, nil
}

func (s *BarcodeService) cacheFoodInDB(ctx context.Context, barcode string, food *domain.FoodBarcodeResponse) error {
	tr := otel.Tracer("services/barcode_service.go")
	ctx, span := tr.Start(ctx, "cacheFoodInDB")
	defer span.End()

	params := db.CreateFoodFromBarcodeParams{
		Barcode:  &barcode,
		Name:     food.Name,
		Calories: intPtrFromInt32Ptr(&food.Calories),
		Protein:  new(float64(food.Protein)),
		Carbs:    new(float64(food.Carbs)),
		Fat:      new(float64(food.Fat)),
	}

	if food.Brand != nil {
		params.Brand = food.Brand
	}
	if food.ServingSize != nil {
		params.ServingSize = intPtrFromInt32Ptr(food.ServingSize)
	}
	if food.ServingUnit != nil {
		params.ServingUnit = food.ServingUnit
	}
	source := "OpenFoodFacts"
	params.Source = &source

	_, err := s.foodRepo.CreateFoodFromBarcode(ctx, params)
	return err
}

func (s *BarcodeService) mapDBToResponse(ctx context.Context, food *db.FoodDatabase) *domain.FoodBarcodeResponse {
	tr := otel.Tracer("services/barcode_service.go")
	_, span := tr.Start(ctx, "mapDBToResponse")
	defer span.End()

	source := "Database"
	var barcode string
	if food.Barcode != nil {
		barcode = *food.Barcode
	}

	calories := int32(0)
	if food.Calories != nil {
		calories = int32(*food.Calories)
	}

	var protein int32
	if food.Protein != nil {
		protein = int32(*food.Protein)
	}
	var carbs int32
	if food.Carbs != nil {
		carbs = int32(*food.Carbs)
	}
	var fat int32
	if food.Fat != nil {
		fat = int32(*food.Fat)
	}

	var servingSize *int32
	if food.ServingSize != nil {
		ss := int32(*food.ServingSize)
		servingSize = &ss
	}

	return &domain.FoodBarcodeResponse{
		Barcode:     barcode,
		Name:        food.Name,
		Brand:       food.Brand,
		Calories:    calories,
		Protein:     protein,
		Carbs:       carbs,
		Fat:         fat,
		ServingSize: servingSize,
		ServingUnit: food.ServingUnit,
		Source:      &source,
	}
}

func intPtrFromInt32Ptr(i *int32) *int {
	if i == nil {
		return nil
	}
	val := int(*i)
	return &val
}
