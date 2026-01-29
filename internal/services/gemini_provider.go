package services

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/jeancarloshp/calorieai/internal/domain"
	"go.opentelemetry.io/otel"
	"google.golang.org/genai"
)

type GeminiProvider struct {
	client *genai.Client
	model  string
	logger domain.Logger
}

func NewGeminiProvider(apiKey string, model string, logger domain.Logger) *GeminiProvider {
	if model == "" {
		model = "gemini-3-flash-preview"
	}
	ctx := context.Background()
	client, err := genai.NewClient(ctx, &genai.ClientConfig{APIKey: apiKey})
	if err != nil {
		logger.Error("failed to create Gemini client", "error", err)
		panic(fmt.Sprintf("failed to create Gemini client: %v", err))
	}
	return &GeminiProvider{
		client: client,
		model:  model,
		logger: logger,
	}
}

func (g *GeminiProvider) RecognizeFood(
	ctx context.Context,
	imageBytes []byte,
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
	          "unit": {"type": "string"}
	        },
	        "required": ["name", "calories", "protein", "carbs", "fat", "quantity", "unit"]
	      }
	    }
	  },
	  "required": ["food_items"]
	}
	
	Return ONLY valid, complete JSON matching this schema. Ensure the JSON is properly closed and valid.`

	parts := []*genai.Part{
		genai.NewPartFromBytes(imageBytes, "image/jpeg"),
		genai.NewPartFromText(prompt),
	}

	contents := []*genai.Content{
		genai.NewContentFromParts(parts, genai.RoleUser),
	}

	config := &genai.GenerateContentConfig{
		Temperature:     genai.Ptr(float32(0.4)),
		TopK:            genai.Ptr(float32(32)),
		TopP:            genai.Ptr(float32(1)),
		MaxOutputTokens: 4096,
	}

	result, err := g.client.Models.GenerateContent(ctx, g.model, contents, config)
	if err != nil {
		g.logger.Error("failed to call Gemini API", "error", err)
		return nil, fmt.Errorf("failed to call Gemini API: %w", err)
	}

	responseText := result.Text()

	// Remove markdown code block formatting if present
	responseText = strings.TrimPrefix(responseText, "```json\n")
	responseText = strings.TrimSuffix(responseText, "\n```")

	// Check if the response is valid JSON
	if !json.Valid([]byte(responseText)) {
		g.logger.Error("invalid JSON response from Gemini", "text", responseText)
		return nil, fmt.Errorf("invalid JSON response from AI provider")
	}

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
	imageBytes []byte,
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

	parts := []*genai.Part{
		genai.NewPartFromBytes(imageBytes, "image/jpeg"),
		genai.NewPartFromText(prompt),
	}

	contents := []*genai.Content{
		genai.NewContentFromParts(parts, genai.RoleUser),
	}

	config := &genai.GenerateContentConfig{
		Temperature:     genai.Ptr(float32(0.4)),
		TopK:            genai.Ptr(float32(32)),
		TopP:            genai.Ptr(float32(1)),
		MaxOutputTokens: 1024,
	}

	result, err := g.client.Models.GenerateContent(ctx, g.model, contents, config)
	if err != nil {
		g.logger.Error("failed to call Gemini API", "error", err)
		return nil, fmt.Errorf("failed to call Gemini API: %w", err)
	}

	responseText := result.Text()

	// Remove markdown code block formatting if present
	responseText = strings.TrimPrefix(responseText, "```json\n")
	responseText = strings.TrimSuffix(responseText, "\n```")

	var resultResp domain.EstimateQuantityResponse
	if err := json.Unmarshal([]byte(responseText), &resultResp); err != nil {
		g.logger.Error("failed to parse quantity estimation", "error", err, "text", responseText)
		return nil, fmt.Errorf("failed to parse quantity estimation: %w", err)
	}

	return &resultResp, nil
}
