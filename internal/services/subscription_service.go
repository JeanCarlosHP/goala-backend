package services

import (
	"context"
	"fmt"
	"time"

	"github.com/jeancarloshp/calorieai/internal/domain"
	"github.com/jeancarloshp/calorieai/internal/domain/enum"
	"github.com/jeancarloshp/calorieai/internal/repositories"
)

type SubscriptionService struct {
	repo   *repositories.SubscriptionRepository
	logger domain.Logger
}

func NewSubscriptionService(repo *repositories.SubscriptionRepository, log domain.Logger) *SubscriptionService {
	return &SubscriptionService{
		repo:   repo,
		logger: log,
	}
}

func (s *SubscriptionService) GetOrCreateSubscription(ctx context.Context, userID string) (*domain.Subscription, error) {
	sub, err := s.repo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get subscription: %w", err)
	}

	if sub != nil {
		return sub, nil
	}

	newSub := &domain.Subscription{
		UserID:           userID,
		RevenueCatUserID: userID,
		IsActive:         false,
		Plan:             enum.PlanFree,
		IsTrial:          false,
	}

	created, err := s.repo.Create(ctx, newSub)
	if err != nil {
		return nil, fmt.Errorf("failed to create subscription: %w", err)
	}

	s.logger.Info("Created new free subscription", map[string]interface{}{
		"user_id": userID,
	})

	return created, nil
}

func (s *SubscriptionService) GetByUserID(ctx context.Context, userID string) (*domain.Subscription, error) {
	return s.repo.GetByUserID(ctx, userID)
}

func (s *SubscriptionService) ValidateAccess(ctx context.Context, userID string) (bool, error) {
	sub, err := s.GetOrCreateSubscription(ctx, userID)
	if err != nil {
		return false, err
	}

	if sub.IsTrial {
		return true, nil
	}

	if !sub.IsActive {
		return false, nil
	}

	if sub.IsExpired() {
		s.logger.Warn("Subscription expired", map[string]interface{}{
			"user_id":            userID,
			"current_period_end": sub.CurrentPeriodEnd,
		})
		return false, nil
	}

	return true, nil
}

func (s *SubscriptionService) ProcessWebhookEvent(ctx context.Context, event *domain.RevenueCatEvent) error {
	processed, err := s.repo.IsEventProcessed(ctx, event.ID)
	if err != nil {
		return fmt.Errorf("failed to check event status: %w", err)
	}

	if processed {
		s.logger.Info("Event already processed, skipping", map[string]interface{}{
			"event_id":   event.ID,
			"event_type": event.Type,
		})
		return nil
	}

	userID := event.AppUserID
	if userID == "" {
		userID = event.OriginalAppUserID
	}

	if userID == "" {
		return fmt.Errorf("no user ID in webhook event")
	}

	plan := s.mapProductIDToPlan(event.ProductID)
	isActive := s.isActiveEvent(event.Type)
	expirationAt := event.ExpirationAt()
	purchasedAt := event.PurchasedAt()

	subscription := &domain.Subscription{
		UserID:                          userID,
		RevenueCatUserID:                userID,
		RevenueCatOriginalTransactionID: event.OriginalTransactionID,
		IsActive:                        isActive,
		Plan:                            plan,
		IsTrial:                         event.IsTrial,
		CurrentPeriodStart:              &purchasedAt,
		CurrentPeriodEnd:                expirationAt,
		LastEventID:                     &event.ID,
		LastEventType:                   stringPtr(event.Type.String()),
		LastEventAt:                     timePtr(time.Now()),
	}

	_, err = s.repo.Upsert(ctx, subscription)
	if err != nil {
		return fmt.Errorf("failed to upsert subscription: %w", err)
	}

	s.logger.Info("Processed webhook event", map[string]interface{}{
		"event_id":   event.ID,
		"event_type": event.Type,
		"user_id":    userID,
		"plan":       plan,
		"is_active":  isActive,
		"is_trial":   event.IsTrial,
	})

	return nil
}

func (s *SubscriptionService) mapProductIDToPlan(productID string) enum.SubscriptionPlan {
	switch productID {
	case "monthly_premium", "premium_monthly":
		return enum.PlanMonthly
	case "yearly_premium", "premium_yearly":
		return enum.PlanYearly
	default:
		return enum.PlanFree
	}
}

func (s *SubscriptionService) isActiveEvent(eventType enum.RevenueCatEventType) bool {
	switch eventType {
	case enum.EventInitialPurchase, enum.EventRenewal, enum.EventUncancellation:
		return true
	case enum.EventCancellation, enum.EventExpiration, enum.EventBillingIssue:
		return false
	default:
		return false
	}
}

func stringPtr(s string) *string {
	return &s
}

func timePtr(t time.Time) *time.Time {
	return &t
}
