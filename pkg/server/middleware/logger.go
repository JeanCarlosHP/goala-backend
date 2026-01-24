package middleware

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

func RequestLogger() fiber.Handler {
	return func(c *fiber.Ctx) error {
		start := time.Now()
		requestID := uuid.New().String()

		c.Locals("requestID", requestID)

		logger := log.With().
			Str("request_id", requestID).
			Str("method", c.Method()).
			Str("path", c.Path()).
			Str("ip", c.IP()).
			Logger()

		c.Locals("logger", &logger)

		err := c.Next()

		duration := time.Since(start)
		status := c.Response().StatusCode()

		logEvent := logger.Info()
		if status >= 400 && status < 500 {
			logEvent = logger.Warn()
		} else if status >= 500 {
			logEvent = logger.Error()
		}

		logEvent.
			Int("status", status).
			Dur("duration_ms", duration).
			Int("response_size", len(c.Response().Body())).
			Msg("Request completed")

		return err
	}
}
