package services

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/jeancarloshp/calorieai/internal/domain"
	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/option"
	"go.opentelemetry.io/otel"
)

type OpenAIProvider struct {
	client openai.Client
	model  string
	logger domain.Logger
}

func NewOpenAIProvider(apiKey string, model string, logger domain.Logger) *OpenAIProvider {
	if model == "" {
		model = "gpt-4o"
	}
	client := openai.NewClient(
		option.WithAPIKey(apiKey),
	)
	return &OpenAIProvider{
		client: client,
		model:  model,
		logger: logger,
	}
}

func (o *OpenAIProvider) RecognizeFood(
	ctx context.Context,
	imageBase64 string,
) ([]domain.RecognizedFoodItem, error) {
	tr := otel.Tracer("services/openai_provider.go")
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

	Return ONLY valid, complete JSON matching this schema. Ensure the JSON is properly closed and valid.`

	imageURL := fmt.Sprintf("data:image/jpeg;base64,%s", imageBase64)

	messages := []openai.ChatCompletionMessageParamUnion{
		openai.UserMessage([]openai.ChatCompletionContentPartUnionParam{
			openai.TextContentPart(prompt),
			openai.ImageContentPart(openai.ChatCompletionContentPartImageImageURLParam{
				URL: imageURL,
			}),
		}),
	}

	chat, err := o.client.Chat.Completions.New(ctx, openai.ChatCompletionNewParams{
		Messages: messages,
		Model:    o.model,
	})

	if err != nil {
		o.logger.Error("failed to call OpenAI API", "error", err)
		return nil, fmt.Errorf("failed to call OpenAI API: %w", err)
	}

	responseText := chat.Choices[0].Message.Content

	// Remove markdown code block formatting if present
	responseText = strings.TrimPrefix(responseText, "```json\n")
	responseText = strings.TrimSuffix(responseText, "\n```")

	// Check if the response is valid JSON
	if !json.Valid([]byte(responseText)) {
		o.logger.Error("invalid JSON response from OpenAI", "text", responseText)
		return nil, fmt.Errorf("invalid JSON response from AI provider")
	}

	var foodData struct {
		FoodItems []domain.RecognizedFoodItem `json:"food_items"`
	}

	if err := json.Unmarshal([]byte(responseText), &foodData); err != nil {
		o.logger.Error("failed to parse food items", "error", err, "text", responseText)
		return nil, fmt.Errorf("failed to parse food items: %w", err)
	}

	return foodData.FoodItems, nil
}

func (o *OpenAIProvider) EstimateQuantity(
	ctx context.Context,
	imageBase64 string,
	req *domain.EstimateQuantityRequest,
) (*domain.EstimateQuantityResponse, error) {
	tr := otel.Tracer("services/openai_provider.go")
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

	imageURL := fmt.Sprintf("data:image/jpeg;base64,%s", imageBase64)

	messages := []openai.ChatCompletionMessageParamUnion{
		openai.UserMessage([]openai.ChatCompletionContentPartUnionParam{
			openai.TextContentPart(prompt),
			openai.ImageContentPart(openai.ChatCompletionContentPartImageImageURLParam{
				URL: imageURL,
			}),
		}),
	}

	chat, err := o.client.Chat.Completions.New(ctx, openai.ChatCompletionNewParams{
		Messages: messages,
		Model:    o.model,
	})

	if err != nil {
		o.logger.Error("failed to call OpenAI API", "error", err)
		return nil, fmt.Errorf("failed to call OpenAI API: %w", err)
	}

	responseText := chat.Choices[0].Message.Content

	// Remove markdown code block formatting if present
	responseText = strings.TrimPrefix(responseText, "```json\n")
	responseText = strings.TrimSuffix(responseText, "\n```")

	var resultResp domain.EstimateQuantityResponse
	if err := json.Unmarshal([]byte(responseText), &resultResp); err != nil {
		o.logger.Error("failed to parse quantity estimation", "error", err, "text", responseText)
		return nil, fmt.Errorf("failed to parse quantity estimation: %w", err)
	}

	return &resultResp, nil
}
