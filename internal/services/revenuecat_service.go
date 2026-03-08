package services

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"

	"github.com/jeancarloshp/calorieai/internal/domain"
	"go.opentelemetry.io/otel"
)

type RevenueCatService struct {
	webhookSecret string
	logger        domain.Logger
}

func NewRevenueCatService(webhookSecret string, log domain.Logger) *RevenueCatService {
	return &RevenueCatService{
		webhookSecret: webhookSecret,
		logger:        log,
	}
}

func (s *RevenueCatService) VerifyWebhookSignature(ctx context.Context, payload []byte, signature string) error {
	tr := otel.Tracer("services/revenuecat_service.go")
	_, span := tr.Start(ctx, "VerifyWebhookSignature")
	defer span.End()

	if s.webhookSecret == "" {
		s.logger.Warn("Webhook secret not configured, skipping signature verification", nil)
		return nil
	}

	mac := hmac.New(sha256.New, []byte(s.webhookSecret))
	mac.Write(payload)
	expectedSignature := hex.EncodeToString(mac.Sum(nil))

	if !hmac.Equal([]byte(signature), []byte(expectedSignature)) {
		return fmt.Errorf("invalid webhook signature")
	}

	return nil
}

func (s *RevenueCatService) ParseWebhook(ctx context.Context, payload []byte) (*domain.RevenueCatWebhook, error) {
	tr := otel.Tracer("services/revenuecat_service.go")
	_, span := tr.Start(ctx, "ParseWebhook")
	defer span.End()

	var webhook domain.RevenueCatWebhook
	if err := json.Unmarshal(payload, &webhook); err != nil {
		return nil, fmt.Errorf("failed to parse webhook: %w", err)
	}

	if !webhook.Event.Type.IsValid() {
		return nil, fmt.Errorf("invalid event type: %s", webhook.Event.Type)
	}

	return &webhook, nil
}
