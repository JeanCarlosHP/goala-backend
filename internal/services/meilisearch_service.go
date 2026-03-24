package services

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jeancarloshp/calorieai/internal/domain"
)

type MeiliSearchService struct {
	baseURL    string
	apiKey     string
	indexName  string
	httpClient *http.Client
	logger     domain.Logger
}

func NewMeiliSearchService(cfg *domain.Config, logger domain.Logger) *MeiliSearchService {
	indexName := cfg.MeiliSearchFoodsIndex
	if indexName == "" {
		indexName = "foods"
	}

	return &MeiliSearchService{
		baseURL:   strings.TrimRight(cfg.MeiliSearchURL, "/"),
		apiKey:    cfg.MeiliSearchAPIKey,
		indexName: indexName,
		httpClient: &http.Client{
			Timeout: 3 * time.Second,
		},
		logger: logger,
	}
}

func (s *MeiliSearchService) Enabled() bool {
	return s.baseURL != ""
}

func (s *MeiliSearchService) SearchFoodIDs(ctx context.Context, query string, limit int) ([]uuid.UUID, error) {
	if !s.Enabled() {
		return nil, nil
	}

	body, _ := json.Marshal(map[string]any{
		"q":     query,
		"limit": limit,
	})

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		fmt.Sprintf("%s/indexes/%s/search", s.baseURL, s.indexName),
		bytes.NewReader(body),
	)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	if s.apiKey != "" {
		req.Header.Set("Authorization", "Bearer "+s.apiKey)
	}

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode >= 300 {
		return nil, fmt.Errorf("meilisearch returned status %d", resp.StatusCode)
	}

	var payload struct {
		Hits []struct {
			ID string `json:"id"`
		} `json:"hits"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		return nil, err
	}

	ids := make([]uuid.UUID, 0, len(payload.Hits))
	for _, hit := range payload.Hits {
		parsed, err := uuid.Parse(hit.ID)
		if err == nil {
			ids = append(ids, parsed)
		}
	}

	return ids, nil
}

func (s *MeiliSearchService) IndexFood(ctx context.Context, food domain.SearchFood) {
	if !s.Enabled() || food.ID == nil {
		return
	}

	body, _ := json.Marshal([]map[string]any{{
		"id":       food.ID.String(),
		"name":     food.Name,
		"brand":    food.Brand,
		"source":   food.Source,
		"calories": food.Calories,
		"protein":  food.Protein,
		"carbs":    food.Carbs,
		"fat":      food.Fat,
	}})

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		fmt.Sprintf("%s/indexes/%s/documents", s.baseURL, s.indexName),
		bytes.NewReader(body),
	)
	if err != nil {
		return
	}
	req.Header.Set("Content-Type", "application/json")
	if s.apiKey != "" {
		req.Header.Set("Authorization", "Bearer "+s.apiKey)
	}

	resp, err := s.httpClient.Do(req)
	if err != nil {
		s.logger.Warn("failed to index food in meilisearch", "error", err)
		return
	}
	defer func() { _ = resp.Body.Close() }()
}
