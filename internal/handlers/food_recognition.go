package handlers

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/jeancarloshp/calorieai/internal/domain"
	"github.com/jeancarloshp/calorieai/internal/services"
	"github.com/rs/zerolog/log"
)

type FoodRecognitionHandler struct {
	foodRecognitionService *services.FoodRecognitionService
	barcodeService         *services.BarcodeService
	validator              *validator.Validate
}

func NewFoodRecognitionHandler(
	foodRecognitionService *services.FoodRecognitionService,
	barcodeService *services.BarcodeService,
) *FoodRecognitionHandler {
	return &FoodRecognitionHandler{
		foodRecognitionService: foodRecognitionService,
		barcodeService:         barcodeService,
		validator:              validator.New(),
	}
}

func (h *FoodRecognitionHandler) RecognizeFood(c *fiber.Ctx) error {
	ctx := context.Background()

	fileHeader, err := c.FormFile("image")
	if err != nil {
		log.Error().Err(err).Msg("Failed to get image from form")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "image file is required",
		})
	}

	req := &domain.FoodRecognitionRequest{
		Name:         c.FormValue("name"),
		Type:         c.FormValue("type"),
		MealLocation: c.FormValue("mealLocation"),
		URI:          c.FormValue("uri"),
	}

	if err := h.validator.Struct(req); err != nil {
		log.Error().Err(err).Msg("Validation failed")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "name, type, and mealLocation are required",
		})
	}

	// Configurar SSE headers
	c.Set("Content-Type", "text/event-stream")
	c.Set("Cache-Control", "no-cache")
	c.Set("Connection", "keep-alive")
	c.Set("Transfer-Encoding", "chunked")

	// Canal para receber progresso
	progressChan := make(chan domain.ProgressUpdate, 10)
	resultChan := make(chan *domain.FoodRecognitionResponse, 1)
	errorChan := make(chan error, 1)

	// Processar em goroutine
	go func() {
		result, err := h.foodRecognitionService.RecognizeFoodWithProgress(ctx, fileHeader, req, progressChan)
		if err != nil {
			errorChan <- err
			return
		}
		resultChan <- result
	}()

	// Enviar eventos SSE
	c.Context().SetBodyStreamWriter(func(w *bufio.Writer) {
		for {
			select {
			case progress := <-progressChan:
				data, _ := json.Marshal(progress)
				fmt.Fprintf(w, "event: progress\ndata: %s\n\n", data)
				w.Flush()

			case result := <-resultChan:
				data, _ := json.Marshal(fiber.Map{
					"success": true,
					"data":    result,
					"message": "food recognized successfully",
				})
				fmt.Fprintf(w, "event: complete\ndata: %s\n\n", data)
				w.Flush()
				return

			case err := <-errorChan:
				log.Error().Err(err).Msg("Failed to recognize food")
				data, _ := json.Marshal(fiber.Map{
					"success": false,
					"message": "failed to recognize food",
				})
				fmt.Fprintf(w, "event: error\ndata: %s\n\n", data)
				w.Flush()
				return

			case <-ctx.Done():
				return
			}
		}
	})

	return nil
}

func (h *FoodRecognitionHandler) GetFoodByBarcode(c *fiber.Ctx) error {
	ctx := context.Background()

	barcode := c.Params("barcode")
	if barcode == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "barcode is required",
		})
	}

	result, err := h.barcodeService.GetFoodByBarcode(ctx, barcode)
	if err != nil {
		log.Error().Err(err).Str("barcode", barcode).Msg("Failed to get food by barcode")
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"message": "food not found for barcode",
		})
	}

	if err := h.validator.Struct(result); err != nil {
		log.Error().Err(err).Msg("Validation failed for barcode response")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "invalid response from barcode service",
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    result,
		"message": "food retrieved successfully",
	})
}

func (h *FoodRecognitionHandler) EstimateQuantity(c *fiber.Ctx) error {
	ctx := context.Background()

	fileHeader, err := c.FormFile("image")
	if err != nil {
		log.Error().Err(err).Msg("Failed to get image from form")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "image file is required",
		})
	}

	req := &domain.EstimateQuantityRequest{
		Name:         c.FormValue("name"),
		Type:         c.FormValue("type"),
		MealLocation: c.FormValue("mealLocation"),
		URI:          c.FormValue("uri"),
	}

	if refSize := c.FormValue("referenceServingSize"); refSize != "" {
		req.ReferenceServingSize = &refSize
	}
	if refUnit := c.FormValue("referenceServingUnit"); refUnit != "" {
		req.ReferenceServingUnit = &refUnit
	}

	if err := h.validator.Struct(req); err != nil {
		log.Error().Err(err).Msg("Validation failed")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "name, type, and mealLocation are required",
		})
	}

	// Configurar SSE headers
	c.Set("Content-Type", "text/event-stream")
	c.Set("Cache-Control", "no-cache")
	c.Set("Connection", "keep-alive")
	c.Set("Transfer-Encoding", "chunked")

	// Canal para receber progresso
	progressChan := make(chan domain.ProgressUpdate, 10)
	resultChan := make(chan *domain.EstimateQuantityResponse, 1)
	errorChan := make(chan error, 1)

	// Processar em goroutine
	go func() {
		result, err := h.foodRecognitionService.EstimateQuantityWithProgress(ctx, fileHeader, req, progressChan)
		if err != nil {
			errorChan <- err
			return
		}
		resultChan <- result
	}()

	// Enviar eventos SSE
	c.Context().SetBodyStreamWriter(func(w *bufio.Writer) {
		for {
			select {
			case progress := <-progressChan:
				data, _ := json.Marshal(progress)
				fmt.Fprintf(w, "event: progress\ndata: %s\n\n", data)
				w.Flush()

			case result := <-resultChan:
				data, _ := json.Marshal(fiber.Map{
					"success": true,
					"data":    result,
					"message": "quantity estimated successfully",
				})
				fmt.Fprintf(w, "event: complete\ndata: %s\n\n", data)
				w.Flush()
				return

			case err := <-errorChan:
				log.Error().Err(err).Msg("Failed to estimate quantity")
				data, _ := json.Marshal(fiber.Map{
					"success": false,
					"message": "failed to estimate quantity",
				})
				fmt.Fprintf(w, "event: error\ndata: %s\n\n", data)
				w.Flush()
				return

			case <-ctx.Done():
				return
			}
		}
	})

	return nil
}
