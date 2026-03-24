package services

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/jeancarloshp/calorieai/internal/repositories"
	"go.opentelemetry.io/otel"

	"github.com/google/uuid"
	"github.com/jeancarloshp/calorieai/internal/domain"
)

type FoodService struct {
	foodRepo         *repositories.FoodRepository
	cache            *RedisCacheService
	meili            *MeiliSearchService
	openFoodFactsURL string
	httpClient       *http.Client
	logger           domain.Logger
}

func NewFoodService(
	foodRepo *repositories.FoodRepository,
	cfg *domain.Config,
	cache *RedisCacheService,
	meili *MeiliSearchService,
	logger domain.Logger,
) *FoodService {
	openFoodFactsURL := strings.TrimRight(cfg.OpenFoodFactsAPIURL, "/")
	if openFoodFactsURL == "" {
		openFoodFactsURL = "https://world.openfoodfacts.org"
	}

	return &FoodService{
		foodRepo:         foodRepo,
		cache:            cache,
		meili:            meili,
		openFoodFactsURL: openFoodFactsURL,
		httpClient: &http.Client{
			Timeout: 8 * time.Second,
		},
		logger: logger,
	}
}

func (s *FoodService) SearchFoods(ctx context.Context, query string) ([]domain.FoodDatabase, error) {
	tr := otel.Tracer("services/food_service.go")
	ctx, span := tr.Start(ctx, "SearchFoods")
	defer span.End()

	if query == "" {
		return []domain.FoodDatabase{}, nil
	}

	return s.foodRepo.SearchFoodDatabase(ctx, query, 20)
}

func (s *FoodService) SearchFoodsManual(
	ctx context.Context,
	userID uuid.UUID,
	req domain.FoodSearchRequest,
) (*domain.FoodSearchResponse, error) {
	tr := otel.Tracer("services/food_service.go")
	ctx, span := tr.Start(ctx, "SearchFoodsManual")
	defer span.End()

	limit := req.Limit
	if limit <= 0 {
		limit = 20
	}

	response := &domain.FoodSearchResponse{
		Query:     strings.TrimSpace(req.Query),
		Foods:     []domain.SearchFood{},
		Recent:    []domain.SearchFood{},
		Favorites: []domain.SearchFood{},
	}

	recent, err := s.foodRepo.GetRecentSearchFoods(ctx, userID, 8)
	if err == nil {
		response.Recent = dedupeFoods(recent)
	}

	favorites, err := s.foodRepo.GetFavoriteFoods(ctx, userID, 8)
	if err == nil {
		response.Favorites = dedupeFoods(favorites)
	}

	if response.Query == "" {
		return response, nil
	}

	cacheKey := "foods:search:" + strings.ToLower(response.Query)
	if s.cache != nil && s.cache.GetJSON(ctx, cacheKey, response) {
		response.Recent = dedupeFoods(response.Recent)
		response.Favorites = dedupeFoods(response.Favorites)
		return response, nil
	}

	localFoods, err := s.searchLocalFoods(ctx, userID, response.Query, limit)
	if err != nil {
		return nil, err
	}

	response.Foods = dedupeFoods(localFoods)
	if len(response.Foods) < min(limit, 8) {
		fallbackFoods, err := s.searchOpenFoodFacts(ctx, response.Query, limit)
		if err == nil {
			response.Foods = dedupeFoods(append(response.Foods, fallbackFoods...))
		}
	}
	if len(response.Foods) > limit {
		response.Foods = response.Foods[:limit]
	}

	if s.cache != nil {
		s.cache.SetJSON(ctx, cacheKey, response, 10*time.Minute)
	}

	return response, nil
}

func (s *FoodService) ToggleFavorite(ctx context.Context, userID, foodID uuid.UUID, favorite bool) error {
	return s.foodRepo.ToggleFavorite(ctx, userID, foodID, favorite)
}

func (s *FoodService) EnsureCatalogFood(ctx context.Context, food domain.SearchFood) (*domain.SearchFood, error) {
	if food.ID != nil && food.Source == "internal" {
		return &food, nil
	}

	stored, err := s.foodRepo.UpsertFoodCatalogEntry(ctx, food)
	if err != nil {
		return nil, err
	}

	if s.meili != nil {
		s.meili.IndexFood(ctx, *stored)
	}

	return stored, nil
}

func (s *FoodService) searchLocalFoods(ctx context.Context, userID uuid.UUID, query string, limit int) ([]domain.SearchFood, error) {
	if s.meili != nil && s.meili.Enabled() {
		ids, err := s.meili.SearchFoodIDs(ctx, query, limit)
		if err == nil && len(ids) > 0 {
			foods, repoErr := s.foodRepo.SearchFoodsByIDs(ctx, ids, userID)
			if repoErr == nil && len(foods) > 0 {
				return foods, nil
			}
		}
	}

	return s.foodRepo.SearchFoodsForAutocomplete(ctx, query, limit, userID)
}

