package services

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/jeancarloshp/calorieai/internal/domain"
	"go.opentelemetry.io/otel"
)

type GeminiProvider struct {
	apiKey     string
	httpClient *http.Client
	logger     domain.Logger
}

func NewGeminiProvider(apiKey string, logger domain.Logger) *GeminiProvider {
	return &GeminiProvider{
		apiKey: apiKey,
		httpClient: &http.Client{
			Timeout: 60 * time.Second,
		},
		logger: logger,
	}
}

func (g *GeminiProvider) RecognizeFood(
	ctx context.Context,
	imageBase64 string,
) ([]domain.RecognizedFoodItem, error) {
	tr := otel.Tracer("services/gemini_provider.go")
	ctx, span := tr.Start(ctx, "RecognizeFood")
	defer span.End()

	prompt := `Analyze this food image and return a JSON array of food items with their nutritional information.
	For each food item, provide:
	- name: food name in English (string)
	- calories: estimated calories (number, can be decimal)
	- protein: protein in grams (number, can be decimal)
	- carbs: carbohydrates in grams (number, can be decimal)
	- fat: fat in grams (number, can be decimal)
	- quantity: estimated quantity (number, can be decimal)
	- unit: unit of measurement (string: g, ml, or serving)
	- confidence: confidence score between 0 and 1 (number)
	
	Use the following JSON schema:
	{
	  "type": "object",
	  "properties": {
	    "food_items": {
	      "type": "array",
	      "items": {
	        "type": "object",
	        "properties": {
	          "name": {"type": "string"},
	          "calories": {"type": "number"},
	          "protein": {"type": "number"},
	          "carbs": {"type": "number"},
	          "fat": {"type": "number"},
	          "quantity": {"type": "number"},
	          "unit": {"type": "string"},
	          "confidence": {"type": "number", "minimum": 0, "maximum": 1}
	        },
	        "required": ["name", "calories", "protein", "carbs", "fat", "quantity", "unit", "confidence"]
	      }
	    }
	  },
	  "required": ["food_items"]
	}
	
	Return ONLY valid JSON matching this schema.`

	requestBody := map[string]interface{}{
		"contents": []map[string]interface{}{
			{
				"parts": []map[string]interface{}{
					{"text": prompt},
					{
						"inline_data": map[string]string{
							"mime_type": "image/jpeg",
							"data":      imageBase64,
						},
					},
				},
			},
		},
		"generationConfig": map[string]interface{}{
			"temperature":     0.4,
			"topK":            32,
			"topP":            1,
			"maxOutputTokens": 2048,
		},
	}

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	url := fmt.Sprintf("https://generativelanguage.googleapis.com/v1beta/models/gemini-3-flash-preview:generateContent?key=%s", g.apiKey)
	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := g.httpClient.Do(req)
	if err != nil {
		g.logger.Error("failed to call Gemini API", "error", err)
		return nil, fmt.Errorf("failed to call Gemini API: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		g.logger.Error("Gemini API returned error", "statusCode", resp.StatusCode, "body", string(bodyBytes))
		return nil, fmt.Errorf("gemini API returned status %d", resp.StatusCode)
	}

	var geminiResp struct {
		Candidates []struct {
			Content struct {
				Parts []struct {
					Text string `json:"text"`
				} `json:"parts"`
			} `json:"content"`
		} `json:"candidates"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&geminiResp); err != nil {
		g.logger.Error("failed to decode Gemini response", "error", err)
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	if len(geminiResp.Candidates) == 0 || len(geminiResp.Candidates[0].Content.Parts) == 0 {
		return nil, fmt.Errorf("no response from Gemini")
	}

	responseText := geminiResp.Candidates[0].Content.Parts[0].Text

	// Remove markdown code block formatting if present
	responseText = strings.TrimPrefix(responseText, "```json\n")
	responseText = strings.TrimSuffix(responseText, "\n```")

	var foodData struct {
		FoodItems []domain.RecognizedFoodItem `json:"food_items"`
	}

	if err := json.Unmarshal([]byte(responseText), &foodData); err != nil {
		g.logger.Error("failed to parse food items", "error", err, "text", responseText)
		return nil, fmt.Errorf("failed to parse food items: %w", err)
	}

	return foodData.FoodItems, nil
}

func (g *GeminiProvider) EstimateQuantity(
	ctx context.Context,
	imageBase64 string,
	req *domain.EstimateQuantityRequest,
) (*domain.EstimateQuantityResponse, error) {
	tr := otel.Tracer("services/gemini_provider.go")
	ctx, span := tr.Start(ctx, "EstimateQuantity")
	defer span.End()

	prompt := fmt.Sprintf(`Analyze this image to estimate the quantity of %s.
	Consider the context: meal location is %s, meal type is %s.
	
	Return ONLY valid JSON in this format:
	{
		"estimatedQuantity": <number>,
		"unit": "<g|ml|serving>",
		"confidence": <0-1>,
		"reasoning": "Brief explanation of the estimation"
	}`, req.Name, req.MealLocation, req.Type)

	if req.ReferenceServingSize != nil && req.ReferenceServingUnit != nil {
		prompt += fmt.Sprintf("\nReference serving: %s %s", *req.ReferenceServingSize, *req.ReferenceServingUnit)
	}

	requestBody := map[string]interface{}{
		"contents": []map[string]interface{}{
			{
				"parts": []map[string]interface{}{
					{"text": prompt},
					{
						"inline_data": map[string]string{
							"mime_type": "image/jpeg",
							"data":      imageBase64,
						},
					},
				},
			},
		},
		"generationConfig": map[string]interface{}{
			"temperature":     0.4,
			"topK":            32,
			"topP":            1,
			"maxOutputTokens": 1024,
		},
	}

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	url := fmt.Sprintf("https://generativelanguage.googleapis.com/v1beta/models/gemini-3-flash-preview:generateContent?key=%s", g.apiKey)
	httpReq, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := g.httpClient.Do(httpReq)
	if err != nil {
		g.logger.Error("failed to call Gemini API", "error", err)
		return nil, fmt.Errorf("failed to call Gemini API: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		g.logger.Error("Gemini API returned error", "statusCode", resp.StatusCode, "body", string(bodyBytes))
		return nil, fmt.Errorf("gemini API returned status %d", resp.StatusCode)
	}

	var geminiResp struct {
		Candidates []struct {
			Content struct {
				Parts []struct {
					Text string `json:"text"`
				} `json:"parts"`
			} `json:"content"`
		} `json:"candidates"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&geminiResp); err != nil {
		g.logger.Error("failed to decode Gemini response", "error", err)
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	if len(geminiResp.Candidates) == 0 || len(geminiResp.Candidates[0].Content.Parts) == 0 {
		return nil, fmt.Errorf("no response from Gemini")
	}

	responseText := geminiResp.Candidates[0].Content.Parts[0].Text

	// Remove markdown code block formatting if present
	responseText = strings.TrimPrefix(responseText, "```json\n")
	responseText = strings.TrimSuffix(responseText, "\n```")

	var result domain.EstimateQuantityResponse
	if err := json.Unmarshal([]byte(responseText), &result); err != nil {
		g.logger.Error("failed to parse quantity estimation", "error", err, "text", responseText)
		return nil, fmt.Errorf("failed to parse quantity estimation: %w", err)
	}

	return &result, nil
}
