package services

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jeancarloshp/calorieai/internal/domain"
	"github.com/jeancarloshp/calorieai/internal/domain/enum"
	"github.com/jeancarloshp/calorieai/internal/repositories"
	"go.opentelemetry.io/otel"
)

type AIUsageService struct {
	usageRepo *repositories.AIUsageRepository
	subRepo   *repositories.SubscriptionRepository
	userRepo  *repositories.UserRepository
	logger    domain.Logger
}

func NewAIUsageService(
	usageRepo *repositories.AIUsageRepository,
	subRepo *repositories.SubscriptionRepository,
	userRepo *repositories.UserRepository,
	log domain.Logger,
) *AIUsageService {
	return &AIUsageService{
		usageRepo: usageRepo,
		subRepo:   subRepo,
		userRepo:  userRepo,
		logger:    log,
	}
}

type QuotaConfig struct {
	Free    int32
	Monthly int32
	Yearly  int32
	Trial   int32
}

var featureQuotas = map[enum.AIFeature]QuotaConfig{
	enum.FeatureFoodRecognition: {
		Free: 10, // Agora é o limite diário fixo
	},
	enum.FeatureMealAnalysis: {
		Free: 5,
	},
	enum.FeatureNutritionAdvice: {
		Free: 3,
	},
}

func (s *AIUsageService) CheckAndIncrementUsage(ctx context.Context, userID string, feature enum.AIFeature) error {
	tr := otel.Tracer("services/ai_usage_service.go")
	ctx, span := tr.Start(ctx, "CheckAndIncrementUsage")
	defer span.End()

	user, err := s.userRepo.GetByID(ctx, uuid.MustParse(userID))
	if err != nil {
		return fmt.Errorf("failed to get user: %w", err)
	}

	quota := s.getQuotaForFeature(ctx, feature)
	periodStart, periodEnd := s.getCurrentPeriod(ctx, user.Timezone)

	usage, err := s.usageRepo.GetByPeriod(ctx, userID, feature, periodStart)
	if err != nil {
		return fmt.Errorf("failed to get usage: %w", err)
	}

	if usage == nil || usage.IsExpired() {
		usage, err = s.usageRepo.CreateOrReset(ctx, userID, feature, quota, periodStart, periodEnd)
		if err != nil {
			return fmt.Errorf("failed to create usage record: %w", err)
		}
	}

	if !usage.HasQuota() {
		return &QuotaExceededError{
			Feature:    feature,
			UsageCount: usage.UsageCount,
			Quota:      usage.Quota,
		}
	}

	_, err = s.usageRepo.Increment(ctx, userID, feature, quota, periodStart, periodEnd)
	if err != nil {
		return fmt.Errorf("failed to increment usage: %w", err)
	}

	s.logger.Info("AI usage incremented", map[string]interface{}{
		"user_id": userID,
		"feature": feature,
		"count":   usage.UsageCount + 1,
		"quota":   quota,
	})

	return nil
}

func (s *AIUsageService) GetUsage(ctx context.Context, userID string, feature enum.AIFeature) (*domain.AIUsage, error) {
	tr := otel.Tracer("services/ai_usage_service.go")
	ctx, span := tr.Start(ctx, "GetUsage")
	defer span.End()

	return s.usageRepo.Get(ctx, userID, feature)
}

func (s *AIUsageService) ListUserUsage(ctx context.Context, userID string) ([]*domain.AIUsage, error) {
	tr := otel.Tracer("services/ai_usage_service.go")
	ctx, span := tr.Start(ctx, "ListUserUsage")
	defer span.End()

	return s.usageRepo.ListByUser(ctx, userID)
}

func (s *AIUsageService) getQuotaForFeature(ctx context.Context, feature enum.AIFeature) int32 {
	tr := otel.Tracer("services/ai_usage_service.go")
	_, span := tr.Start(ctx, "getQuotaForFeature")
	defer span.End()

	quotaConfig, exists := featureQuotas[feature]
	if !exists {
		return 0
	}

	return quotaConfig.Free // Agora sempre retorna o limite diário fixo
}

func (s *AIUsageService) getCurrentPeriod(ctx context.Context, timezone string) (time.Time, time.Time) {
	tr := otel.Tracer("services/ai_usage_service.go")
	_, span := tr.Start(ctx, "getCurrentPeriod")
	defer span.End()

	loc, err := time.LoadLocation(timezone)
	if err != nil {
		// Fallback para UTC se timezone inválida
		loc = time.UTC
	}

	now := time.Now().In(loc)
	startOfDay := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, loc)
	endOfDay := startOfDay.Add(24 * time.Hour)

	return startOfDay, endOfDay
}

type QuotaExceededError struct {
	Feature    enum.AIFeature
	UsageCount int32
	Quota      int32
}

func (e *QuotaExceededError) Error() string {
	return fmt.Sprintf("quota exceeded for feature %s: %d/%d", e.Feature, e.UsageCount, e.Quota)
}

func IsQuotaExceededError(err error) bool {
	_, ok := err.(*QuotaExceededError)
	return ok
}
