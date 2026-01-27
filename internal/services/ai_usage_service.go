package services

import (
	"context"
	"fmt"
	"time"

	"github.com/jeancarloshp/calorieai/internal/domain"
	"github.com/jeancarloshp/calorieai/internal/domain/enum"
	"github.com/jeancarloshp/calorieai/internal/repositories"
	"go.opentelemetry.io/otel"
)

type AIUsageService struct {
	usageRepo *repositories.AIUsageRepository
	subRepo   *repositories.SubscriptionRepository
	logger    domain.Logger
}

func NewAIUsageService(
	usageRepo *repositories.AIUsageRepository,
	subRepo *repositories.SubscriptionRepository,
	log domain.Logger,
) *AIUsageService {
	return &AIUsageService{
		usageRepo: usageRepo,
		subRepo:   subRepo,
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
		Free:    10,
		Monthly: 100,
		Yearly:  100,
		Trial:   50,
	},
	enum.FeatureMealAnalysis: {
		Free:    5,
		Monthly: 50,
		Yearly:  50,
		Trial:   25,
	},
	enum.FeatureNutritionAdvice: {
		Free:    3,
		Monthly: 30,
		Yearly:  30,
		Trial:   15,
	},
}

func (s *AIUsageService) CheckAndIncrementUsage(ctx context.Context, userID string, feature enum.AIFeature) error {
	tr := otel.Tracer("services/ai_usage_service.go")
	ctx, span := tr.Start(ctx, "CheckAndIncrementUsage")
	defer span.End()

	sub, err := s.subRepo.GetByUserID(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to get subscription: %w", err)
	}

	if sub == nil {
		sub = &domain.Subscription{
			UserID:   userID,
			IsActive: false,
			Plan:     enum.PlanFree,
			IsTrial:  false,
		}
	}

	quota := s.getQuotaForPlan(ctx, feature, sub)
	periodStart, periodEnd := s.getCurrentPeriod(ctx, sub)

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

func (s *AIUsageService) getQuotaForPlan(ctx context.Context, feature enum.AIFeature, sub *domain.Subscription) int32 {
	tr := otel.Tracer("services/ai_usage_service.go")
	ctx, span := tr.Start(ctx, "getQuotaForPlan")
	defer span.End()

	quotaConfig, exists := featureQuotas[feature]
	if !exists {
		return 0
	}

	if sub.IsTrial {
		return quotaConfig.Trial
	}

	switch sub.Plan {
	case enum.PlanMonthly:
		return quotaConfig.Monthly
	case enum.PlanYearly:
		return quotaConfig.Yearly
	case enum.PlanFree:
		return quotaConfig.Free
	default:
		return quotaConfig.Free
	}
}

func (s *AIUsageService) getCurrentPeriod(ctx context.Context, sub *domain.Subscription) (time.Time, time.Time) {
	tr := otel.Tracer("services/ai_usage_service.go")
	ctx, span := tr.Start(ctx, "getQuotaForPlan")
	defer span.End()

	now := time.Now()

	if sub.CurrentPeriodStart != nil && sub.CurrentPeriodEnd != nil {
		if now.Before(*sub.CurrentPeriodEnd) {
			return *sub.CurrentPeriodStart, *sub.CurrentPeriodEnd
		}
	}

	startOfDay := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
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
