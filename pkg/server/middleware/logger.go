package middleware

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/jeancarloshp/calorieai/internal/domain"
	"go.opentelemetry.io/otel"
)

func RequestLogger(logger domain.Logger) fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx := c.UserContext()

		tr := otel.Tracer("middleware/logger.go")
		ctx, span := tr.Start(ctx, "RequestLogger")
		defer span.End()

		start := time.Now()
		requestID := uuid.New().String()

		c.Locals("requestID", requestID)

		c.Locals("logger", &logger)

		err := c.Next()

		duration := time.Since(start)
		status := c.Response().StatusCode()

		if status >= 400 && status < 500 {
			logger.Warn(
				"Client error occurred",
				"request_id", requestID,
				"method", c.Method(),
				"path", c.Path(),
				"ip", c.IP(),
			)
		} else if status >= 500 {
			logger.Error(
				"Server error occurred",
				"request_id", requestID,
				"method", c.Method(),
				"path", c.Path(),
				"ip", c.IP(),
			)
		}

		logger.Info(
			"Request completed",
			"request_id", requestID,
			"method", c.Method(),
			"path", c.Path(),
			"status", status,
			"duration_ms", duration.Milliseconds(),
			"response_size", len(c.Response().Body()),
			"ip", c.IP(),
		)

		return err
	}
}