func (s *FoodService) searchOpenFoodFacts(ctx context.Context, query string, limit int) ([]domain.SearchFood, error) {
	requestURL := fmt.Sprintf(
		"%s/cgi/search.pl?search_terms=%s&search_simple=1&action=process&json=1&page_size=%d",
		s.openFoodFactsURL,
		url.QueryEscape(query),
		limit,
	)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, requestURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "CalorieAI/1.0")

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode >= 300 {
		return nil, fmt.Errorf("openfoodfacts returned status %d", resp.StatusCode)
	}

	var payload struct {
		Products []struct {
			Code        string `json:"code"`
			ProductName string `json:"product_name"`
			Brands      string `json:"brands"`
			ServingSize string `json:"serving_size"`
			Nutriments  struct {
				EnergyKcal100g float64 `json:"energy-kcal_100g"`
				Proteins100g   float64 `json:"proteins_100g"`
				Carbs100g      float64 `json:"carbohydrates_100g"`
				Fat100g        float64 `json:"fat_100g"`
			} `json:"nutriments"`
		} `json:"products"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		return nil, err
	}

	foods := make([]domain.SearchFood, 0, len(payload.Products))
	for _, product := range payload.Products {
		name := strings.TrimSpace(product.ProductName)
		if name == "" {
			continue
		}

		externalID := strings.TrimSpace(product.Code)
		brand := optionalTrimmedString(product.Brands)
		item := domain.SearchFood{
			ExternalID: optionalTrimmedString(externalID),
			Name:       name,
			Brand:      brand,
			Calories:   product.Nutriments.EnergyKcal100g,
			Protein:    product.Nutriments.Proteins100g,
			Carbs:      product.Nutriments.Carbs100g,
			Fat:        product.Nutriments.Fat100g,
			Source:     "openfoodfacts",
			Portions:   []domain.FoodPortion{{Name: "g", Grams: 1}},
		}

		if grams, label := parseServingSizePortion(product.ServingSize); grams > 0 && label != "" {
			item.Portions = append(item.Portions, domain.FoodPortion{
				Name:  label,
				Grams: grams,
			})
		}

		foods = append(foods, item)
	}

	return foods, nil
}

func dedupeFoods(items []domain.SearchFood) []domain.SearchFood {
	seen := make(map[string]struct{})
	result := make([]domain.SearchFood, 0, len(items))
	for _, item := range items {
		key := strings.ToLower(item.Name + "|" + stringValue(item.Brand) + "|" + item.Source)
		if item.ExternalID != nil && *item.ExternalID != "" {
			key = item.Source + "|" + *item.ExternalID
		}
		if _, exists := seen[key]; exists {
			continue
		}
		seen[key] = struct{}{}
		result = append(result, item)
	}
	return result
}

func parseServingSizePortion(raw string) (float64, string) {
	raw = strings.TrimSpace(strings.ToLower(raw))
	switch {
	case strings.HasSuffix(raw, "g"):
		value := strings.TrimSpace(strings.TrimSuffix(raw, "g"))
		return parsePositiveFloat(value), "serving"
	case strings.HasSuffix(raw, "ml"):
		value := strings.TrimSpace(strings.TrimSuffix(raw, "ml"))
		return parsePositiveFloat(value), "serving"
	default:
		return 0, ""
	}
}

func parsePositiveFloat(value string) float64 {
	parsed, err := strconv.ParseFloat(value, 64)
	if err != nil || parsed <= 0 {
		return 0
	}
	return parsed
}

func optionalTrimmedString(value string) *string {
	value = strings.TrimSpace(value)
	if value == "" {
		return nil
	}
	return &value
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func stringValue(value *string) string {
	if value == nil {
		return ""
	}
	return *value
}

func (s *FoodService) GetRecentFoods(ctx context.Context, userID uuid.UUID) ([]domain.RecentFood, error) {
	tr := otel.Tracer("services/food_service.go")
	ctx, span := tr.Start(ctx, "GetRecentFoods")
	defer span.End()

	return s.foodRepo.GetRecentFoods(ctx, userID, 20)
}

func (s *FoodService) CreateFoodItem(ctx context.Context, req *domain.CreateFoodItemRequest) (*domain.FoodItem, error) {
	tr := otel.Tracer("services/food_service.go")
	ctx, span := tr.Start(ctx, "CreateFoodItem")
	defer span.End()

	return s.foodRepo.CreateStandalone(ctx, req)
}

func (s *FoodService) GetFoodItem(ctx context.Context, id uuid.UUID) (*domain.FoodItem, error) {
	tr := otel.Tracer("services/food_service.go")
	ctx, span := tr.Start(ctx, "GetFoodItem")
	defer span.End()

	return s.foodRepo.GetByID(ctx, id)
}

func (s *FoodService) UpdateFoodItem(ctx context.Context, id uuid.UUID, req *domain.UpdateFoodItemRequest) (*domain.FoodItem, error) {
	tr := otel.Tracer("services/food_service.go")
	ctx, span := tr.Start(ctx, "UpdateFoodItem")
	defer span.End()

	_, err := s.foodRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return s.foodRepo.Update(ctx, id, req)
}

func (s *FoodService) DeleteFoodItem(ctx context.Context, id uuid.UUID) error {
	tr := otel.Tracer("services/food_service.go")
	ctx, span := tr.Start(ctx, "DeleteFoodItem")
	defer span.End()

	_, err := s.foodRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	return s.foodRepo.Delete(ctx, id)
}
